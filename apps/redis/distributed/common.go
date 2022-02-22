package distributed

import (
	"github.com/getevo/evo"
	"github.com/getevo/evo-ng/lib/network"
	"github.com/getevo/evo-ng/lib/proc"
	"os"
)

var appName = ""

func GetAppName() string {
	if appName != "" {
		return appName
	}

	hostname, _ := os.Hostname()
	appName = proc.Name() + ":" + evo.GetConfig().Server.Port + "#" + hostname
	if config, err := network.GetConfig(); err != nil {
		appName += "#" + config.LocalIP.String()
	}

	return appName
}
