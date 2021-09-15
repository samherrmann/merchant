package exec

import (
	"fmt"
	"os"
	"os/exec"
)

// RunTextEditor opens filename in the text editor.
func RunTextEditor(filename string) error {
	name := TextEditorCmd[0]
	args := append(TextEditorCmd[1:], filename)
	return runCommand(name, args...)
}

// RunSpreadsheetEditor opens filename in the spreadsheet editor.
func RunSpreadsheetEditor(filename string) error {
	name := SpreadsheetEditorCmd[0]
	args := append(SpreadsheetEditorCmd[1:], filename)
	return runCommand(name, args...)
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cannot open %v: %w", name, err)
	}
	return nil
}
