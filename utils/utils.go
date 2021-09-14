package utils

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func OpenFileInTextEditor(filename string) error {
	cmd := exec.Command("code", filename)
	return cmd.Run()
}

func OpenFileInSpreadsheetEditor(filename string) error {
	cmd := exec.Command("soffice.exe", "--calc", filename)
	return cmd.Run()
}

func ParseID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

func RemoveExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
