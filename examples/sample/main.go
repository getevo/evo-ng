package main

//GENERATED BY EVO-NG
import (
	"github.com/getevo/evo-ng"
	"github.com/getevo/evo-ng/examples/sample/apps/myapp"
)

func main() {
	//Register EVO
	evo.Engine()

	//Register apps/myapp
	myapp.Register()

}
