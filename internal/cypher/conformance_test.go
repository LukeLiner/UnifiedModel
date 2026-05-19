package cypher

import (
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

type conformanceCase struct {
	Name    string           `yaml:"name"`
	Query   string           `yaml:"query"`
	Params  map[string]any   `yaml:"params"`
	Columns []string         `yaml:"columns"`
	Rows    []map[string]any `yaml:"rows"`
}

func TestConformanceCorpus(t *testing.T) {
	raw, err := os.ReadFile("testdata/conformance.yaml")
	if err != nil {
		t.Fatalf("read conformance corpus: %v", err)
	}
	var cases []conformanceCase
	if err := yaml.Unmarshal(raw, &cases); err != nil {
		t.Fatalf("parse conformance corpus: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := Execute(tc.Query, testGraph(), tc.Params, Options{Limit: 100})
			if err != nil {
				t.Fatalf("execute: %v", err)
			}
			if !reflect.DeepEqual(result.Columns, tc.Columns) {
				t.Fatalf("columns mismatch\nwant: %#v\n got: %#v", tc.Columns, result.Columns)
			}
			if !reflect.DeepEqual(normalizeRows(result.Rows), normalizeRows(tc.Rows)) {
				t.Fatalf("rows mismatch\nwant: %#v\n got: %#v", normalizeRows(tc.Rows), normalizeRows(result.Rows))
			}
		})
	}
}

func normalizeRows(rows []map[string]any) []map[string]any {
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		next := make(map[string]any, len(row))
		for key, value := range row {
			next[key] = normalizeValue(value)
		}
		out = append(out, next)
	}
	return out
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case int:
		return int64(typed)
	case int64:
		return typed
	case []string:
		values := make([]any, 0, len(typed))
		for _, item := range typed {
			values = append(values, item)
		}
		return values
	case []any:
		values := make([]any, 0, len(typed))
		for _, item := range typed {
			values = append(values, normalizeValue(item))
		}
		return values
	default:
		return typed
	}
}
