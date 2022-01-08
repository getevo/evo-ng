package ng

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/getevo/evo-ng/cmd/evo/watcher"
	"github.com/getevo/evo-ng/lib/file"
	"github.com/getevo/evo-ng/lib/proc"
	"github.com/getevo/evo/lib/log"
	"github.com/tidwall/sjson"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Build struct {
	BinName     string
	BuildArgs   []string
	WorkingDir  string
	ProgramArgs []string
	Run         bool
}

var BuildCfg Build

func Watcher() {
	fmt.Println("Hot Reload Mode")
	BuildCfg = Build{}
	BuildCfg.WorkingDir = proc.WorkingDir()
	BuildCfg.BinName = filepath.Base(BuildCfg.WorkingDir)
	if runtime.GOOS == "windows" {
		// check if it already has the .exe extension
		if !strings.HasSuffix(BuildCfg.BinName, ".exe") {
			BuildCfg.BinName += ".exe"
		}
	}

	var onBuild = false
	var builder = watcher.NewBuilder(BuildCfg.WorkingDir, BuildCfg.BinName, BuildCfg.WorkingDir, BuildCfg.BuildArgs)
	var runner = watcher.NewRunner(os.Stdout, os.Stderr, filepath.Join(BuildCfg.WorkingDir, builder.Binary()), os.Args[1:])
	watcher.NewWatcher(BuildCfg.WorkingDir, func() {
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
			updateVersion()
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
		updateVersion()
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

func updateVersion() {
	f, err := os.Open(BuildCfg.WorkingDir + "/" + BuildCfg.BinName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		log.Fatal(err)
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	if skeleton.Version.Hash != hash {
		skeleton.Version.Hash = hash
		skeleton.Version.Build++
		b, err := file.ReadFile("./app.json")
		if err != nil {
			log.Fatal(err)
		}

		b, err = sjson.SetBytes(b, "version.hash", hash)
		if err != nil {
			log.Fatal(err)
		}
		b, err = sjson.SetBytes(b, "version.build", skeleton.Version.Build)
		if err != nil {
			log.Fatal(err)
		}
		b, err = sjson.SetBytes(b, "version.date", time.Now())
		if err != nil {
			log.Fatal(err)
		}
		file.Write(BuildCfg.WorkingDir+"/app.json", b)
	}
}
