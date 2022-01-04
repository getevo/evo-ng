package geoip

import (
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo/lib/text"
)

//Context Extend GeoIP Context
type Context struct {
	Test string `json:"test"`
}

func (context *Context) Extend(request *evo.Context) error {
	context.Test = text.Random(10)
	//fmt.Println("TEST")
	return nil
}

func Register() error {

	return nil
}
