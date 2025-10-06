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

					repo, err := gorm.NewRepository()
					if err != nil {
						return fmt.Errorf("failed to create repository: %w", err)
					}

					service := flow.NewService(ulid.NewIDMaker(), repo)

					mux := goahttp.NewMuxer()
					requestDecoder := goahttp.RequestDecoder
					responseEncoder := goahttp.ResponseEncoder

					handler := NewHandler(service)
					endpoints := task.NewEndpoints(handler)
					taskServer := taskserver.New(endpoints, mux, requestDecoder, responseEncoder, nil, nil)
					taskServer.Mount(mux)

					server := &http.Server{Addr: ":8080", Handler: mux} //nolint:exhaustruct,gosec

					err = server.ListenAndServe()
					if err != nil {
						return fmt.Errorf("failed to listen and serve: %w", err)
					}

					return nil
				},
			},
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
