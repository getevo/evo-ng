package evo

import (
	"fmt"
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	os.Args = append(os.Args, "-etcd", "10.12.47.89:2379", "-etcd-prefix", "/db")
	fmt.Println(os.Args)
	var cfg = struct {
		TTL                  int `config:"ttl"`
		MaximumLagOnFailover int `config:"maximum_lag_on_failover"`
		RetryTimeout         int `config:"retry_timeout"`
		Postgresql           struct {
			UsePgRewind bool `config:"use_pg_rewind"`
		} `config:"postgresql"`
		LoopWait int `config:"loop_wait"`
	}{}

	var err = ParseConfig(&cfg)
	fmt.Println(cfg)
	fmt.Println(err)
}
