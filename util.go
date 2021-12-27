package evo

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func Register(registrable ...interface{}) {
	for _, fn := range registrable {
		var err = CallFn(fn)
		if err != nil {
			panic(STrace(err, 2))
		}
	}
}

func Trace(obj interface{}, i int) {
	fmt.Println(STrace(obj, i+1))
}

func STrace(obj interface{}, i int) string {
	_, file, line, _ := runtime.Caller(i)
	var path = filepath.Base(filepath.Dir(file)) + "/" + filepath.Base(file)
	return fmt.Sprint(time.Now().Format("2006-01-02 15:04:06 "), path+":"+fmt.Sprint(line)+":", fmt.Sprintf("%+v", obj))
}

func CallFn(fn interface{}) error {
	if v, ok := fn.(func()); ok {
		v()
		return nil
	} else if v, ok := fn.(func() error); ok {
		return v()
	}
	return fmt.Errorf("invalid function %+v", fn)
}

func Panic(obj interface{}) {
	fmt.Println("PANIC", STrace(obj, 2))
	os.Exit(1)
}
