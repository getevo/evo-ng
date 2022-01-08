package action

import (
	"bytes"
	"encoding/json"
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/cmd/evo/ng"
	"github.com/getevo/evo-ng/lib/file"
	"github.com/getevo/evo-ng/lib/proc"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

func CreateApp() {
	var wd = proc.WorkingDir()
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		proc.Die(err)
	}
	if len(files) > 0 {
		proc.Die("directory is not empty")
	}

	var app = ng.Skeleton{}
	app.App = filepath.Base(wd)
	app.Version = ng.Version{
		Major: 0,
		Minor: 0,
	}

	app.HotReload = true
	app.Debug = true
	app.Include = []string{}

	app.Config = []string{
		"./config.yml",
	}

	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetIndent("", "    ")
	if err := enc.Encode(app); err != nil {
		proc.Die(err)
	}

	err = file.Write(wd+"/app.json", buffer.Bytes())
	if err != nil {
		proc.Die(err)
	}

	var config = evo.Configuration{}
	var b, _ = yaml.Marshal(config.Default())
	err = file.Write(wd+"/config.yml", b)
	if err != nil {
		proc.Die(err)
	}

	proc.Die()
}
