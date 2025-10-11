package trace_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/focus/internal/app/trace"
	"github.com/neatflowcv/focus/internal/pkg/repository/memory"
	"github.com/stretchr/testify/require"
)

type ServiceData struct {
	repo *memory.Repository
}

func newService(t *testing.T) (*trace.Service, *ServiceData) {
	t.Helper()

	repo := memory.NewRepository()
	service := trace.NewService(repo)

	return service, &ServiceData{
		repo: repo,
	}
}

func TestServiceCreateTrace(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	err := service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "1",
		ParentID: "",
	})

	require.NoError(t, err)
}

func TestServiceDeleteTrace(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "3",
		ParentID: "",
	})
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "2",
		ParentID: "3",
	})
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "1",
		ParentID: "2",
	})
	_ = service.SetActual(t.Context(), &trace.SetActualInput{
		ID:     "2",
		Actual: 5 * time.Second,
	})
	_ = service.SetActual(t.Context(), &trace.SetActualInput{
		ID:     "1",
		Actual: 10 * time.Second,
	})

	err := service.DeleteTrace(t.Context(), &trace.DeleteTraceInput{
		ID: "1",
	})

	require.NoError(t, err)
	require.Len(t, data.repo.Traces, 2)
	require.Equal(t, 5*time.Second, data.repo.Traces["2"].Actual())
	require.Equal(t, 5*time.Second, data.repo.Traces["3"].Actual())
}

func TestServiceDeleteTrace_Error(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	err := service.DeleteTrace(t.Context(), &trace.DeleteTraceInput{
		ID: "1",
	})

	require.ErrorIs(t, err, trace.ErrTraceNotFound)
}

func TestServiceSetActual(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "3",
		ParentID: "",
	})
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "2",
		ParentID: "3",
	})
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "1",
		ParentID: "2",
	})

	err := service.SetActual(t.Context(), &trace.SetActualInput{
		ID:     "1",
		Actual: 10 * time.Second,
	})

	require.NoError(t, err)
	require.Equal(t, 10*time.Second, data.repo.Traces["1"].Actual())
	require.Equal(t, 10*time.Second, data.repo.Traces["2"].Actual())
	require.Equal(t, 10*time.Second, data.repo.Traces["3"].Actual())
}

func TestServiceSetActual_Error(t *testing.T) {
	t.Parallel()

	service, _ := newService(t)

	err := service.SetActual(t.Context(), &trace.SetActualInput{
		ID:     "1",
		Actual: 10 * time.Second,
	})

	require.ErrorIs(t, err, trace.ErrTraceNotFound)
}

func TestService_Actual(t *testing.T) {
	t.Parallel()

	service, data := newService(t)
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "3",
		ParentID: "",
	})
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "2",
		ParentID: "3",
	})
	_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
		ID:       "1",
		ParentID: "2",
	})
	_ = service.SetActual(t.Context(), &trace.SetActualInput{
		ID:     "2",
		Actual: 5 * time.Second,
	})
	_ = service.SetActual(t.Context(), &trace.SetActualInput{
		ID:     "1",
		Actual: 10 * time.Second,
	})

	require.Len(t, data.repo.Traces, 3)
	require.Equal(t, 10*time.Second, data.repo.Traces["1"].Actual())
	require.Equal(t, 15*time.Second, data.repo.Traces["2"].Actual())
	require.Equal(t, 15*time.Second, data.repo.Traces["3"].Actual())
}

func TestServiceListTraces(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(t)

		out, err := service.ListTraces(t.Context(), &trace.ListTracesInput{
			IDs: nil,
		})

		require.NoError(t, err)
		require.Empty(t, out.Traces)
	})

	t.Run("one", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(t)
		_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
			ID:       "1",
			ParentID: "",
		})

		out, err := service.ListTraces(t.Context(), &trace.ListTracesInput{
			IDs: []string{"1"},
		})

		require.NoError(t, err)
		require.Len(t, out.Traces, 1)
	})

	t.Run("two", func(t *testing.T) {
		t.Parallel()

		service, _ := newService(t)
		_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
			ID:       "1",
			ParentID: "",
		})
		_ = service.CreateTrace(t.Context(), &trace.CreateTraceInput{
			ID:       "2",
			ParentID: "",
		})

		out, err := service.ListTraces(t.Context(), &trace.ListTracesInput{
			IDs: []string{"1", "2"},
		})

		require.NoError(t, err)
		require.Len(t, out.Traces, 2)
	})
}
