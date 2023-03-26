// Package osutil provides utilities that complement the os package in the
// standard library.
package osutil

import (
	"os"
	"path/filepath"
)

// MakeUserDir creates a directory along with any necessary parents with
// drwx------ permissions.
func MakeUserDir(elem ...string) (string, error) {
	dir := filepath.Join(elem...)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}
