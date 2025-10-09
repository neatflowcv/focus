package extra_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/focus/internal/app/extra"
	"github.com/neatflowcv/focus/internal/pkg/domain"
	"github.com/neatflowcv/focus/internal/pkg/repository/memory"
	"github.com/stretchr/testify/require"
)

type ServiceData struct {
	repo *memory.Repository
}

func newService(t *testing.T) (*extra.Service, *ServiceData) {
	t.Helper()

	repo := memory.NewRepository()
	service := extra.NewService(repo)

	return service, &ServiceData{
		repo: repo,
	}
}

func TestServiceCreateExtra(t *testing.T) {
	t.Parallel()

	service, data := newService(t)

	_, err := service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "test",
		ParentID: "",
	})

	require.NoError(t, err)
	require.Len(t, data.repo.Extras, 1)
	require.Equal(t, time.Duration(0), data.repo.Extras["test"].ActualTime())
	require.Equal(t, time.Duration(0), data.repo.Extras["test"].EstimatedTime())
	require.Equal(t, time.Time{}, data.repo.Extras["test"].StartedAt())
	require.True(t, data.repo.Extras["test"].Leaf())
	require.Equal(t, domain.TaskStatusTodo, data.repo.Extras["test"].Status())
}

func TestServiceDeleteExtra(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "test",
		ParentID: "",
	})

	err := service.DeleteExtra(t.Context(), &extra.DeleteExtraInput{
		ID: "test",
	})

	require.NoError(t, err)
}

func TestServiceDeleteExtra_Error(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	err := service.DeleteExtra(t.Context(), &extra.DeleteExtraInput{
		ID: "test",
	})

	require.ErrorIs(t, err, extra.ErrExtraNotFound)
}

func TestServiceUpdateEstimatedTime(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "parent",
		ParentID: "",
	})
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "child",
		ParentID: "parent",
	})

	err := service.UpdateEstimatedTime(t.Context(), &extra.UpdateEstimatedTimeInput{
		ID:            "child",
		EstimatedTime: time.Hour,
	})

	require.NoError(t, err)
	require.Equal(t, time.Hour, data.repo.Extras["child"].EstimatedTime())
	require.Equal(t, time.Duration(0), data.repo.Extras["parent"].EstimatedTime())
}

func TestServiceUpdateActualTime(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "parent",
		ParentID: "",
	})
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "child",
		ParentID: "parent",
	})

	err := service.UpdateActualTime(t.Context(), &extra.UpdateActualTimeInput{
		ID:         "child",
		ActualTime: time.Hour,
	})

	require.NoError(t, err)
	require.Equal(t, time.Hour, data.repo.Extras["child"].AccActualTime())
	require.Equal(t, time.Hour, data.repo.Extras["parent"].AccActualTime())
}

func TestServiceListExtras0(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	ret, err := service.ListExtras(t.Context(), &extra.ListExtrasInput{
		IDs: []string{"test"},
	})

	require.NoError(t, err)
	require.Empty(t, ret.Extras)
}

func TestServiceListExtras1(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "test",
		ParentID: "",
	})

	ret, err := service.ListExtras(t.Context(), &extra.ListExtrasInput{
		IDs: []string{"test"},
	})

	require.NoError(t, err)
	require.Len(t, ret.Extras, 1)
}

func TestServiceListExtras2(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "test1",
		ParentID: "",
	})
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "test2",
		ParentID: "",
	})

	ret, err := service.ListExtras(t.Context(), &extra.ListExtrasInput{
		IDs: []string{"test1", "test2"},
	})

	require.NoError(t, err)
	require.Len(t, ret.Extras, 2)
}

func TestServiceCheckLeaf(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "parent",
		ParentID: "",
	})
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "child",
		ParentID: "parent",
	})

	require.True(t, data.repo.Extras["child"].Leaf())
	require.False(t, data.repo.Extras["parent"].Leaf())
}

func TestServiceCheckStatus1(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "parent",
		ParentID: "",
	})
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "child",
		ParentID: "parent",
	})

	require.Equal(t, domain.TaskStatusTodo, data.repo.Extras["parent"].Status())
	require.Equal(t, domain.TaskStatusTodo, data.repo.Extras["child"].Status())
}

func TestServiceCheckStatus2(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "parent",
		ParentID: "",
	})
	_ = service.SetDone(t.Context(), &extra.SetDoneInput{
		ID: "parent",
	})
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "child",
		ParentID: "parent",
	})

	require.Equal(t, domain.TaskStatusTodo, data.repo.Extras["parent"].Status())
	require.Equal(t, domain.TaskStatusTodo, data.repo.Extras["child"].Status())
}

func TestServiceCheckStatus3(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "parent",
		ParentID: "",
	})
	_ = service.SetDone(t.Context(), &extra.SetDoneInput{
		ID: "parent",
	})
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID:       "child",
		ParentID: "",
	})
	_ = service.UpdateParent(t.Context(), &extra.UpdateParentInput{
		ID:       "child",
		ParentID: "parent",
	})

	require.Equal(t, domain.TaskStatusTodo, data.repo.Extras["parent"].Status())
	require.Equal(t, domain.TaskStatusTodo, data.repo.Extras["child"].Status())
}
