// Package exec runs external commands.
package editor

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		cmd  []string
	}{
		{
			name: "empty command",
			cmd:  nil,
		},
		{
			name: "name only",
			cmd:  []string{"foo"},
		},
		{
			name: "name and args",
			cmd:  []string{"foo", "--flag"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			editor := New(tt.cmd...)
			got := append([]string{editor.name}, editor.args...)
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%s", tt.cmd) {
				t.Errorf("got %v, want %v", got, tt.cmd)
			}
		})
	}
}

func Test_editor_Open(t *testing.T) {
	editor := New("echo", "hello world")

	var b bytes.Buffer
	editor.output = &b

	if err := editor.Open(""); err != nil {
		t.Fatal(err)
	}

	got := b.String()
	want := "hello world \n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
