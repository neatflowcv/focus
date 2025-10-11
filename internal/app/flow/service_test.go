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

	const (
		username = "test"
		title    = "test"
		parentID = ""
		nextID   = ""
	)

	service, data := newService(t)
	now := time.Now()
	_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
		Username: username,
	})

	ret, err := service.CreateTask(t.Context(), &flow.CreateTaskInput{
		Username: username,
		Title:    title,
		Now:      now,
		ParentID: parentID,
		NextID:   nextID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, ret.ID)
	require.Equal(t, domain.TaskID(parentID), data.repo.Tasks[username][domain.TaskID(ret.ID)].ParentID())
	require.Equal(t, domain.TaskID(nextID), data.repo.Tasks[username][domain.TaskID(ret.ID)].NextID())
	require.Equal(t, title, data.repo.Tasks[username][domain.TaskID(ret.ID)].Title())
	require.Equal(t, now, ret.CreatedAt)
	require.Equal(t, now, data.repo.Tasks[username][domain.TaskID(ret.ID)].CreatedAt())
	require.Equal(t, uint64(1), ret.Version)
	require.Equal(t, uint64(1), data.repo.Tasks[username][domain.TaskID(ret.ID)].Version())
}

func TestServiceCreateTask_Error(t *testing.T) {
	t.Parallel()

	t.Run("unknown parent", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(t)
		_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
			Username: "test",
		})

		_, err := service.CreateTask(t.Context(), &flow.CreateTaskInput{
			Username: "test",
			Title:    "test",
			Now:      time.Now(),
			ParentID: "unknown",
			NextID:   "",
		})

		require.ErrorIs(t, err, flow.ErrParentTaskNotFound)
	})

	t.Run("unknown next", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(t)
		_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
			Username: "test",
		})

		_, err := service.CreateTask(t.Context(), &flow.CreateTaskInput{
			Username: "test",
			Title:    "test",
			Now:      time.Now(),
			ParentID: "",
			NextID:   "unknown",
		})

		require.ErrorIs(t, err, flow.ErrNextTaskNotFound)
	})
}

func TestServiceListTasks(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(t)
		_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
			Username: "test",
		})

		ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
			Username: "test",
			ParentID: "",
		})

		require.NoError(t, err)
		require.Empty(t, ret.Tasks)
	})

	t.Run("one", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(t)
		_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
			Username: "test",
		})
		_, _ = service.CreateTask(t.Context(), &flow.CreateTaskInput{
			Username: "test",
			Title:    "test",
			Now:      time.Now(),
			ParentID: "",
			NextID:   "",
		})

		ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
			Username: "test",
			ParentID: "",
		})

		require.NoError(t, err)
		require.Len(t, ret.Tasks, 1)
	})

	t.Run("two", func(t *testing.T) {
		t.Parallel()

		const (
			firstTitle  = "first"
			secondTitle = "second"
		)

		service, _ := newService(t)
		_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
			Username: "test",
		})
		_, _ = service.CreateTask(t.Context(), &flow.CreateTaskInput{
			Username: "test",
			Title:    firstTitle,
			Now:      time.Now(),
			ParentID: "",
			NextID:   "",
		})
		_, _ = service.CreateTask(t.Context(), &flow.CreateTaskInput{
			Username: "test",
			Title:    secondTitle,
			Now:      time.Now(),
			ParentID: "",
			NextID:   "",
		})

		ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
			Username: "test",
			ParentID: "",
		})

		require.NoError(t, err)
		require.Len(t, ret.Tasks, 2)
		require.Equal(t, firstTitle, ret.Tasks[0].Title)
		require.Equal(t, secondTitle, ret.Tasks[1].Title)
	})

	t.Run("two-reverse", func(t *testing.T) {
		t.Parallel()

		const (
			firstTitle  = "first"
			secondTitle = "second"
		)

		service, _ := newService(t)
		_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
			Username: "test",
		})
		firstTask, _ := service.CreateTask(t.Context(), &flow.CreateTaskInput{
			Username: "test",
			Title:    firstTitle,
			Now:      time.Now(),
			ParentID: "",
			NextID:   "",
		})
		_, _ = service.CreateTask(t.Context(), &flow.CreateTaskInput{
			Username: "test",
			Title:    secondTitle,
			Now:      time.Now(),
			ParentID: "",
			NextID:   firstTask.ID,
		})

		ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
			Username: "test",
			ParentID: "",
		})

		require.NoError(t, err)
		require.Len(t, ret.Tasks, 2)
		require.Equal(t, secondTitle, ret.Tasks[0].Title)
		require.Equal(t, firstTitle, ret.Tasks[1].Title)
	})
}

func TestServiceDeleteTask(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
		Username: "test",
	})
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
	_ = service.CreateRootDummy(t.Context(), &flow.CreateRootDummyInput{
		Username: "test",
	})

	err := service.DeleteTask(t.Context(), &flow.DeleteTaskInput{
		Username: "test",
		TaskID:   "test",
	})

	require.Error(t, err)
}
