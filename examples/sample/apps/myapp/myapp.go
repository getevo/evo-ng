package myapp

import (
	"fmt"
	"github.com/getevo/evo-ng/examples/sample/http"
	"os"
)

func Register() error {
	fmt.Println("hello!")
	fmt.Println(os.Args[1:])
	http.Get("/", func(context *http.Context) error {
		context.WriteResponse("Hey!")
		return nil
	})

	http.Use("/test", func(context *http.Context) error {

		return fmt.Errorf("you ask for /test/*")
	})
	http.Asset("/asset", "./assets")

	group := http.Group("/a")

	group.Get("/b", func(context *http.Context) error {
		context.WriteResponse("Hey")
		return nil
	})
	return nil
}
