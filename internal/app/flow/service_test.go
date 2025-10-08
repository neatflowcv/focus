package flow_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/eventbus"
	"github.com/neatflowcv/focus/internal/pkg/idmaker/ulid"
	"github.com/neatflowcv/focus/internal/pkg/repository/memory"
	"github.com/stretchr/testify/require"
)

type ServiceData struct {
	bus     *eventbus.Bus
	idmaker *ulid.IDMaker
	repo    *memory.Repository
}

func newService(t *testing.T) (*flow.Service, *ServiceData) {
	t.Helper()

	data := &ServiceData{
		bus:     eventbus.NewBus(),
		idmaker: ulid.NewIDMaker(),
		repo:    memory.NewRepository(),
	}

	return flow.NewService(data.bus, data.idmaker, data.repo), data
}

func TestServiceCreateTask(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	now := time.Now()

	ret, err := service.CreateTask(t.Context(), &flow.CreateTaskInput{
		Username: "test",
		Title:    "test",
		Now:      now,
	})

	require.NoError(t, err)
	require.Equal(t, "test", ret.Task.Title)
	require.Equal(t, now, ret.Task.CreatedAt)
	require.NotEmpty(t, ret.Task.ID)
	require.Equal(t, string(domain.TaskStatusTodo), ret.Task.Status)
}

func TestServiceListTasksWithNoData(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
		Username: "test",
		IDs:      []string{"unknown"},
	})

	require.NoError(t, err)
	require.Empty(t, ret)
}

func TestServiceListTasks(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	now := time.Now()
	task, _ := service.CreateTask(t.Context(), &flow.CreateTaskInput{
		Username: "test",
		Title:    "test",
		Now:      now,
	})

	ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
		Username: "test",
		IDs:      []string{task.Task.ID},
	})

	require.NoError(t, err)
	require.Len(t, ret.Tasks, 1)
	require.Equal(t, "test", ret.Tasks[0].Title)
	require.Equal(t, now, ret.Tasks[0].CreatedAt)
	require.Equal(t, string(domain.TaskStatusTodo), ret.Tasks[0].Status)
}

func TestServiceDeleteTask(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	task, _ := service.CreateTask(t.Context(), &flow.CreateTaskInput{
		Username: "test",
		Title:    "test",
		Now:      time.Now(),
	})

	err := service.DeleteTask(t.Context(), &flow.DeleteTaskInput{
		Username: "test",
		TaskID:   task.Task.ID,
	})

	require.NoError(t, err)
	require.Empty(t, data.repo.Tasks["test"][domain.TaskID(task.Task.ID)])
}

func TestServiceDeleteTask_Error(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	err := service.DeleteTask(t.Context(), &flow.DeleteTaskInput{
		Username: "test",
		TaskID:   "test",
	})

	require.Error(t, err)
}
