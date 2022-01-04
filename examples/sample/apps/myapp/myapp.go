package myapp

import (
	"fmt"
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/examples/sample/http"
	"os"
)

func Register() error {
	fmt.Println("hello!")
	fmt.Println(os.Args[1:])
	evo.RegisterView("myapp","./apps/myapp/views")
	return nil
}

var group = http.Group("/a")

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
		context.N
		return context.View("myapp","test","name","John Doe",map[string]interface{}{
			"a":"A",
			"b":"B",
		})
	})
	return nil
}

func Ready() {

}
