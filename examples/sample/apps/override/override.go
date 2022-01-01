package override

import "github.com/getevo/evo-ng/examples/sample/http"

func OverRide() error {
	http.Get("/a/b", func(context *http.Context) error {
		context.WriteResponse("over ride!")
		return nil
	})
	return nil
}
