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

	err = db.AutoMigrate(&Task{}, &TaskRelation{}, &TaskAugment{}) //nolint:exhaustruct
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &Repository{db: db}, nil
}

func (r *Repository) CountSubtasks(ctx context.Context, username string, id domain.TaskID) (int, error) {
	tasks, err := gorm.G[TaskRelation](r.db).
		Where(&TaskRelation{ParentID: string(id)}). //nolint:exhaustruct
		Find(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count subtasks: %w", err)
	}

	return len(tasks), nil
}

func (r *Repository) CreateTask(ctx context.Context, username string, task *domain.Task) error { //nolint:cyclop,funlen
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := gorm.G[Task](tx).Create(ctx, FromDomainTask(task, username))
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		err = gorm.G[TaskRelation](tx).Create(ctx, FromDomainTaskRelation(task))
		if err != nil {
			return fmt.Errorf("failed to create task relation: %w", err)
		}

		err = gorm.G[TaskAugment](tx).Create(ctx, &TaskAugment{
			TaskID:             string(task.ID()),
			Leaf:               true,
			AllDescendantsDone: true,
			TotalEstimatedTime: 0,
			TotalActualTime:    0,
		})
		if err != nil {
			return fmt.Errorf("failed to create task augment: %w", err)
		}

		parentID := task.ParentID()
		for parentID != "" {
			parentAugment, err := gorm.G[TaskAugment](tx).
				Where(&TaskAugment{TaskID: string(parentID)}). //nolint:exhaustruct
				Take(ctx)
			if err != nil {
				return fmt.Errorf("failed to get parent task augment: %w", err)
			}

			if !parentAugment.Leaf && !parentAugment.AllDescendantsDone {
				// 변함이 없으므로, 종료
				break
			}

			parentAugment.Leaf = false
			parentAugment.AllDescendantsDone = false

			_, err = gorm.G[TaskAugment](tx).
				Where(&TaskAugment{TaskID: string(parentID)}). //nolint:exhaustruct
				Select("leaf", "all_descendants_done").
				Updates(ctx, parentAugment)
			if err != nil {
				return fmt.Errorf("failed to update parent task augment: %w", err)
			}

			parentTask, err := gorm.G[TaskRelation](tx).
				Where(&TaskRelation{TaskID: string(parentID)}). //nolint:exhaustruct
				Take(ctx)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					break
				}

				return fmt.Errorf("failed to get parent task: %w", err)
			}

			parentID = domain.TaskID(parentTask.ParentID)
		}

		return nil
	})
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

	relation, err := gorm.G[TaskRelation](r.db).Where(&TaskRelation{TaskID: string(id)}).First(ctx) //nolint:exhaustruct
	if err != nil {
		return nil, fmt.Errorf("failed to get task relation: %w", err)
	}

	return task.ToDomain(domain.TaskID(relation.ParentID)), nil
}

func (r *Repository) ListSubTasks(
	ctx context.Context,
	username string,
	parentID domain.TaskID,
) ([]*domain.Task, error) {
	searchParentID := string(parentID)

	relations, err := gorm.G[TaskRelation](r.db).
		Where(&TaskRelation{ParentID: searchParentID}). //nolint:exhaustruct
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list task relations: %w", err)
	}

	var taskIDs []string
	for _, relation := range relations {
		taskIDs = append(taskIDs, relation.TaskID)
	}

	tasks, err := gorm.G[Task](r.db).
		Where("id IN (?)", taskIDs).
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	var ret []*domain.Task
	for _, task := range tasks {
		ret = append(ret, task.ToDomain(parentID))
	}

	return ret, nil
}

func (r *Repository) ListDescendantsTasks(
	ctx context.Context,
	username string,
	parentID domain.TaskID,
) ([]*domain.Task, error) {
	var stack []string

	stack = append(stack, string(parentID))

	var ret []*domain.Task

	for len(stack) > 0 {
		searchParentID := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		relations, err := gorm.G[TaskRelation](r.db).
			Where(&TaskRelation{ //nolint:exhaustruct
				ParentID: searchParentID,
			}).
			Find(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list task relations: %w", err)
		}

		var taskIDs []string
		for _, relation := range relations {
			taskIDs = append(taskIDs, relation.TaskID)
		}

		tasks, err := gorm.G[Task](r.db).
			Where("id IN (?)", taskIDs).
			Find(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list tasks: %w", err)
		}

		for _, task := range tasks {
			ret = append(ret, task.ToDomain(domain.TaskID(searchParentID)))
		}
	}

	return ret, nil
}
