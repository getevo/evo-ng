package proc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var execFile = filepath.Base(os.Args[0])
var appName = strings.TrimSuffix(execFile, filepath.Ext(execFile))
var tempDir = os.TempDir() + "/" + appName
var hostName, _ = os.Hostname()

var App = struct {
	ExecFile string `json:"exec_file"`
	AppName  string `json:"app_name"`
	TempDir  string `json:"temp_dir"`
	HostName string `json:"host_name"`
}{
	ExecFile: execFile,
	AppName:  appName,
	TempDir:  tempDir,
	HostName: hostName,
}

func AppName() string {
	return appName
}

func Name() string {
	return filepath.Base(os.Args[0])
}

func AppDir() string {
	return filepath.Dir(os.Args[0])
}

func Args() []string {
	return os.Args[1:]
}

func Pid() int {
	return os.Getpid()
}

func Die(message ...interface{}) {
	fmt.Println(message...)
	os.Exit(0)
}

func TempDir() string {
	os.MkdirAll(tempDir, os.ModePerm)
	return tempDir
}

func AppID() string {
	return hostName
}

func WorkingDir() string {
	path, _ := os.Getwd()
	return path
}
