package memory

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

var (
	_ repository.Repository         = (*Repository)(nil)
	_ repository.RelationRepository = (*Repository)(nil)
	_ repository.ExtraRepository    = (*Repository)(nil)
)

type Repository struct {
	Tasks     map[string]map[domain.TaskID]*domain.Task
	Relations map[domain.RelationID]*domain.Relation
	Children  map[domain.RelationID][]domain.RelationID // parentID -> children
	Extras    map[domain.ExtraID]*domain.Extra
}

func NewRepository() *Repository {
	return &Repository{
		Tasks:     make(map[string]map[domain.TaskID]*domain.Task),
		Relations: make(map[domain.RelationID]*domain.Relation),
		Children:  make(map[domain.RelationID][]domain.RelationID),
		Extras:    make(map[domain.ExtraID]*domain.Extra),
	}
}

func (r *Repository) CreateTask(ctx context.Context, username string, task *domain.Task) error {
	if _, ok := r.Tasks[username]; !ok {
		r.Tasks[username] = make(map[domain.TaskID]*domain.Task)
	}

	r.Tasks[username][task.ID()] = task

	return nil
}

func (r *Repository) GetTask(ctx context.Context, username string, id domain.TaskID) (*domain.Task, error) {
	task, ok := r.Tasks[username][id]
	if !ok {
		return nil, repository.ErrTaskNotFound
	}

	return task, nil
}

func (r *Repository) CreateRelation(ctx context.Context, relation *domain.Relation) error {
	if _, ok := r.Relations[relation.ID()]; ok {
		return repository.ErrRelationAlreadyExists
	}

	r.Relations[relation.ID()] = relation
	r.Children[relation.ParentID()] = append(r.Children[relation.ParentID()], relation.ID())

	return nil
}

func (r *Repository) GetRelation(ctx context.Context, id domain.RelationID) (*domain.Relation, error) {
	relation, ok := r.Relations[id]
	if !ok {
		return nil, repository.ErrRelationNotFound
	}

	return relation, nil
}

func (r *Repository) ListChildrenRelations(ctx context.Context, id domain.RelationID) ([]*domain.Relation, error) {
	var ret []*domain.Relation
	for _, childID := range r.Children[id] {
		ret = append(ret, r.Relations[childID])
	}

	return ret, nil
}

func (r *Repository) DeleteRelation(ctx context.Context, relation *domain.Relation) error {
	if _, ok := r.Relations[relation.ID()]; !ok {
		return repository.ErrRelationNotFound
	}

	delete(r.Relations, relation.ID())
	delete(r.Children, relation.ParentID())

	return nil
}

func (r *Repository) UpdateRelation(ctx context.Context, relation *domain.Relation) error {
	if _, ok := r.Relations[relation.ID()]; !ok {
		return repository.ErrRelationNotFound
	}

	r.Relations[relation.ID()] = relation

	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, username string, task *domain.Task) error {
	if _, ok := r.Tasks[username][task.ID()]; !ok {
		return repository.ErrTaskNotFound
	}

	delete(r.Tasks[username], task.ID())

	return nil
}

func (r *Repository) CreateExtra(ctx context.Context, extra *domain.Extra) error {
	if _, ok := r.Extras[extra.ID()]; ok {
		return repository.ErrExtraAlreadyExists
	}

	r.Extras[extra.ID()] = extra

	return nil
}

func (r *Repository) DeleteExtra(ctx context.Context, extra *domain.Extra) error {
	if _, ok := r.Extras[extra.ID()]; !ok {
		return repository.ErrExtraNotFound
	}

	delete(r.Extras, extra.ID())

	return nil
}

func (r *Repository) GetExtra(ctx context.Context, id domain.ExtraID) (*domain.Extra, error) {
	extra, ok := r.Extras[id]
	if !ok {
		return nil, repository.ErrExtraNotFound
	}

	return extra, nil
}

func (r *Repository) ListExtras(ctx context.Context, ids []domain.ExtraID) ([]*domain.Extra, error) {
	var ret []*domain.Extra

	for _, id := range ids {
		extra, ok := r.Extras[id]
		if !ok {
			continue
		}

		ret = append(ret, extra)
	}

	return ret, nil
}

func (r *Repository) UpdateExtra(ctx context.Context, extra *domain.Extra) error {
	if _, ok := r.Extras[extra.ID()]; !ok {
		return repository.ErrExtraNotFound
	}

	r.Extras[extra.ID()] = extra

	return nil
}
