package evo

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/getevo/evo-ng/internal/args"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
	"time"
)

type Configuration struct {
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

func (c Configuration) Default() Configuration {
	if err := defaults.Set(&c); err != nil {
		panic(err)
	}
	return c
}

func ParseConfig(params ...interface{}) error {
	var path = ""
	var category = ""
	var out int
	for i, item := range params {
		switch val := item.(type) {
		case string:
			if path == "" {
				path = val
			} else {
				category = val
			}
		default:
			if reflect.TypeOf(val).Kind() == reflect.Struct {
				out = i
			}
		}

	}

	if path == "" {
		if v := args.Get("-c"); v != "" {
			path = v
		} else {
			path = "config.yml"
		}
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not load config file at %s", path)
	}
	m := map[string]interface{}{}
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	var toDecode interface{}
	if category != "" {
		if v, ok := m[category]; ok {
			toDecode = v
		} else {
			return fmt.Errorf("cannot find %s in %s", category, path)
		}
	} else {
		toDecode = m
	}

	cfg := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &params[out],
		TagName:          "yaml",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}
	err = decoder.Decode(toDecode)
	if err != nil {
		return err
	}
	return nil

}
