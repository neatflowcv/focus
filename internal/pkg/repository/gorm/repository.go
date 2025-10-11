package gorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	_ repository.Repository      = (*Repository)(nil)
	_ repository.ExtraRepository = (*Repository)(nil)
	_ repository.TraceRepository = (*Repository)(nil)
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() (*Repository, error) {
	dsn :=
		"host=127.0.0.1 user=focus password=password dbname=focus port=5432 sslmode=disable TimeZone=Asia/Seoul"

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{ //nolint:exhaustruct
				DSN:                  dsn,
				PreferSimpleProtocol: true,
			},
		),
		&gorm.Config{}) //nolint:exhaustruct
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.AutoMigrate(&Task{}, &Extra{}, &Trace{}) //nolint:exhaustruct
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) CreateTasks(ctx context.Context, username string, dTasks ...*domain.Task) error {
	var tasks []*Task
	for _, task := range dTasks {
		tasks = append(tasks, FromDomainTask(task, username))
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, task := range tasks {
			err := gorm.G[Task](tx).Create(ctx, task)
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create tasks: %w", err)
	}

	return nil
}

func (r *Repository) GetTask(ctx context.Context, username string, id domain.TaskID) (*domain.Task, error) {
	task, err := gorm.G[Task](r.db).
		Where(&Task{ID: string(id), Username: username}). //nolint:exhaustruct
		Take(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrTaskNotFound
		}

		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task.ToDomain(), nil
}

func (r *Repository) DeleteTasks(ctx context.Context, username string, tasks ...*domain.Task) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, task := range tasks {
			affected, err := gorm.G[Task](tx).
				Where(
					&Task{ //nolint:exhaustruct
						ID:       string(task.ID()),
						Username: username,
						Version:  task.Version(),
					},
				).
				Delete(ctx)
			if err != nil {
				return fmt.Errorf("failed to delete task: %w", err)
			}

			if affected == 0 {
				return repository.ErrTaskNotFound
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to delete tasks: %w", err)
	}

	return nil
}

func (r *Repository) ListTasks(ctx context.Context, username string, parentID domain.TaskID) ([]*domain.Task, error) {
	var tasks []Task

	tasks, err := gorm.G[Task](r.db).
		Where(
			&Task{ //nolint:exhaustruct
				Username: username,
				ParentID: sql.NullString{String: string(parentID), Valid: true},
			},
		).
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return ToDomainTasks(tasks), nil
}

func (r *Repository) UpdateTasks(ctx context.Context, username string, dTasks ...*domain.Task) error {
	var tasks []*Task
	for _, task := range dTasks {
		tasks = append(tasks, FromDomainTask(task, username))
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, task := range tasks {
			oldVersion := task.Version
			task.Version++

			affected, err := gorm.G[Task](tx).
				Where(
					&Task{ //nolint:exhaustruct
						ID:       task.ID,
						Username: username,
						Version:  oldVersion,
					},
				).
				Updates(ctx, *task)
			if err != nil {
				return fmt.Errorf("failed to update task: %w", err)
			}

			if affected == 0 {
				return repository.ErrTaskNotFound
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update tasks: %w", err)
	}

	return nil
}

func (r *Repository) CreateExtra(ctx context.Context, dExtra *domain.Extra) error {
	extra := FromDomainExtra(dExtra)

	err := gorm.G[Extra](r.db).Create(ctx, extra)
	if err != nil {
		return fmt.Errorf("failed to create extra: %w", err)
	}

	return nil
}

func (r *Repository) DeleteExtra(ctx context.Context, dExtra *domain.Extra) error {
	extra := FromDomainExtra(dExtra)

	affected, err := gorm.G[Extra](r.db).
		Where(extra).
		Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete extra: %w", err)
	}

	if affected == 0 {
		return repository.ErrExtraNotFound
	}

	return nil
}

func (r *Repository) GetExtra(ctx context.Context, id domain.ExtraID) (*domain.Extra, error) {
	extra, err := gorm.G[Extra](r.db).
		Where(&Extra{ID: string(id)}). //nolint:exhaustruct
		Take(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get extra: %w", err)
	}

	return extra.ToDomain(), nil
}

func (r *Repository) ListExtras(ctx context.Context, ids []domain.ExtraID) ([]*domain.Extra, error) {
	var ret []*domain.Extra

	for _, id := range ids {
		extra, err := gorm.G[Extra](r.db).
			Where(&Extra{ID: string(id)}). //nolint:exhaustruct
			Take(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list extras: %w", err)
		}

		ret = append(ret, extra.ToDomain())
	}

	return ret, nil
}

func (r *Repository) UpdateExtra(ctx context.Context, extra *domain.Extra) error {
	affected, err := gorm.G[Extra](r.db).
		Where(&Extra{ID: string(extra.ID())}). //nolint:exhaustruct
		Updates(ctx, *FromDomainExtra(extra))
	if err != nil {
		return fmt.Errorf("failed to update extra: %w", err)
	}

	if affected == 0 {
		return repository.ErrExtraNotFound
	}

	return nil
}

func (r *Repository) CreateTrace(ctx context.Context, dTrace *domain.Trace) error {
	trace := FromDomainTrace(dTrace)

	err := gorm.G[Trace](r.db).Create(ctx, trace)
	if err != nil {
		return fmt.Errorf("failed to create trace: %w", err)
	}

	return nil
}

func (r *Repository) DeleteTrace(ctx context.Context, trace *domain.Trace) error {
	affected, err := gorm.G[Trace](r.db).
		Where(&Trace{ID: string(trace.ID())}). //nolint:exhaustruct
		Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete trace: %w", err)
	}

	if affected == 0 {
		return repository.ErrTraceNotFound
	}

	return nil
}

func (r *Repository) GetTrace(ctx context.Context, id domain.TraceID) (*domain.Trace, error) {
	trace, err := gorm.G[Trace](r.db).
		Where(&Trace{ID: string(id)}). //nolint:exhaustruct
		Take(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get trace: %w", err)
	}

	return trace.ToDomain(), nil
}

func (r *Repository) UpdateTraces(ctx context.Context, dTraces ...*domain.Trace) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, dTrace := range dTraces {
			trace := FromDomainTrace(dTrace)

			affected, err := gorm.G[Trace](r.db).
				Where(&Trace{ID: trace.ID}). //nolint:exhaustruct
				Updates(ctx, *trace)
			if err != nil {
				return fmt.Errorf("failed to update traces: %w", err)
			}

			if affected == 0 {
				return repository.ErrTraceNotFound
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update traces: %w", err)
	}

	return nil
}

func (r *Repository) ListTraces(ctx context.Context, ids []domain.TraceID) ([]*domain.Trace, error) {
	traces, err := gorm.G[Trace](r.db).
		Where("id IN ?", ids).
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list traces: %w", err)
	}

	var ret []*domain.Trace
	for _, trace := range traces {
		ret = append(ret, trace.ToDomain())
	}

	return ret, nil
}
