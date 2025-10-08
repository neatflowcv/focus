package gorm

import (
	"context"
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

func (r *Repository) CreateTask(ctx context.Context, username string, task *domain.Task) error {
	err := gorm.G[Task](r.db).Create(ctx, FromDomainTask(task, username))
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
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
