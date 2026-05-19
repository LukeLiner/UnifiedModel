package umodel

import (
	"context"
	"sync"

	apperrors "github.com/alibaba/UnifiedModel/pkg/errors"
	"github.com/alibaba/UnifiedModel/pkg/model"
)

type graphStore interface {
	PutUModelElements(ctx context.Context, batch model.UModelElementBatch) (model.WriteResult, error)
	GetUModelSnapshot(ctx context.Context, req model.UModelSnapshotRequest) (model.UModelSnapshot, error)
}

type Service struct {
	graph              graphStore
	mu                 sync.RWMutex
	indexes            map[string]*schemaIndex
	commonSchemaLoader CommonSchemaLoader
}

func NewService(graph graphStore) *Service {
	return &Service{
		graph:   graph,
		indexes: make(map[string]*schemaIndex),
	}
}

type CommonSchemaLoader interface {
	LoadCommonSchemaPacks(ctx context.Context, workspace string, packs []string) ([]model.UModelElement, error)
}

func (s *Service) SetCommonSchemaLoader(loader CommonSchemaLoader) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commonSchemaLoader = loader
}

func (s *Service) Validate(ctx context.Context, workspace string, elements []model.UModelElement) (model.ValidationResult, error) {
	for _, element := range elements {
		if element.Kind == "" || element.Domain == "" || element.Name == "" {
			return model.ValidationResult{
				Valid: false,
				Errors: []model.ErrorDetail{{
					Field:  "kind/domain/name",
					Reason: "umodel element kind, domain, and name are required",
				}},
			}, nil
		}
	}
	return model.ValidationResult{Valid: true}, nil
}

func (s *Service) PutElements(ctx context.Context, batch model.UModelElementBatch) (model.WriteResult, error) {
	if batch.Workspace == "" {
		return model.WriteResult{}, apperrors.New(apperrors.CodeInvalidArgument, "workspace is required")
	}
	validation, err := s.Validate(ctx, batch.Workspace, batch.Elements)
	if err != nil {
		return model.WriteResult{}, err
	}
	if !validation.Valid {
		return model.WriteResult{}, apperrors.WithDetails(apperrors.CodeValidationFailed, "umodel validation failed", map[string]string{
			"field": validation.Errors[0].Field,
		})
	}
	result, err := s.graph.PutUModelElements(ctx, batch)
	if err != nil {
		return model.WriteResult{}, err
	}
	s.mergeIndex(batch.Workspace, batch.Elements)
	return result, nil
}

func (s *Service) DeleteElements(ctx context.Context, workspace string, ids []string) (model.WriteResult, error) {
	items := make([]model.BatchItemResult, 0, len(ids))
	for _, id := range ids {
		items = append(items, model.BatchItemResult{
			ID:      id,
			OK:      false,
			Code:    string(apperrors.CodeNotImplemented),
			Message: "delete dependency checks are not implemented in the current service",
		})
	}
	return model.WriteResult{Failed: len(items), Items: items}, nil
}

func (s *Service) RebuildIndex(ctx context.Context, workspace string) error {
	if workspace == "" {
		return apperrors.New(apperrors.CodeInvalidArgument, "workspace is required")
	}
	snapshot, err := s.graph.GetUModelSnapshot(ctx, model.UModelSnapshotRequest{Workspace: workspace})
	if err != nil {
		return err
	}
	s.replaceIndex(workspace, snapshot.Elements)
	return nil
}

func (s *Service) ResolveEntitySet(ctx context.Context, ref model.EntityTypeRef) (model.EntitySetSchema, error) {
	if ref.Domain == "" || ref.Name == "" {
		return model.EntitySetSchema{}, apperrors.New(apperrors.CodeInvalidArgument, "entity set ref requires domain and name")
	}
	if element, ok := s.findIndexedElement("entity_set", ref.Domain, ref.Name); ok {
		return model.EntitySetSchema{Ref: ref, Fields: fieldMapFromElement(element)}, nil
	}
	return model.EntitySetSchema{Ref: ref}, nil
}

func (s *Service) ResolveRelationType(ctx context.Context, ref model.RelationTypeRef) (model.RelationSchema, error) {
	if ref.Type == "" {
		return model.RelationSchema{}, apperrors.New(apperrors.CodeInvalidArgument, "relation type is required")
	}
	if element, ok := s.findIndexedRelation(ref); ok {
		return model.RelationSchema{Ref: ref, Fields: fieldMapFromElement(element)}, nil
	}
	return model.RelationSchema{Ref: ref}, nil
}

func (s *Service) ValidateEntityPayload(ctx context.Context, payload model.EntityPayload) (model.ValidationResult, error) {
	required := []string{"__domain__", "__entity_type__", "__entity_id__", "__method__", "__first_observed_time__", "__last_observed_time__"}
	for _, field := range required {
		if payload[field] == nil || payload[field] == "" {
			return model.ValidationResult{
				Valid:  false,
				Errors: []model.ErrorDetail{{Field: field, Reason: "required CMS 2.0 entity field is missing"}},
			}, nil
		}
	}
	if id, ok := payload["__entity_id__"].(string); !ok || !model.IsEntityID(id) {
		return model.ValidationResult{
			Valid:  false,
			Errors: []model.ErrorDetail{{Field: "__entity_id__", Reason: "entity id must be a 128-bit lowercase hex string"}},
		}, nil
	}
	if !validWriteMethod(payload["__method__"]) {
		return model.ValidationResult{
			Valid:  false,
			Errors: []model.ErrorDetail{{Field: "__method__", Reason: "method must be Create, Update, Expire, or Delete"}},
		}, nil
	}
	return model.ValidationResult{Valid: true}, nil
}

func (s *Service) ValidateRelationPayload(ctx context.Context, payload model.RelationPayload) (model.ValidationResult, error) {
	required := []string{"__src_domain__", "__src_entity_type__", "__src_entity_id__", "__dest_domain__", "__dest_entity_type__", "__dest_entity_id__", "__relation_type__", "__method__", "__first_observed_time__", "__last_observed_time__"}
	for _, field := range required {
		if payload[field] == nil || payload[field] == "" {
			return model.ValidationResult{
				Valid:  false,
				Errors: []model.ErrorDetail{{Field: field, Reason: "required CMS 2.0 relation field is missing"}},
			}, nil
		}
	}
	for _, field := range []string{"__src_entity_id__", "__dest_entity_id__"} {
		if id, ok := payload[field].(string); !ok || !model.IsEntityID(id) {
			return model.ValidationResult{
				Valid:  false,
				Errors: []model.ErrorDetail{{Field: field, Reason: "entity id must be a 128-bit lowercase hex string"}},
			}, nil
		}
	}
	if !validWriteMethod(payload["__method__"]) {
		return model.ValidationResult{
			Valid:  false,
			Errors: []model.ErrorDetail{{Field: "__method__", Reason: "method must be Create, Update, Expire, or Delete"}},
		}, nil
	}
	return model.ValidationResult{Valid: true}, nil
}

func (s *Service) SnapshotVersion(ctx context.Context, workspace string) (model.SchemaVersion, error) {
	snapshot, err := s.graph.GetUModelSnapshot(ctx, model.UModelSnapshotRequest{Workspace: workspace})
	if err != nil {
		return model.SchemaVersion{}, err
	}
	return model.SchemaVersion{Workspace: workspace, Version: snapshot.Version}, nil
}

func validWriteMethod(value any) bool {
	method, ok := value.(string)
	if !ok {
		return false
	}
	switch method {
	case "Create", "Update", "Expire", "Delete":
		return true
	default:
		return false
	}
}
