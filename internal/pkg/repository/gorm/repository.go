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
	_ repository.Repository         = (*Repository)(nil)
	_ repository.RelationRepository = (*Repository)(nil)
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

func (r *Repository) CreateRelation(ctx context.Context, relation *domain.Relation) error {
	err := gorm.G[Relation](r.db).Create(ctx, FromDomainRelation(relation))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return repository.ErrRelationAlreadyExists
		}

		return fmt.Errorf("failed to create relation: %w", err)
	}

	return nil
}

func (r *Repository) DeleteRelation(ctx context.Context, relation *domain.Relation) error {
	affected, err := gorm.G[Relation](r.db).
		Where(&Relation{ID: string(relation.ID())}). //nolint:exhaustruct
		Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete relation: %w", err)
	}

	if affected == 0 {
		return repository.ErrRelationNotFound
	}

	return nil
}

func (r *Repository) GetRelation(ctx context.Context, id domain.RelationID) (*domain.Relation, error) {
	relation, err := gorm.G[Relation](r.db).
		Where(&Relation{ID: string(id)}). //nolint:exhaustruct
		Take(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrRelationNotFound
		}

		return nil, fmt.Errorf("failed to get relation: %w", err)
	}

	return relation.ToDomain(), nil
}

func (r *Repository) ListChildrenRelations(ctx context.Context, id domain.RelationID) ([]*domain.Relation, error) {
	relations, err := gorm.G[Relation](r.db).
		Where(&Relation{ParentID: sql.NullString{String: string(id), Valid: true}}). //nolint:exhaustruct
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list children relations: %w", err)
	}

	return ToDomainRelations(relations), nil
}

func (r *Repository) UpdateRelation(ctx context.Context, dRelataion *domain.Relation) error {
	relation := FromDomainRelation(dRelataion)
	relation.Version++

	affected, err := gorm.G[Relation](r.db).
		Where(&Relation{ID: string(dRelataion.ID()), Version: dRelataion.Version()}). //nolint:exhaustruct
		Updates(ctx, *relation)
	if err != nil {
		return fmt.Errorf("failed to update relation: %w", err)
	}

	if affected == 0 {
		return repository.ErrRelationBusy
	}

	return nil
}
