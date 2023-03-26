// Package config manages the application configurations.
package config

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestConfig_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		json    string
		wantErr bool
	}{
		{
			name:    "should fail when JSON is invalid",
			json:    ``,
			wantErr: true,
		},
		{
			name: "should unmarshal empty object",
			json: `{}`,
			config: &Config{
				TextEditor:        DefaultTextEditor,
				SpreadsheetEditor: DefaultSpreadsheetEditor,
			},
			wantErr: false,
		},
		{
			name: "should set custom editors",
			json: `{
				"textEditor": ["foo"],
				"spreadsheetEditor": ["bar"]
			}`,
			config: &Config{
				TextEditor:        []string{"foo"},
				SpreadsheetEditor: []string{"bar"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			err := c.UnmarshalJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Fatalf("got error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(c, tt.config) {
				t.Fatalf("got %+v, want %+v", c, tt.config)
			}
		})
	}
}

func Test_load(t *testing.T) {
	t.Run("should return error if file does not exist", func(t *testing.T) {
		_, err := load(t.TempDir())
		if err == nil {
			t.Fatal(err)
		}
	})
	t.Run("should not return error if file does exist", func(t *testing.T) {
		dir := t.TempDir()
		if _, err := initFile(dir); err != nil {
			t.Fatal(err)
		}
		_, err := load(dir)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("should return error if file does not contain JSON", func(t *testing.T) {
		dir := t.TempDir()
		file, err := os.Create(joinFilename(dir))
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		_, err = load(dir)
		if err == nil {
			t.Fatal(err)
		}
	})
}

func Test_initFile(t *testing.T) {
	t.Run("should create file if it doesn't exist", func(t *testing.T) {
		dir := t.TempDir()
		if _, err := initFile(dir); err != nil {
			t.Fatal(err)
		}
		_, err := os.Stat(dir)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("should not create file if it already exists", func(t *testing.T) {
		dir := t.TempDir()
		// Create the file.
		if _, err := initFile(dir); err != nil {
			t.Fatal(err)
		}
		// Try to create it again.
		if _, err := initFile(dir); !os.IsExist(err) {
			t.Fatal("expected os.ErrExist")
		}
	})

	t.Run("should write JSON", func(t *testing.T) {
		dir := t.TempDir()
		// Create the file.
		if _, err := initFile(dir); err != nil {
			t.Fatal(err)
		}
		b, err := os.ReadFile(joinFilename(dir))
		if err != nil {
			t.Fatal(err)
		}
		raw := &json.RawMessage{}
		if err := json.Unmarshal(b, raw); err != nil {
			t.Fatal(err)
		}
	})
}
