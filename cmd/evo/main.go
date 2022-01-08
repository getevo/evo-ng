package main

import (
	"github.com/getevo/evo-ng/cmd/evo/action"
	"github.com/getevo/evo-ng/cmd/evo/ng"
	"github.com/getevo/evo-ng/install"
	"github.com/getevo/evo-ng/lib/args"
)

const (
	WINDOWS = "windows"
	LINUX   = "linux"
	DARWIN  = "darwin"
)

func main() {
	//check if is installed
	install.Install()

	switch args.Get("create") {
	case "app":
		action.CreateApp()
	case "":
	default:
		panic("unrecognized parameter")
	}

	ng.Start()
}
