package myapp

import (
	"fmt"
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/examples/sample/http"
	"github.com/getevo/evo-ng/websocket"
	"log"
	"os"
	"time"
)

func Register() error {
	fmt.Println("hello!")
	fmt.Println(os.Args[1:])
	evo.RegisterView("myapp", "./apps/myapp/views")
	go func() {
		for {
			time.Sleep(1 * time.Second)
			message <- fmt.Sprint(time.Now().Unix())
		}
	}()
	return nil
}

var group = http.Group("/a")
var message = make(chan string)

func Router() error {
	http.Get("/", func(context *http.Context) error {
		context.WriteResponse("Hey!")
		return nil
	})

	http.Use("/test", func(context *http.Context) error {

		return fmt.Errorf("you ask for /test/*")
	})
	http.Asset("/asset", "./assets")

	group.Get("/b", func(context *http.Context) error {
		fmt.Println("something reaches here")
		context.WriteResponse("Hey")
		return nil
	})

	http.Get("/view", func(context *http.Context) error {
		context.Message.Error("Test error")
		return context.View("myapp", "test", "name", "John Doe", map[string]interface{}{
			"a": "A",
			"b": "B",
		})
	})

	http.Get("/panic", func(context *http.Context) error {
		var m map[string]interface{}
		return m["1"].(error)
	})

	http.WebSocket("/ws", func(context *http.Context, c *websocket.Conn) error {

		for {
			msg := <-message

			if err := c.WriteMessage(1, []byte(msg)); err != nil {
				log.Println("write:", err)
				break
			}
		}
		return nil
	})
	return nil
}

func Ready() {

}
