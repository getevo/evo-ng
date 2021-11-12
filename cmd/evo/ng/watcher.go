package ng

import (
	"fmt"
	"github.com/getevo/evo-ng/cmd/evo/watcher"
	"github.com/getevo/evo/lib/gpath"
	"github.com/getevo/evo/lib/log"
	"os"
	"path/filepath"
	"time"
)


type Build struct {
	BinName     string
	BuildArgs   []string
	WorkingDir  string
	ProgramArgs []string
	Run         bool
}

func Watcher()  {
	fmt.Println("Hot Reload Mode")
	cfg := Build{}
	cfg.WorkingDir = gpath.WorkingDir()
	var onBuild = false
	var builder = watcher.NewBuilder(cfg.WorkingDir, cfg.BinName, cfg.WorkingDir, cfg.BuildArgs)
	var runner = watcher.NewRunner(os.Stdout, os.Stderr, filepath.Join(cfg.WorkingDir, builder.Binary()), os.Args[1:])
	watcher.NewWatcher(cfg.WorkingDir, func() {
		if onBuild {
			fmt.Println("skip build due another build")
			return
		}
		runner.Kill()
		var counter = 0
		for runner.IsRunning() {
			counter++
			if counter > 2 {
				fmt.Println("Unable to kill process. try again ...")
				runner.Kill()
				break
			}
		}
		onBuild = true
		err := builder.Build()
		onBuild = false
		if err != nil {

			fmt.Println("\n\nBUILD FAILED:")
			log.Error(err)
		} else {

			onBuild = false
			_, err = runner.Run()
			if err != nil {
				log.Error(err)
			}
		}
	})

	onBuild = true
	err := builder.Build()
	onBuild = false
	if err != nil {
		log.Error(err)
	} else {

		_, err = runner.Run()
		if err != nil {
			fmt.Println("\n\nBUILD FAILED:")
			log.Error(err)
			//fmt.Println("")
		}
	}

	for {
		time.Sleep(1 * time.Minute)
	}
}
