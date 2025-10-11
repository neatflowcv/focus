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
		ParentID: "",
		NextID:   "",
	})

	require.NoError(t, err)
	require.NotEmpty(t, ret.ID)
	require.Equal(t, now, ret.CreatedAt)
}

func TestServiceListTasksWithNoData(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
		Username: "test",
		ParentID: "unknown",
	})

	require.NoError(t, err)
	require.Empty(t, ret)
}

func TestServiceListTasks(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	now := time.Now()
	_, _ = service.CreateTask(t.Context(), &flow.CreateTaskInput{
		Username: "test",
		Title:    "test",
		Now:      now,
		ParentID: "",
		NextID:   "",
	})

	ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
		Username: "test",
		ParentID: "",
	})

	require.NoError(t, err)
	require.Len(t, ret.Tasks, 1)
	require.Equal(t, "test", ret.Tasks[0].Title)
	require.Equal(t, now, ret.Tasks[0].CreatedAt)
}

func TestServiceDeleteTask(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	task, _ := service.CreateTask(t.Context(), &flow.CreateTaskInput{
		Username: "test",
		Title:    "test",
		Now:      time.Now(),
		ParentID: "",
		NextID:   "",
	})

	err := service.DeleteTask(t.Context(), &flow.DeleteTaskInput{
		Username: "test",
		TaskID:   task.ID,
	})

	require.NoError(t, err)
	require.Empty(t, data.repo.Tasks["test"][domain.TaskID(task.ID)])
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
