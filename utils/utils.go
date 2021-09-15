package utils

import (
	"path/filepath"
	"strconv"
	"strings"
)

func ParseID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

func RemoveExt(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
