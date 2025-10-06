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

var _ repository.Repository = (*Repository)(nil)

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

	err = db.AutoMigrate(&Task{}) //nolint:exhaustruct
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) CountSubtasks(ctx context.Context, username string, id domain.TaskID) (int, error) {
	searchParentID := sql.NullString{String: string(id), Valid: id != ""}

	tasks, err := gorm.G[Task](r.db).
		Where(&Task{ParentID: searchParentID, Username: username}). //nolint:exhaustruct
		Find(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count subtasks: %w", err)
	}

	return len(tasks), nil
}

func (r *Repository) CreateTask(ctx context.Context, username string, task *domain.Task) error {
	err := gorm.G[Task](r.db).Create(ctx, FromDomainTask(task, username))
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (r *Repository) GetTask(ctx context.Context, username string, id domain.TaskID) (*domain.Task, error) {
	task, err := gorm.G[Task](r.db).Where(&Task{ID: string(id), Username: username}).First(ctx) //nolint:exhaustruct
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrTaskNotFound
		}

		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task.ToDomain(), nil
}

func (r *Repository) ListSubTasks(
	ctx context.Context,
	username string,
	parentID domain.TaskID,
) ([]*domain.Task, error) {
	searchParentID := sql.NullString{String: string(parentID), Valid: parentID != ""}

	tasks, err := gorm.G[Task](r.db).
		Where(&Task{ParentID: searchParentID, Username: username}). //nolint:exhaustruct
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	var ret []*domain.Task
	for _, task := range tasks {
		ret = append(ret, task.ToDomain())
	}

	return ret, nil
}

func (r *Repository) ListDescendantsTasks(
	ctx context.Context,
	username string,
	parentID domain.TaskID,
) ([]*domain.Task, error) {
	var stack []string

	var ret []*domain.Task

	stack = append(stack, string(parentID))
	for len(stack) > 0 {
		searchParentID := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		tasks, err := gorm.G[Task](r.db).
			Where(&Task{ //nolint:exhaustruct
				ParentID: sql.NullString{String: searchParentID, Valid: searchParentID != ""},
				Username: username,
			}).
			Find(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tasks: %w", err)
		}

		for _, task := range tasks {
			ret = append(ret, task.ToDomain())
			stack = append(stack, task.ID)
		}
	}

	return ret, nil
}
