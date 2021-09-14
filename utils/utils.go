package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func OpenFileInTextEditor(filename string) error {
	return runCommand("code", filename)
}

func OpenFileInSpreadsheetEditor(filename string) error {
	return runCommand("soffice.exe", "--calc", filename)
}

func ParseID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

func RemoveExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
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
