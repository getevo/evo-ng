package myapp

import (
	"fmt"
	"github.com/getevo/evo-ng/examples/sample/request"
	"os"
)

func Register() error {
	fmt.Println("hello!")
	fmt.Println(os.Args[1:])
	request.Get("/", func(context *request.Context) error {
		context.Base.Write("Hey!")
		return nil
	})
	return nil
}
