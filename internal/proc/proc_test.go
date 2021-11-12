package proc

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(Name())
}

func TestPid(t *testing.T) {
	fmt.Println(Pid())
}


func TestArgs(t *testing.T) {
	fmt.Println(Args())
}


func TestDie(t *testing.T) {
	Die(1,2,3)
}

func TestTempDir(t *testing.T) {
	fmt.Println( TempDir() )
}
