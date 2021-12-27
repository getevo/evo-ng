package geoip

import (
	"github.com/getevo/evo/lib/text"
)

//Context Extend GeoIP Context
type Context struct {
	Test string `json:"test"`
}

func (context *Context) Extend() error {
	context.Test = text.Random(10)
	return nil
}

func Register() error {

	return nil
}
