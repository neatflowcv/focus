package flow_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/idmaker/ulid"
	"github.com/neatflowcv/focus/internal/pkg/repository/memory"
	"github.com/stretchr/testify/require"
)

type ServiceData struct {
	idmaker *ulid.IDMaker
	repo    *memory.Repository
}

func newService(t *testing.T) (*flow.Service, *ServiceData) { //nolint:unparam
	t.Helper()

	data := &ServiceData{
		idmaker: ulid.NewIDMaker(),
		repo:    memory.NewRepository(),
	}

	return flow.NewService(data.idmaker, data.repo), data
}

func TestServiceCreateTask(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	now := time.Now()

	ret, err := service.CreateTask(t.Context(), &flow.CreateTaskInput{
		Username: "test",
		ParentID: "",
		Title:    "test",
	}, now)

	require.NoError(t, err)
	require.Equal(t, "test", ret.Title())
	require.Empty(t, ret.ParentID())
	require.Equal(t, now, ret.CreatedAt())
	require.NotEmpty(t, ret.ID())
	require.Equal(t, domain.TaskStatusTodo, ret.Status())
	require.InEpsilon(t, 10.0, ret.Order(), 0.05)
}

func TestServiceListTasksWithNoData(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
		Username:  "test",
		ParentID:  "",
		Recursive: false,
	})

	require.NoError(t, err)
	require.Empty(t, ret)
}

func TestServiceListTasksWithNoDataRecursive(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
		Username:  "test",
		ParentID:  "",
		Recursive: true,
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
		ParentID: "",
		Title:    "test",
	}, now)

	ret, err := service.ListTasks(t.Context(), &flow.ListTasksInput{
		Username:  "test",
		ParentID:  "",
		Recursive: false,
	})

	require.NoError(t, err)
	require.Len(t, ret, 1)
	require.Equal(t, "test", ret[0].Title())
	require.Empty(t, ret[0].ParentID())
	require.Equal(t, now, ret[0].CreatedAt())
	require.Equal(t, domain.TaskStatusTodo, ret[0].Status())
	require.InEpsilon(t, 10.0, ret[0].Order(), 0.05)
}
