package install

import (
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/lib/args"
	"github.com/getevo/evo-ng/lib/file"
	"github.com/getevo/evo-ng/lib/proc"
	"log"
	"os"
	"runtime"
)

const (
	WINDOWS = "windows"
	LINUX   = "linux"
	DARWIN  = "darwin"
)

func Install() {
	switch runtime.GOOS {
	case WINDOWS:
		windows()
	case LINUX:
		linux()
	case DARWIN:
		darwin()
	}

}

func windows() {
	var dir = os.Getenv("windir") + "/System32/evo-ng.exe"
	if args.Exists("-update") || !file.IsFileExist(dir) {
		err := file.CopyFile(os.Args[0], dir)
		if err != nil {
			log.Panicf(err.Error())
		}
		proc.Die("EVO-NG has successfully installed")
	}
}

func darwin() {
	evo.Panic("macos is not supported")
}

func linux() {
	var dir = "/usr/bin/evo-ng"
	if args.Exists("-update") || !file.IsFileExist(dir) {
		err := file.CopyFile(os.Args[0], dir)
		if err != nil {
			log.Panicf(err.Error())
		}
		proc.Die("EVO-NG has successfully installed")
	}
}
