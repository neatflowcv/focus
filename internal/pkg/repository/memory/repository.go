package memory

import (
	"context"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository"
)

var (
	_ repository.Repository      = (*Repository)(nil)
	_ repository.ExtraRepository = (*Repository)(nil)
	_ repository.TraceRepository = (*Repository)(nil)
)

type Repository struct {
	Tasks    map[string]map[domain.TaskID]*domain.Task
	Children map[domain.TaskID][]domain.TaskID
	Extras   map[domain.ExtraID]*domain.Extra
	Traces   map[domain.TraceID]*domain.Trace
}

func NewRepository() *Repository {
	return &Repository{
		Tasks:    make(map[string]map[domain.TaskID]*domain.Task),
		Children: make(map[domain.TaskID][]domain.TaskID),
		Extras:   make(map[domain.ExtraID]*domain.Extra),
		Traces:   make(map[domain.TraceID]*domain.Trace),
	}
}

func (r *Repository) CreateTask(ctx context.Context, username string, task *domain.Task) error {
	if _, ok := r.Tasks[username]; !ok {
		r.Tasks[username] = make(map[domain.TaskID]*domain.Task)
	}

	r.Tasks[username][task.ID()] = task
	r.Children[task.ParentID()] = append(r.Children[task.ParentID()], task.ID())

	return nil
}

func (r *Repository) GetTask(ctx context.Context, username string, id domain.TaskID) (*domain.Task, error) {
	task, ok := r.Tasks[username][id]
	if !ok {
		return nil, repository.ErrTaskNotFound
	}

	return task, nil
}

func (r *Repository) ListTasks(ctx context.Context, username string, parentID domain.TaskID) ([]*domain.Task, error) {
	var ret []*domain.Task

	for _, id := range r.Children[parentID] {
		task, ok := r.Tasks[username][id]
		if !ok {
			panic("task not found")
		}

		ret = append(ret, task)
	}

	return ret, nil
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

func (r *Repository) CreateTrace(ctx context.Context, trace *domain.Trace) error {
	if _, ok := r.Traces[trace.ID()]; ok {
		return repository.ErrTraceAlreadyExists
	}

	r.Traces[trace.ID()] = trace

	return nil
}

func (r *Repository) GetTrace(ctx context.Context, id domain.TraceID) (*domain.Trace, error) {
	trace, ok := r.Traces[id]
	if !ok {
		return nil, repository.ErrTraceNotFound
	}

	return trace, nil
}

func (r *Repository) UpdateTraces(ctx context.Context, traces ...*domain.Trace) error {
	for _, trace := range traces {
		if _, ok := r.Traces[trace.ID()]; !ok {
			return repository.ErrTraceNotFound
		}
	}

	for _, trace := range traces {
		r.Traces[trace.ID()] = trace
	}

	return nil
}

func (r *Repository) DeleteTrace(ctx context.Context, trace *domain.Trace) error {
	if _, ok := r.Traces[trace.ID()]; !ok {
		return repository.ErrTraceNotFound
	}

	delete(r.Traces, trace.ID())

	return nil
}
