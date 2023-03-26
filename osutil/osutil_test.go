package osutil

import (
	"os"
	"testing"
)

func TestMakeUserDir(t *testing.T) {
	t.Run("should create directory", func(t *testing.T) {
		tempDir := t.TempDir()
		dir, err := MakeUserDir(tempDir, "foo", "bar")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(dir); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("should not return error if dir already exists", func(t *testing.T) {
		tempDir := t.TempDir()
		dir, err := MakeUserDir(tempDir, "foo", "bar")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(dir); err != nil {
			t.Fatal(err)
		}
		_, err = MakeUserDir(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("should only give owner permission", func(t *testing.T) {
		tempDir := t.TempDir()
		dir, err := MakeUserDir(tempDir, "foo", "bar")
		if err != nil {
			t.Fatal(err)
		}
		stat, err := os.Stat(dir)
		if err != nil {
			t.Fatal(err)
		}
		got := stat.Mode().String()
		want := "drwx------"
		if got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	})
}
