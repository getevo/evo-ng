package myapp

import (
	"fmt"
	"github.com/getevo/evo-ng/examples/sample/http"
	"os"
)

func Register() error {
	fmt.Println("hello!")
	fmt.Println(os.Args[1:])
	return nil
}

var group = http.Group("/a")

func Routers() error {
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

	return nil
}

func Ready() {

}
