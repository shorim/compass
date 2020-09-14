package eventdef

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/pkg/errors"
)

const eventAPIDefTable string = `"public"."event_api_definitions"`

var (
	idColumn      = "id"
	tenantColumn  = "tenant_id"
	bundleColumn  = "bundle_id"
	apiDefColumns = []string{idColumn, "od_id", tenantColumn, bundleColumn, "title", "short_description", "description", "group_name",
		"event_definitions", "tags", "documentation", "changelog_entries", "logo", "image", "url", "release_status", "last_updated", "extensions",
		"spec_data", "spec_format", "spec_type", "version", "version_deprecated", "version_deprecated_since", "version_for_removal"}
	idColumns        = []string{"id"}
	updatableColumns = []string{"title", "short_description", "description", "group_name",
		"event_definitions", "tags", "documentation", "changelog_entries", "logo", "image", "url", "release_status", "last_updated", "extensions",
		"spec_data", "spec_format", "spec_type", "version", "version_deprecated", "version_deprecated_since", "version_for_removal"}
)

//go:generate mockery -name=EventAPIDefinitionConverter -output=automock -outpkg=automock -case=underscore
type EventAPIDefinitionConverter interface {
	FromEntity(entity Entity) (model.EventDefinition, error)
	ToEntity(apiModel model.EventDefinition) (Entity, error)
	EventSpecFromEntity(specEnt EntitySpec) *model.EventSpec
}

//go:generate mockery -name=SpecRepository -output=automock -outpkg=automock -case=underscore
type SpecRepository interface {
	ListForAPI(ctx context.Context, tenantID, apiID string) ([]*model.Spec, error)
	ListForEvent(ctx context.Context, tenantID, eventID string) ([]*model.Spec, error)
	Exists(ctx context.Context, tenant, id string) (bool, error)
	GetByID(ctx context.Context, tenantID, id string) (*model.Spec, error)
	CreateMany(ctx context.Context, item []*model.Spec) error
	Create(ctx context.Context, item *model.Spec) error
	Update(ctx context.Context, item *model.Spec) error
	Delete(ctx context.Context, tenantID string, id string) error
}

type pgRepository struct {
	singleGetter    repo.SingleGetter
	pageableQuerier repo.PageableQuerier
	creator         repo.Creator
	updater         repo.Updater
	deleter         repo.Deleter
	existQuerier    repo.ExistQuerier
	conv            EventAPIDefinitionConverter
	specRepo        SpecRepository
}

func NewRepository(conv EventAPIDefinitionConverter, specRepo SpecRepository) *pgRepository {
	return &pgRepository{
		singleGetter:    repo.NewSingleGetter(resource.EventDefinition, eventAPIDefTable, tenantColumn, apiDefColumns),
		pageableQuerier: repo.NewPageableQuerier(resource.EventDefinition, eventAPIDefTable, tenantColumn, apiDefColumns),
		creator:         repo.NewCreator(resource.EventDefinition, eventAPIDefTable, apiDefColumns),
		updater:         repo.NewUpdater(resource.EventDefinition, eventAPIDefTable, updatableColumns, tenantColumn, idColumns),
		deleter:         repo.NewDeleter(resource.EventDefinition, eventAPIDefTable, tenantColumn),
		existQuerier:    repo.NewExistQuerier(resource.EventDefinition, eventAPIDefTable, tenantColumn),
		conv:            conv,
		specRepo:        specRepo,
	}
}

type EventAPIDefCollection []Entity

func (r EventAPIDefCollection) Len() int {
	return len(r)
}

func (r *pgRepository) GetByID(ctx context.Context, tenantID string, id string) (*model.EventDefinition, error) {
	var eventEntity Entity
	err := r.singleGetter.Get(ctx, tenantID, repo.Conditions{repo.NewEqualCondition("id", id)}, repo.NoOrderBy, &eventEntity)
	if err != nil {
		return nil, errors.Wrap(err, "while getting EventDefinition")
	}

	eventModel, err := r.conv.FromEntity(eventEntity)
	if err != nil {
		return nil, err
	}

	if !eventEntity.SpecData.Valid {
		specs, err := r.specRepo.ListForEvent(ctx, tenantID, id)
		if err != nil {
			return nil, err
		}
		apiSpecs := make([]*model.EventSpec, 0, 0)
		for _, spec := range specs {
			apiSpecs = append(apiSpecs, spec.ToEventSpec())
		}
		eventModel.Specs = append(eventModel.Specs, apiSpecs...)
	} else {
		eventModel.Specs = append(eventModel.Specs, r.conv.EventSpecFromEntity(eventEntity.EntitySpec))
	}

	return &eventModel, nil
}

func (r *pgRepository) ExistsByCondition(ctx context.Context, tenant string, conds repo.Conditions) (bool, error) {
	return r.existQuerier.Exists(ctx, tenant, conds)
}

func (r *pgRepository) GetByConditions(ctx context.Context, tenant string, conds repo.Conditions) (*model.EventDefinition, error) {
	var eventAPIDefEntity Entity
	err := r.singleGetter.Get(ctx, tenant, conds, repo.NoOrderBy, &eventAPIDefEntity)
	if err != nil {
		return nil, errors.Wrap(err, "while getting EventDefinition")
	}

	eventAPIDefModel, err := r.conv.FromEntity(eventAPIDefEntity)
	if err != nil {
		return nil, errors.Wrap(err, "while creating EventDefinition entity to model")
	}

	return &eventAPIDefModel, nil
}

func (r *pgRepository) GetForBundle(ctx context.Context, tenant string, id string, bundleID string) (*model.EventDefinition, error) {
	var ent Entity

	conditions := repo.Conditions{
		repo.NewEqualCondition(idColumn, id),
		repo.NewEqualCondition(bundleColumn, bundleID),
	}
	if err := r.singleGetter.Get(ctx, tenant, conditions, repo.NoOrderBy, &ent); err != nil {
		return nil, err
	}

	eventAPIModel, err := r.conv.FromEntity(ent)
	if err != nil {
		return nil, errors.Wrap(err, "while creating event definition model from entity")
	}

	return &eventAPIModel, nil
}

func (r *pgRepository) ListForBundle(ctx context.Context, tenantID string, bundleID string, pageSize int, cursor string) (*model.EventDefinitionPage, error) {
	conditions := repo.Conditions{
		repo.NewEqualCondition(bundleColumn, bundleID),
	}

	return r.list(ctx, tenantID, pageSize, cursor, conditions)
}

func (r *pgRepository) list(ctx context.Context, tenant string, pageSize int, cursor string, conditions repo.Conditions) (*model.EventDefinitionPage, error) {
	var eventCollection EventAPIDefCollection
	page, totalCount, err := r.pageableQuerier.List(ctx, tenant, pageSize, cursor, "id", &eventCollection, conditions...)
	if err != nil {
		return nil, err
	}

	var items []*model.EventDefinition

	for _, eventEnt := range eventCollection {
		m, err := r.conv.FromEntity(eventEnt)
		if err != nil {
			return nil, err
		}
		if !eventEnt.SpecData.Valid {
			specs, err := r.specRepo.ListForEvent(ctx, tenant, m.ID)
			if err != nil {
				return nil, err
			}
			eventSpecs := make([]*model.EventSpec, 0, 0)
			for _, spec := range specs {
				eventSpecs = append(eventSpecs, spec.ToEventSpec())
			}
			m.Specs = append(m.Specs, eventSpecs...)
		} else {
			m.Specs = append(m.Specs, r.conv.EventSpecFromEntity(eventEnt.EntitySpec))
		}
		items = append(items, &m)
	}

	return &model.EventDefinitionPage{
		Data:       items,
		TotalCount: totalCount,
		PageInfo:   page,
	}, nil
}

func (r *pgRepository) Create(ctx context.Context, item *model.EventDefinition) error {
	if item == nil {
		return apperrors.NewInternalError("item cannot be nil")
	}

	entity, err := r.conv.ToEntity(*item)
	if err != nil {
		return errors.Wrap(err, "while creating EventDefinition model to entity")
	}

	err = r.creator.Create(ctx, entity)
	if err != nil {
		return errors.Wrap(err, "while saving entity to db")
	}

	return nil
}

func (r *pgRepository) CreateMany(ctx context.Context, items []*model.EventDefinition) error {
	for index, item := range items {
		entity, err := r.conv.ToEntity(*item)
		if err != nil {
			return errors.Wrapf(err, "while creating %d item", index)
		}
		err = r.creator.Create(ctx, entity)
		if err != nil {
			return errors.Wrapf(err, "while persisting %d item", index)
		}
	}

	return nil
}

func (r *pgRepository) Update(ctx context.Context, item *model.EventDefinition) error {
	if item == nil {
		return apperrors.NewInternalError("item cannot be nil")
	}

	entity, err := r.conv.ToEntity(*item)
	if err != nil {
		return errors.Wrap(err, "while converting model to entity")
	}

	return r.updater.UpdateSingle(ctx, entity)
}

func (r *pgRepository) Exists(ctx context.Context, tenantID, id string) (bool, error) {
	return r.existQuerier.Exists(ctx, tenantID, repo.Conditions{repo.NewEqualCondition(idColumn, id)})
}

func (r *pgRepository) Delete(ctx context.Context, tenantID string, id string) error {
	return r.deleter.DeleteOne(ctx, tenantID, repo.Conditions{repo.NewEqualCondition(idColumn, id)})
}
