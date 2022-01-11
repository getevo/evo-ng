package evo

import (
	"fmt"
	"github.com/fatih/color"
	"path/filepath"
	"runtime"
	"time"
)

func Register(registrable ...interface{}) {
	for _, fn := range registrable {
		var err = CallFn(fn)
		if err != nil {
			Panic(err)
		}
	}
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

func Trace(obj interface{}, i int) {
	fmt.Println(STrace(obj, i+1))
}

func STrace(obj interface{}, i int) string {
	_, file, line, _ := runtime.Caller(i)
	var path = filepath.Base(filepath.Dir(file)) + "/" + filepath.Base(file)
	return fmt.Sprint(time.Now().Format("2006-01-02 15:04:06 "), path+":"+fmt.Sprint(line)+":", fmt.Sprintf("%+v", obj))
}

func Panic(obj interface{}) {
	color.Red(STrace(obj, 2))
	panic(obj)
}
