package watcher

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Builder is a interface for the build process
type Builder interface {
	Build() error
	Binary() string
}

type builder struct {
	dir       string
	binary    string
	errors    string
	wd        string
	buildArgs []string
	run       bool
}

// NewBuilder creates a new builder
func NewBuilder(dir string, bin string, wd string, buildArgs []string) Builder {
	// resolve bin name by current folder name
	if bin == "" {
		bin = filepath.Base(wd)
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == OSWindows {
		// check if it already has the .exe extension
		if !strings.HasSuffix(bin, ".exe") {
			bin += ".exe"
		}
	}

	return &builder{dir: dir, binary: bin, wd: wd, buildArgs: buildArgs}
}

// Binary returns its build binary's path
func (b *builder) Binary() string {
	return b.binary
}

// Build the Golang project set for this builder
func (b *builder) Build() error {
	fmt.Println("Building program")
	args := append([]string{"go", "build", "-o", filepath.Join(b.wd, b.binary)}, b.buildArgs...)
	fmt.Println("Build command", args)

	command := exec.Command(args[0], args[1:]...) // nolint gas
	command.Dir = b.dir

	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed with %v\n%s", err, output)
	}

	if !command.ProcessState.Success() {
		return fmt.Errorf("error building: %s", output)
	}

	return nil
}