package evo

import (
	"fmt"
	"github.com/getevo/evo-ng/internal/file"
	"github.com/getevo/evo-ng/internal/generic"
	websocket2 "github.com/getevo/evo-ng/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
	"time"
)

type ContextInterface interface {
	New() ContextInterface
	Get(*Context, interface{}) error
	WebSocket(*Context, interface{}, *websocket2.Conn) error
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

var app *fiber.App
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

var Config = Configuration{}

func Engine() {
	var err = ParseConfig(&Config)
	if err != nil {
		panic(err)
	}
	Events.Register()
	Events.On("exit", func() {
		go func() {
			var j = 7
			for {
				time.Sleep(1 * time.Second)
				if j == 5 {
					fmt.Println("Something prevents the app from graceful shutdown")
				}
				if j < 6 {
					fmt.Printf("Force shutdown after %d seconds \n", j)
				}
				j--
				if j == 0 {
					os.Exit(1)
				}
			}
		}()

	})
	var config = fiber.Config{
		CaseSensitive:         true,
		StrictRouting:         false,
		Immutable:             Config.WebServer.Immutable,
		ServerHeader:          Config.WebServer.ServerHeader,
		UnescapePath:          Config.WebServer.UnescapePath,
		BodyLimit:             int(generic.Parse(Config.WebServer.BodyLimit).SizeInBytes()),
		DisableStartupMessage: true,
		DisableKeepalive:      Config.WebServer.DisableKeepalive,
		ReadTimeout:           Config.WebServer.ReadTimeout,
		WriteTimeout:          Config.WebServer.WriteTimeout,
		IdleTimeout:           Config.WebServer.IdleTimeout,
		ReadBufferSize:        Config.WebServer.ReadBufferSize,
		WriteBufferSize:       Config.WebServer.WriteBufferSize,
	}
	if config.ReadTimeout < 1 {
		config.ReadTimeout = 0
	}
	if config.WriteTimeout < 1 {
		config.WriteTimeout = 0
	}
	if config.IdleTimeout < 1 {
		config.IdleTimeout = 0
	}
	if config.ReadBufferSize < 1 {
		config.ReadBufferSize = 4096
	}
	if config.WriteBufferSize < 1 {
		config.WriteBufferSize = 4096
	}
	if config.BodyLimit < 1 {
		config.WriteBufferSize = 4 * 1024 * 1024
	}
	app = fiber.New(config)
	if config.ETag {
		app.Use(etag.New())
	}
	if Config.WebServer.Debug {
		app.Use(logger.New())
	}

	if Config.WebServer.AllowOrigins != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     Config.WebServer.AllowOrigins,
			AllowMethods:     Config.WebServer.AllowMethods,
			AllowHeaders:     Config.WebServer.AllowHeaders,
			AllowCredentials: Config.WebServer.AllowCredentials,
			ExposeHeaders:    Config.WebServer.ExposeHeaders,
			MaxAge:           Config.WebServer.PreflightMaxCacheAge,
		}))
	}

	if Config.WebServer.CompressLevel > -1 {
		if Config.WebServer.CompressLevel > 2 {
			Config.WebServer.CompressLevel = 2
		}
		app.Use(compress.New(compress.Config{
			Level: compress.Level(Config.WebServer.CompressLevel),
		}))
	}

	var recoverConfig = recover.Config{
		EnableStackTrace: Config.WebServer.Debug,
	}
	app.Use(recover.New(recoverConfig))

	/*	app.Use("/webs",func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})

		app.Get("/webs",websocket.New(func(c *websocket.Conn) {
			fmt.Println("here")
			// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
			var (
				mt  int
				msg []byte
				err error
			)
			for {
				if mt, msg, err = c.ReadMessage(); err != nil {
					log.Println("read:", err)
					break
				}
				fmt.Printf("recv: %s \n", msg)

				if err = c.WriteMessage(mt, []byte(c.Fiber.BaseURL()+":"+string(msg))); err != nil {
					log.Println("write:", err)
					break
				}
			}

		}))*/
}

func UseContext(context ContextInterface) {
	contextInstance = context
}

func WebSocket(url string, callback interface{}, params ...interface{}) {

	app.Get(url, func(c *fiber.Ctx) error {
		if websocket2.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}, websocket2.New(func(c *websocket2.Conn) {
		var ct = Context{
			fiber: c.Fiber,
		}
		contextInstance.New().WebSocket(&ct, callback, c)
		/*	var (
				mt  int
				msg []byte
				err error
			)
			for {
				if mt, msg, err = c.ReadMessage(); err != nil {
					log.Println("read:", err)
					break
				}
				if err = c.WriteMessage(mt, []byte( c.Fiber.BaseURL()+":"+string(msg) )); err != nil {
					log.Println("write:", err)
					break
				}
			}*/

	}))

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

func Run(ready func()) {
	ready()
	Events.Trigger("ready")
	fmt.Println("Starting server", Config.WebServer.Bind+":"+Config.WebServer.Port, "...")
	var err = app.Listen(Config.WebServer.Bind + ":" + Config.WebServer.Port)
	if err != nil {
		fmt.Println(err)
	}
}
