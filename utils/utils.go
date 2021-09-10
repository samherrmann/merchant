package utils

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func RunVSCode(filename string) error {
	cmd := exec.Command("code", filename)
	return cmd.Run()
}

func ParseID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

func RemoveExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
