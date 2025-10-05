package main

import (
	"log"
	"runtime/debug"

	_ "goa.design/goa/v3/codegen"
	_ "goa.design/goa/v3/codegen/generator"
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
}
