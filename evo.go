package evo

import (
	"fmt"
	"github.com/getevo/evo-ng/internal/file"
	"github.com/gofiber/fiber/v2"
	"time"
)

type ContextInterface interface {
	New() ContextInterface
	Get(*Context, interface{}) error
	Post(*Context, interface{}) error
	Put(*Context, interface{}) error
	Patch(*Context, interface{}) error
	All(*Context, interface{}) error
	Delete(*Context, interface{}) error
	Options(*Context, interface{}) error
	Head(*Context, interface{}) error
	Connect(*Context, interface{}) error
	Use(*Context, interface{}) error
}

var app = fiber.New()
var contextInstance ContextInterface

type Version struct {
	Major int       `json:"major"`
	Minor int       `json:"minor"`
	Build int       `json:"build"`
	Hash  string    `json:"hash"`
	Date  time.Time `json:"date"`
}

type Context struct {
	fiber *fiber.Ctx
}

type AssetConfig fiber.Static
type Handler func(*fiber.Ctx) error

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
	app.Post(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func Patch(url string, callback interface{}, params ...interface{}) {
	app.Patch(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func Put(url string, callback interface{}, params ...interface{}) {
	app.Put(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func All(url string, callback interface{}, params ...interface{}) {
	app.All(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func Delete(url string, callback interface{}, params ...interface{}) {
	app.Delete(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func Options(url string, callback interface{}, params ...interface{}) {
	app.Options(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func Connect(url string, callback interface{}, params ...interface{}) {
	app.Connect(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func Head(url string, callback interface{}, params ...interface{}) {
	app.Head(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Get(&ct, callback)
	})
}

func Use(url string, callback interface{}, params ...interface{}) {
	app.Use(url, func(ctx *fiber.Ctx) error {
		var ct = Context{
			fiber: ctx,
		}
		return contextInstance.New().Use(&ct, callback)
	})
}

// Asset set static file path.
// first arg takes url to match
// second arg takes local path
// returns error if the localPath does not exists
func Asset(url, localPath string, config ...AssetConfig) error {
	if !file.IsDir(localPath) {
		return fmt.Errorf("static file path %s does not exists", localPath)
	}
	if len(config) > 0 {
		app.Static(url, localPath, fiber.Static(config[0]))
	} else {
		app.Static(url, localPath)
	}
	app.Use(url, func(c *fiber.Ctx) error {
		c.Status(fiber.StatusNotFound)
		return fmt.Errorf("requested URL %s was not found on this server", c.Path())
	})
	return nil
}

func Run() {
	app.Listen(":80")
}
