// Package utils is a general utility package.
//
// TODO: Remove this package. Generic utility packages are bad practice because
// they have an undefined scope.
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
