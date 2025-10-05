package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	_ "goa.design/goa/v3/codegen"
	_ "goa.design/goa/v3/codegen/generator"

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
