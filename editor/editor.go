// Package editor provides the ability to open files in an external editor.
package editor

import (
	"io"
	"os"
	"os/exec"
)

type Editor interface {
	Open(filename string) error
}

// New creates a new Editor from the given command.
func New(cmd ...string) *editor {
	name := ""
	if len(cmd) > 0 {
		name = cmd[0]
	}
	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	return &editor{
		name:   name,
		args:   args,
		input:  os.Stdin,
		output: os.Stdout,
		errs:   os.Stderr,
	}
}

type editor struct {
	name   string
	args   []string
	input  io.Reader
	output io.Writer
	errs   io.Writer
}

func (e *editor) Open(filename string) error {
	args := e.args
	args = append(args, filename)
	cmd := exec.Command(e.name, args...)
	cmd.Stdin = e.input
	cmd.Stdout = e.output
	cmd.Stderr = e.errs
	return cmd.Run()
}
