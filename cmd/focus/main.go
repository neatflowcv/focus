package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	_ "goa.design/goa/v3/codegen"
	_ "goa.design/goa/v3/codegen/generator"
	goahttp "goa.design/goa/v3/http"

	taskserver "github.com/neatflowcv/focus/gen/http/task/server"
	"github.com/neatflowcv/focus/gen/task"
	"github.com/neatflowcv/focus/internal/app/extra"
	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/focus/internal/app/trace"
	"github.com/neatflowcv/focus/internal/pkg/eventbus"
	"github.com/neatflowcv/focus/internal/pkg/idmaker/ulid"
	"github.com/neatflowcv/focus/internal/pkg/repository/gorm"
	"github.com/urfave/cli/v3"
)

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return info.Main.Version
}

func main() {
	log.Println("version", version())

	app := &cli.Command{ //nolint:exhaustruct
		Name: "focus",
		Commands: []*cli.Command{
			{
				Name: "run",
				Action: func(ctx context.Context, c *cli.Command) error {
					log.Println("running")

					return run()
				},
			},
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run() error { //nolint:cyclop,funlen
	repo, err := gorm.NewRepository()
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	bus := eventbus.NewBus()

	flowService := flow.NewService(bus, ulid.NewIDMaker(), repo)
	extraService := extra.NewService(bus, repo)
	traceService := trace.NewService(repo)

	server := newServer(flowService, extraService, traceService)

	bus.TaskCreated.Subscribe(func(ctx context.Context, event *eventbus.TaskCreatedEvent) {
		err := extraService.CreateExtra(ctx, &extra.CreateExtraInput{
			ID:       event.TaskID,
			ParentID: event.ParentID,
		})
		if err != nil {
			log.Printf("failed to create extra: %v", err)
		}
	})
	bus.TaskDeleted.Subscribe(func(ctx context.Context, event *eventbus.TaskDeletedEvent) {
		err := extraService.DeleteExtra(ctx, &extra.DeleteExtraInput{
			ID: event.TaskID,
		})
		if err != nil {
			log.Printf("failed to delete extra: %v", err)
		}
	})
	bus.TaskRelationUpdated.Subscribe(func(ctx context.Context, event *eventbus.TaskRelationUpdatedEvent) {
		if event.OldParentID == event.NewParentID {
			return
		}

		err := extraService.UpdateParent(ctx, &extra.UpdateParentInput{
			ID:       event.TaskID,
			ParentID: event.NewParentID,
		})
		if err != nil {
			log.Printf("failed to update parent extra: %v", err)
		}
	})

	bus.TaskCreated.Subscribe(func(ctx context.Context, event *eventbus.TaskCreatedEvent) {
		err := traceService.CreateTrace(ctx, &trace.CreateTraceInput{
			ID:       event.TaskID,
			ParentID: event.ParentID,
		})
		if err != nil {
			log.Printf("failed to create trace: %v", err)
		}
	})
	bus.TaskDeleted.Subscribe(func(ctx context.Context, event *eventbus.TaskDeletedEvent) {
		err := traceService.DeleteTrace(ctx, &trace.DeleteTraceInput{
			ID: event.TaskID,
		})
		if err != nil {
			log.Printf("failed to delete trace: %v", err)
		}
	})
	bus.TaskRelationUpdated.Subscribe(func(ctx context.Context, event *eventbus.TaskRelationUpdatedEvent) {
		if event.OldParentID == event.NewParentID {
			return
		}

		err := traceService.UpdateParent(ctx, &trace.UpdateParentInput{
			ID:       event.TaskID,
			ParentID: event.NewParentID,
		})
		if err != nil {
			log.Printf("failed to update parent trace: %v", err)
		}
	})

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}

func newServer(flowService *flow.Service, extraService *extra.Service, traceService *trace.Service) *http.Server {
	mux := goahttp.NewMuxer()
	requestDecoder := goahttp.RequestDecoder
	responseEncoder := goahttp.ResponseEncoder

	handler := NewHandler(flowService, extraService, traceService)
	endpoints := task.NewEndpoints(handler)
	taskServer := taskserver.New(endpoints, mux, requestDecoder, responseEncoder, nil, nil)
	taskServer.Mount(mux)

	return &http.Server{ //nolint:exhaustruct,gosec
		Addr:    ":8080",
		Handler: mux,
	}
}
