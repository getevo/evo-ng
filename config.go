package evo

import (
	"github.com/creasty/defaults"
	"time"
)

type Config struct {
	WebServer struct {
		Bind                         string   `default:"0.0.0.0" yaml:"bind" json:"bind"`
		Port                         string   `default:"8080" yaml:"port" json:"port"`
		StrictRouting                bool     `default:"false" yaml:"strict_routing" json:"strict_routing"`
		ServerHeader                 string   `default:"evo-ng" yaml:"server_header" json:"server_header"`
		ETag                         bool     `default:"true" yaml:"etag" json:"etag"`
		Recover                      bool     `default:"true" yaml:"recover" json:"recover"`
		BodyLimit                    string   `default:"4mb" yaml:"body_limit" json:"body_limit"`
		ProxyHeader                  string   `default:"" yaml:"proxy_header" json:"proxy_header"`
		DisableKeepalive             bool     `default:"true" yaml:"disable_keepalive" json:"disable_keepalive"`
		DisablePreParseMultipartForm bool     `default:"false" yaml:"disable_pre_parse_multipart_form" json:"disable_pre_parse_multipart_form"`
		StaticDir                    []string `default:"[\"./web\"]" yaml:"static_dir" json:"static_dir"`
		Debug                        bool     `default:"true" json:"debug" yaml:"debug"`
	} `yaml:"web_server" json:"web_server"`

	Database struct {
		Enabled          bool          `yaml:"enabled" json:"enabled"`
		Type             string        `yaml:"type" json:"type"`
		Username         string        `yaml:"user" json:"username"`
		Password         string        `yaml:"pass" json:"password"`
		Server           string        `yaml:"server" json:"server"`
		Cache            string        `yaml:"cache" json:"cache"`
		CacheSize        string        `yaml:"cache_size" json:"cache_size"`
		Debug            string        `yaml:"debug" json:"debug"`
		Database         string        `yaml:"database" json:"database"`
		Schema           string        `yaml:"schema" json:"schema"`
		SSLMode          string        `yaml:"ssl_mode" json:"ssl_mode"`
		Params           string        `yaml:"params" json:"params"`
		MaxOpenConns     int           `yaml:"max_open_connections" json:"max_open_connections"`
		MaxIdleConns     int           `yaml:"max_idle_connections" json:"max_idle_connections"`
		ConnMaxLifeTime  time.Duration `yaml:"connection_max_lifetime" json:"connection_max_lifetime"`
		ConnectionString string        `yaml:"connection_string" json:"connection_string"`
	} `yaml:"database" json:"database"`
}

func (c Config) Default() Config {
	if err := defaults.Set(&c); err != nil {
		panic(err)
	}
	return c
}
