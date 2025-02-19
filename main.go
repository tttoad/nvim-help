package main

import (
	"log"
	"os"

	"nvim-help/internal/action"
)

const (
	RUsege = "command to be executed"
)

func main() {
	e := action.NewExector()
	e.Register(action.NewVersionAction())
	e.Register(action.NewModPath())
	e.Register(action.NewDebugByDocker())
	e.Register(action.NewYamlEdit())
	e.Register(action.NewTagsEdit())
	if len(os.Args) < 2 {
		log.Fatal("need args")
	}
	e.Run(os.Args[1], os.Args[2:])
}
