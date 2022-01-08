package redis

import (
	"fmt"
	"github.com/getevo/evo-ng"
)

var Config struct {
	Redis struct {
		AppID  string   `json:"app_id" yaml:"app_id"`
		NodeID string   `json:"node_id" yaml:"node_id"`
		Server []string `json:"server" yaml:"server"`
	} `json:"redis" yaml:"redis"`
}

func Register() error {
	evo.ParseConfig(&Config)
	fmt.Printf("%+v", Config)
	return nil
}
