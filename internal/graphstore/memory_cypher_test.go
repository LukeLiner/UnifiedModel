package graphstore

import (
	"context"
	"testing"

	apperrors "github.com/alibaba/UnifiedModel/pkg/errors"
	"github.com/alibaba/UnifiedModel/pkg/model"
)

func TestMemoryCypherQueryTopoPath(t *testing.T) {
	store := NewMemoryStore()
	writeCypherRelationFixture(t, store)

	plan := model.TopoQueryPlan{
		Workspace: "demo",
		Limit:     20,
		GraphCall: &model.GraphCallPlan{
			Name:   "cypher",
			Cypher: "match (svc:`apm@apm.service` {__entity_id__: $svc}) optional match path = (svc)-[r*1..2]-(neighbor) with svc, neighbor, relationships(path) as rels where neighbor is null or coalesce(neighbor.__deleted__, false) = false return svc.__entity_id__ as service, neighbor.__entity_id__ as neighbor, [rel in rels | type(rel)] as relation_types, size(rels) as hops order by hops, neighbor limit 20",
		},
		Params: map[string]any{"svc": "54013ba69c196820e56801f1ef5aad54"},
	}
	result, err := store.QueryTopo(context.Background(), plan)
	if err != nil {
		t.Fatalf("query topo: %v", err)
	}
	if len(result.Rows) != 1 || result.Rows[0]["service"] != "54013ba69c196820e56801f1ef5aad54" || result.Rows[0]["neighbor"] != "177627f91af678a9b03e993f1a91917f" {
		t.Fatalf("unexpected rows: %+v", result.Rows)
	}
	if types, ok := result.Rows[0]["relation_types"].([]string); !ok || len(types) != 1 || types[0] != "calls" {
		t.Fatalf("unexpected relation types: %#v", result.Rows[0]["relation_types"])
	}
}

func TestMemoryCypherRejectsReadableEntityID(t *testing.T) {
	store := NewMemoryStore()
	_, err := store.QueryTopo(context.Background(), model.TopoQueryPlan{
		Workspace: "demo",
		Limit:     10,
		GraphCall: &model.GraphCallPlan{Name: "cypher", Cypher: "MATCH (svc:`apm@apm.service` {__entity_id__: 'cart'}) RETURN svc"},
	})
	if !apperrors.IsCode(err, apperrors.CodeQueryPlanError) {
		t.Fatalf("expected non-hex entity id to be rejected, got %v", err)
	}
}

func writeCypherRelationFixture(t *testing.T, store *MemoryStore) {
	t.Helper()
	if _, err := store.WriteRelations(context.Background(), model.RelationWriteBatch{
		Workspace: "demo",
		Relations: []model.RelationPayload{{
			"__src_domain__":       "apm",
			"__src_entity_type__":  "apm.service",
			"__src_entity_id__":    "54013ba69c196820e56801f1ef5aad54",
			"__dest_domain__":      "apm",
			"__dest_entity_type__": "apm.service",
			"__dest_entity_id__":   "177627f91af678a9b03e993f1a91917f",
			"__relation_type__":    "calls",
			"__method__":           "Update",
		}},
	}); err != nil {
		t.Fatalf("write relation: %v", err)
	}
}
