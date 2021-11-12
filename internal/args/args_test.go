package args

import (
	"fmt"
	"os"
	"testing"
)

func TestArgs(t *testing.T) {
	os.Args = []string{"-a","abcd","-b","cde", "-c"}
	fmt.Println("-a","=",Get("-a"))
	fmt.Println("-b","=",Get("-b"))
	fmt.Println("-c","=",Get("-c"))
}