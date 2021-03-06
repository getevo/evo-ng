package watcher

import (
	"fmt"
	"os/exec"
	"path/filepath"
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

	return &builder{dir: dir, binary: bin, wd: wd, buildArgs: buildArgs}
}

// Binary returns its build binary's path
func (b *builder) Binary() string {
	return b.binary
}

// Build the Golang project set for this builder
func (b *builder) Build() error {
	fmt.Println("Building program")
	var command = exec.Command("go", "mod", "vendor")
	command.Dir = b.dir
	output, err := command.CombinedOutput()
	if !command.ProcessState.Success() {
		return fmt.Errorf("error building: %s", output)
	}

	args := append([]string{"go", "build", "-o", filepath.Join(b.wd, b.binary)}, b.buildArgs...)
	fmt.Println("Build command", args)

	command = exec.Command(args[0], args[1:]...) // nolint gas
	command.Dir = b.dir

	output, err = command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed with %v\n%s", err, output)
	}

	if !command.ProcessState.Success() {
		return fmt.Errorf("error building: %s", output)
	}

	return nil
}
