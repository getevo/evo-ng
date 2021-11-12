package myapp

import (
	"fmt"
	"os"
)

func Register()  {
	fmt.Println("hello!")
	fmt.Println(os.Args[1:])
}
