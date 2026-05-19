package query

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	apperrors "github.com/alibaba/UnifiedModel/pkg/errors"
	"github.com/alibaba/UnifiedModel/pkg/model"
)

type Executor struct {
	graph graphStore
}

func NewExecutor(graph graphStore) *Executor {
	return &Executor{graph: graph}
}

func (e *Executor) Execute(ctx context.Context, workspace string, plan model.QueryPlan) (model.QueryResult, error) {
	plan.Workspace = workspace

	var result model.QueryResult
	var err error
	switch plan.Source {
	case ".umodel":
		result, err = e.executeUModel(ctx, workspace, plan)
	case ".entity":
		result, err = e.graph.QueryEntities(ctx, model.EntityQueryPlan(plan))
	case ".topo":
		result, err = e.graph.QueryTopo(ctx, model.TopoQueryPlan(plan))
	default:
		return model.QueryResult{}, apperrors.New(apperrors.CodeQueryPlanError, "unsupported query source")
	}
	if err != nil {
		return model.QueryResult{}, err
	}

	rows, columns := applyPipeline(plan.Source, result.Rows, result.Columns, plan)
	result.Rows = rows
	result.Columns = columns
	result.Page = model.PageRequest{Limit: plan.Limit}
	return result, nil
}

func (e *Executor) executeUModel(ctx context.Context, workspace string, plan model.QueryPlan) (model.QueryResult, error) {
	snapshot, err := e.graph.GetUModelSnapshot(ctx, model.UModelSnapshotRequest{Workspace: workspace})
	if err != nil {
		return model.QueryResult{}, err
	}
	rows := make([]map[string]any, 0, len(snapshot.Elements))
	for _, element := range snapshot.Elements {
		rows = append(rows, map[string]any{
			"kind":    element.Kind,
			"domain":  element.Domain,
			"name":    element.Name,
			"version": element.Version,
			"spec":    element.Spec,
		})
	}
	return model.QueryResult{
		Columns: []string{"kind", "domain", "name", "version"},
		Rows:    rows,
		Page:    model.PageRequest{Limit: plan.Limit},
	}, nil
}

func applyPipeline(source string, rows []map[string]any, columns []string, plan model.QueryPlan) ([]map[string]any, []string) {
	if !hasOperator(plan.Pipeline, "with") && len(plan.Filters) > 0 {
		rows = filterRows(source, rows, plan.Filters)
	}

	for _, operator := range plan.Pipeline {
		switch operator.Name {
		case "with":
			rows = filterRows(source, rows, plan.Filters)
		case "where":
			if operator.Predicate != nil {
				rows = filterPredicate(rows, *operator.Predicate)
			}
		case "project":
			rows, columns = projectRows(rows, operator.Project)
		case "sort":
			if operator.Sort != nil {
				sortRows(rows, *operator.Sort)
			}
		case "limit":
			rows = limitRows(rows, operator.Limit)
		}
	}

	rows = limitRows(rows, plan.Limit)
	return rows, columns
}

func hasOperator(operators []model.QueryPipelineOperator, name string) bool {
	for _, operator := range operators {
		if operator.Name == name {
			return true
		}
	}
	return false
}

func filterRows(source string, rows []map[string]any, filters map[string]any) []map[string]any {
	if len(filters) == 0 {
		return rows
	}
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		if rowMatchesFilters(source, row, filters) {
			out = append(out, row)
		}
	}
	return out
}

func rowMatchesFilters(source string, row map[string]any, filters map[string]any) bool {
	switch source {
	case ".umodel":
		if _, ok := filters["id"]; ok {
			return false
		}
		if !stringMatches(rowString(row, "kind"), filters["kind"]) {
			return false
		}
		if !stringMatches(rowString(row, "domain"), filters["domain"]) {
			return false
		}
		if !stringMatches(rowString(row, "name"), filters["name"]) {
			return false
		}
	case ".entity":
		if !stringMatches(rowString(row, "__domain__"), filters["domain"]) {
			return false
		}
		if !stringMatches(rowString(row, "__entity_type__"), filters["name"]) {
			return false
		}
		if !matchesIDs(rowString(row, "__entity_id__"), filters["ids"]) {
			return false
		}
	case ".topo":
		relationType := rowString(row, "__relation_type__")
		if relationType == "" {
			relationType = rowString(row, "relation")
		}
		if !stringMatches(relationType, coalesce(filters["relation_type"], filters["type"])) {
			return false
		}
		if !stringMatches(rowString(row, "src"), filters["src"]) {
			return false
		}
		if !stringMatches(rowString(row, "dest"), filters["dest"]) {
			return false
		}
	}

	query := stringFilter(filters["query"])
	if query != "" && !rowContains(row, query) {
		return false
	}
	return true
}

func filterPredicate(rows []map[string]any, predicate model.QueryPredicate) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		if predicateMatches(row, predicate) {
			out = append(out, row)
		}
	}
	return out
}

func predicateMatches(row map[string]any, predicate model.QueryPredicate) bool {
	left, ok := row[predicate.Field]
	if !ok {
		return false
	}
	switch predicate.Op {
	case "=", "==":
		return compareEqual(left, predicate.Value)
	case "!=":
		return !compareEqual(left, predicate.Value)
	case "contains", "~":
		return containsFold(stringValue(left), stringValue(predicate.Value))
	case ">", ">=", "<", "<=":
		return compareOrdered(left, predicate.Value, predicate.Op)
	default:
		return false
	}
}

func projectRows(rows []map[string]any, fields []string) ([]map[string]any, []string) {
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		next := make(map[string]any, len(fields))
		for _, field := range fields {
			next[field] = row[field]
		}
		out = append(out, next)
	}
	return out, append([]string(nil), fields...)
}

func sortRows(rows []map[string]any, sortSpec model.QuerySort) {
	sort.SliceStable(rows, func(i, j int) bool {
		cmp := compareForSort(rows[i][sortSpec.Field], rows[j][sortSpec.Field])
		if sortSpec.Desc {
			return cmp > 0
		}
		return cmp < 0
	})
}

func limitRows(rows []map[string]any, limit int) []map[string]any {
	if limit <= 0 || len(rows) <= limit {
		return rows
	}
	return rows[:limit]
}

func compareEqual(left, right any) bool {
	if lf, ok := floatValue(left); ok {
		if rf, ok := floatValue(right); ok {
			return lf == rf
		}
	}
	return stringValue(left) == stringValue(right)
}

func compareOrdered(left, right any, op string) bool {
	lf, lok := floatValue(left)
	rf, rok := floatValue(right)
	if lok && rok {
		switch op {
		case ">":
			return lf > rf
		case ">=":
			return lf >= rf
		case "<":
			return lf < rf
		case "<=":
			return lf <= rf
		}
	}
	lv := stringValue(left)
	rv := stringValue(right)
	switch op {
	case ">":
		return lv > rv
	case ">=":
		return lv >= rv
	case "<":
		return lv < rv
	case "<=":
		return lv <= rv
	default:
		return false
	}
}

func compareForSort(left, right any) int {
	if lf, ok := floatValue(left); ok {
		if rf, ok := floatValue(right); ok {
			switch {
			case lf < rf:
				return -1
			case lf > rf:
				return 1
			default:
				return 0
			}
		}
	}
	lv := stringValue(left)
	rv := stringValue(right)
	switch {
	case lv < rv:
		return -1
	case lv > rv:
		return 1
	default:
		return 0
	}
}

func floatValue(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case int32:
		return float64(typed), true
	case float64:
		return typed, true
	case float32:
		return float64(typed), true
	case string:
		n, err := strconv.ParseFloat(typed, 64)
		return n, err == nil
	default:
		return 0, false
	}
}

func rowContains(row map[string]any, query string) bool {
	for _, value := range row {
		if containsFold(stringValue(value), query) {
			return true
		}
	}
	return false
}

func rowString(row map[string]any, key string) string {
	return stringValue(row[key])
}

func matchesIDs(value string, filter any) bool {
	if filter == nil {
		return true
	}
	switch ids := filter.(type) {
	case []string:
		for _, id := range ids {
			if id == value {
				return true
			}
		}
		return false
	case []any:
		for _, id := range ids {
			if stringValue(id) == value {
				return true
			}
		}
		return false
	default:
		return stringValue(filter) == "" || stringValue(filter) == value
	}
}

func stringMatches(value string, filter any) bool {
	expected := stringFilter(filter)
	if expected == "" || expected == "*" {
		return true
	}
	if strings.HasSuffix(expected, "*") {
		return strings.HasPrefix(value, strings.TrimSuffix(expected, "*"))
	}
	return value == expected
}

func stringFilter(value any) string {
	if value == nil {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	case nil:
		return ""
	default:
		return fmt.Sprint(value)
	}
}

func containsFold(value, query string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(query))
}

func coalesce(values ...any) any {
	for _, value := range values {
		if stringValue(value) != "" {
			return value
		}
	}
	return nil
}
