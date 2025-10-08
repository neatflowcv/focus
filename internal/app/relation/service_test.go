package relation_test

import (
	"testing"

	"github.com/neatflowcv/focus/internal/app/relation"
	"github.com/neatflowcv/focus/internal/pkg/eventbus"
	"github.com/neatflowcv/focus/internal/pkg/repository/memory"
	"github.com/stretchr/testify/require"
)

type ServiceData struct {
	bus  *eventbus.Bus
	repo *memory.Repository
}

func newService(t *testing.T) (*relation.Service, *ServiceData) {
	t.Helper()

	data := &ServiceData{
		bus:  eventbus.NewBus(),
		repo: memory.NewRepository(),
	}

	return relation.NewService(data.bus, data.repo), data
}

func TestServiceCreateChildDummy(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	err := service.CreateChildDummy(t.Context(), &relation.CreateChildDummyInput{
		ID: "1",
	})

	require.NoError(t, err)
}

func TestServiceDeleteChildDummy(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	_ = service.CreateChildDummy(t.Context(), &relation.CreateChildDummyInput{
		ID: "1",
	})

	err := service.DeleteChildDummy(t.Context(), &relation.DeleteChildDummyInput{
		ID: "1",
	})

	require.NoError(t, err)
}

func TestServiceCreateRelation(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_ = service.CreateChildDummy(t.Context(), &relation.CreateChildDummyInput{
		ID: "1",
	})

	err := service.CreateRelation(t.Context(), &relation.CreateRelationInput{
		ID:       "2",
		ParentID: "1",
	})

	require.NoError(t, err)
	require.Len(t, data.repo.Relations, 2)
	require.Equal(t, "2", string(data.repo.Relations["2"].ID()))
	require.Empty(t, string(data.repo.Relations["2"].NextID()))
	require.Equal(t, "1", string(data.repo.Relations["2"].ParentID()))
}

func TestServiceListChildren(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	_ = service.CreateChildDummy(t.Context(), &relation.CreateChildDummyInput{
		ID: "1",
	})

	ret, err := service.ListChildren(t.Context(), &relation.ListChildrenInput{
		ParentID: "1",
	})

	require.NoError(t, err)
	require.Empty(t, ret.IDs)
}

func TestServiceListChildrenWithData(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	_ = service.CreateChildDummy(t.Context(), &relation.CreateChildDummyInput{
		ID: "0",
	})
	_ = service.CreateRelation(t.Context(), &relation.CreateRelationInput{
		ID:       "1",
		ParentID: "0",
	})

	ret, err := service.ListChildren(t.Context(), &relation.ListChildrenInput{
		ParentID: "0",
	})

	require.NoError(t, err)
	require.Equal(t, []string{"1"}, ret.IDs)
}

func TestServiceListChildrenWithData2(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)
	_ = service.CreateChildDummy(t.Context(), &relation.CreateChildDummyInput{
		ID: "",
	})
	_ = service.CreateRelation(t.Context(), &relation.CreateRelationInput{
		ID:       "1",
		ParentID: "",
	})
	_ = service.CreateRelation(t.Context(), &relation.CreateRelationInput{
		ID:       "2",
		ParentID: "",
	})

	ret, err := service.ListChildren(t.Context(), &relation.ListChildrenInput{
		ParentID: "",
	})

	require.NoError(t, err)
	require.Equal(t, []string{"2", "1"}, ret.IDs)
}
