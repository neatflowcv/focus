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
	"github.com/neatflowcv/focus/internal/app/flow"
	"github.com/neatflowcv/focus/internal/app/relation"
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

					return run(ctx)
				},
			},
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	repo, err := gorm.NewRepository()
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	bus := eventbus.NewBus()

	flowService := flow.NewService(bus, ulid.NewIDMaker(), repo)
	relationService := relation.NewService(bus, repo)

	_ = relationService.CreateChildDummy(ctx, &relation.CreateChildDummyInput{
		ID: "",
	})

	bus.TaskCreated.Subscribe(ctx, func(ctx context.Context, event *eventbus.TaskCreatedEvent) {
		err := relationService.CreateChildDummy(ctx, &relation.CreateChildDummyInput{
			ID: event.TaskID,
		})
		if err != nil {
			log.Printf("failed to create child dummy: %v", err)
		}
	})
	bus.TaskDeleted.Subscribe(ctx, func(ctx context.Context, event *eventbus.TaskDeletedEvent) {
		err := relationService.DeleteChildDummy(ctx, &relation.DeleteChildDummyInput{
			ID: event.TaskID,
		})
		if err != nil {
			log.Printf("failed to delete child dummy: %v", err)
		}
	})

	server := newServer(flowService, relationService)

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}

func newServer(service *flow.Service, relationService *relation.Service) *http.Server {
	mux := goahttp.NewMuxer()
	requestDecoder := goahttp.RequestDecoder
	responseEncoder := goahttp.ResponseEncoder

	handler := NewHandler(service, relationService)
	endpoints := task.NewEndpoints(handler)
	taskServer := taskserver.New(endpoints, mux, requestDecoder, responseEncoder, nil, nil)
	taskServer.Mount(mux)

	return &http.Server{ //nolint:exhaustruct,gosec
		Addr:    ":8080",
		Handler: mux,
	}
}
