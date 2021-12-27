package evo

import (
	"github.com/gofiber/fiber/v2"
)

type ContextInterface interface {
	New() ContextInterface
	Get(*Context, interface{}) error
	Post(*Context, interface{}) error
	Put(*Context, interface{}) error
	Push(*Context, interface{}) error
	All(*Context, interface{}) error
	Delete(*Context, interface{}) error
	Options(*Context, interface{}) error
}

var app = fiber.New()
var contextInstance ContextInterface

type Context struct {
	fiber *fiber.Ctx
}

func Engine() {

}

func UseContext(context ContextInterface) {
	contextInstance = context
}

func Get(url string, callback interface{}, params ...interface{}) {
	app.Get(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func Post(url string, callback interface{}, params ...interface{}) {

}

func Push(url string, callback interface{}, params ...interface{}) {

}

func Put(url string, callback interface{}, params ...interface{}) {

}

func All(url string, callback interface{}, params ...interface{}) {

}

func Delete(url string, callback interface{}, params ...interface{}) {

}

func Options(url string, callback interface{}, params ...interface{}) {

}

func Run() {
	app.Listen(":80")
}
