package evo

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/getevo/evo-ng/lib/args"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"reflect"
	"time"
)

type Configuration struct {
	WebServer struct {
		Bind                         string        `default:"0.0.0.0" config:"bind" json:"bind"`
		Port                         string        `default:"8080" config:"port" json:"port"`
		StrictRouting                bool          `default:"false" config:"strict_routing" json:"strict_routing"`
		ServerHeader                 string        `default:"evo-ng" config:"server_header" json:"server_header"`
		ETag                         bool          `default:"true" config:"etag" json:"etag"`
		Immutable                    bool          `default:"false" config:"immutable" json:"immutable"`
		UnescapePath                 bool          `default:"false" config:"unescape_path" json:"unescape_path"`
		Recover                      bool          `default:"true" config:"recover" json:"recover"`
		BodyLimit                    string        `default:"4mb" config:"body_limit" json:"body_limit"`
		ProxyHeader                  string        `default:"x-forwarded-for" config:"proxy_header" json:"proxy_header"`
		DisableKeepalive             bool          `default:"true" config:"disable_keepalive" json:"disable_keepalive"`
		DisablePreParseMultipartForm bool          `default:"false" config:"disable_pre_parse_multipart_form" json:"disable_pre_parse_multipart_form"`
		StaticDir                    []string      `default:"[\"./web\"]" config:"static_dir" json:"static_dir"`
		ReadTimeout                  time.Duration `default:"0" config:"read_timeout" json:"read_timeout"`
		WriteTimeout                 time.Duration `default:"0" config:"write_timeout" json:"write_timeout"`
		IdleTimeout                  time.Duration `default:"0" config:"idle_timeout" json:"idle_timeout"`
		ReadBufferSize               int           `default:"4096" config:"read_buffer_size" json:"read_buffer_size"`
		WriteBufferSize              int           `default:"4096" config:"write_buffer_size" json:"write_buffer_size"`
		CORS                         bool          `default:"false" config:"cors_control" json:"cors_control"`
		AllowOrigins                 string        `default:"" config:"allow_origins" json:"allow_origins"`
		AllowHeaders                 string        `default:"" config:"allow_headers" json:"allow_headers"`
		AllowMethods                 string        `default:"" config:"allow_methods" json:"allow_methods"`
		AllowCredentials             bool          `default:"false" config:"allow_credentials" json:"allow_credentials"`
		ExposeHeaders                string        `default:"" config:"expose_headers" json:"expose_headers"`
		PreflightMaxCacheAge         int           `default:"0" config:"preflight_max_cache_age" json:"preflight_max_cache_age"`
		Debug                        bool          `default:"true" json:"debug" config:"debug"`
		CompressLevel                int           `default:"-1" json:"compress_level" config:"compress_level"`
		WebSocket                    bool          `default:"false" json:"websocket" config:"websocket"`
	} `config:"web_server" json:"web_server"`

	Database struct {
		Enabled                                  bool     `default:"true" config:"enabled" json:"enabled"`
		Dialect                                  string   `default:"mysql" config:"dialect" json:"dialect"`
		DSN                                      string   `default:"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local" config:"dsn" json:"dsn"`
		Database                                 string   `default:"dbname" config:"database" json:"database"`
		Replicas                                 []string `config:"replicas" json:"replicas"`
		SkipDefaultTransaction                   bool     `default:"false" config:"skip_default_transaction" json:"skip_default_transaction"`
		FullSaveAssociations                     bool     `default:"true" config:"full_save_associations" json:"full_save_associations"`
		DisableAutomaticPing                     bool     `default:"false" config:"disable_automatic_ping" json:"disable_automatic_ping"`
		DisableForeignKeyConstraintWhenMigrating bool     `default:"false" config:"disable_foreign_key_constraint_when_migrating" json:"disable_foreign_key_constraint_when_migrating"`
		DisableNestedTransaction                 bool     `default:"false" config:"disable_nested_transaction" json:"disable_nested_transaction"`
		CreateBatchSize                          int      `default:"100" config:"create_batch_size" json:"create_batch_size"`
		QueryFields                              bool     `default:"true" config:"query_fields" json:"query_fields"`
		StmtCache                                bool     `default:"true" config:"stmt_cache" json:"stmt_cache"`
		TablePrefix                              string   `default:"" config:"table_prefix" json:"table_prefix"`
		Cache                                    string   `config:"cache" json:"cache"`
		CacheSize                                string   `config:"cache_size" json:"cache_size"`
		Debug                                    string   `config:"debug" json:"debug"`
		MaxOpenConns                             int      `default:"100" config:"max_open_connections" json:"max_open_connections"`
		MaxIdleConns                             int      `default:"10" config:"max_idle_connections" json:"max_idle_connections"`
		ConnMaxLifeTime                          string   `default:"1h" config:"connection_max_lifetime" json:"connection_max_lifetime"`
		ConnMaxIdleTime                          string   `default:"10m" config:"connection_max_idle_time" json:"connection_max_idle_time"`
		L10n                                     bool     `default:"false" config:"l10n" json:"l10n"`
	} `config:"database" json:"database"`
}

func (c Configuration) Default() Configuration {
	if err := defaults.Set(&c); err != nil {
		Panic(err)
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
		TagName:          "json",
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
