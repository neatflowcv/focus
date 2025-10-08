package extra_test

import (
	"errors"
	"testing"
	"time"

	"github.com/neatflowcv/focus/internal/app/extra"
	"github.com/neatflowcv/focus/internal/app/relation"
	"github.com/neatflowcv/focus/internal/pkg/eventbus"
	"github.com/neatflowcv/focus/internal/pkg/repository/memory"
	"github.com/stretchr/testify/require"
)

type ServiceData struct {
	repo            *memory.Repository
	relationService *relation.Service
	extraService    *extra.Service
}

func newService(t *testing.T) (*extra.Service, *ServiceData) {
	t.Helper()

	bus := eventbus.NewBus()
	repo := memory.NewRepository()
	relationService := relation.NewService(bus, repo)
	service := extra.NewService(repo, relationService)

	return service, &ServiceData{
		repo:            repo,
		relationService: relationService,
		extraService:    service,
	}
}

func TestServiceCreateExtra(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_, _ = service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID: "",
	})
	_ = data.relationService.CreateChildDummy(t.Context(), &relation.CreateChildDummyInput{
		ID: "",
	})
	_ = data.relationService.CreateRelation(t.Context(), &relation.CreateRelationInput{
		ID:       "test",
		ParentID: "",
	})

	ret, err := service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID: "test",
	})

	require.NoError(t, err)
	require.Equal(t, time.Duration(0), ret.Extra.ActualTime)
	require.Equal(t, time.Duration(0), ret.Extra.EstimatedTime)
	require.Equal(t, time.Time{}, ret.Extra.StartedAt)
	require.True(t, ret.Extra.Leaf)
}

func TestServiceCreateExtra_Error(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	ret, err := service.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID: "test",
	})
	_ = ret

	require.ErrorIs(t, err, extra.ErrPreconditionFailed)
}

func TestServiceDeleteExtra(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_ = data

	err := service.DeleteExtra(t.Context(), &extra.DeleteExtraInput{
		ID: "test",
	})

	require.NoError(t, err)
}

func TestServiceUpdateEstimatedTime(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	declareRelation(t, data, "test", "")

	err := service.UpdateEstimatedTime(t.Context(), &extra.UpdateEstimatedTimeInput{
		ID:            "test",
		EstimatedTime: time.Hour,
	})

	require.NoError(t, err)
	require.Equal(t, time.Hour, data.repo.Extras["test"].EstimatedTime())
}

func TestServiceUpdateActualTime(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	declareRelation(t, data, "parent", "")
	declareRelation(t, data, "child", "parent")

	err := service.UpdateActualTime(t.Context(), &extra.UpdateActualTimeInput{
		ID:         "child",
		ActualTime: time.Hour,
	})

	require.NoError(t, err)
	require.Equal(t, time.Hour, data.repo.Extras["child"].ActualTime())
	require.Equal(t, time.Hour, data.repo.Extras["parent"].ActualTime())
}

func TestServiceListExtras(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_ = data

	ret, err := service.ListExtras(t.Context(), &extra.ListExtrasInput{
		IDs: []string{"test"},
	})

	require.NoError(t, err)
	require.Empty(t, ret.Extras)
}

func TestServiceCheckLeaf(t *testing.T) {
	t.Parallel()

	_, data := newService(t)

	declareRelation(t, data, "parent", "")
	declareRelation(t, data, "child", "parent")

	require.True(t, data.repo.Extras["child"].Leaf())
	require.False(t, data.repo.Extras["parent"].Leaf())
}

func declareRelation(t *testing.T, data *ServiceData, id string, parentID string) {
	t.Helper()

	err := data.relationService.CreateChildDummy(t.Context(), &relation.CreateChildDummyInput{
		ID: parentID,
	})
	if err != nil && !errors.Is(err, relation.ErrRelationAlreadyExists) {
		require.NoError(t, err)
	}

	err = data.relationService.CreateRelation(t.Context(), &relation.CreateRelationInput{
		ID:       id,
		ParentID: parentID,
	})
	require.NoError(t, err)

	_, _ = data.extraService.CreateExtra(t.Context(), &extra.CreateExtraInput{
		ID: id,
	})
}
