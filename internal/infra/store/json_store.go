package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// JSONStore provides atomic write operations for JSON files.
type JSONStore struct {
	baseDir string
}

// NewJSONStore creates a new JSONStore with the given base directory.
func NewJSONStore(baseDir string) *JSONStore {
	return &JSONStore{baseDir: baseDir}
}

// Read reads a JSON file and unmarshals it into the target.
func (js *JSONStore) Read(filename string, target interface{}) error {
	path := filepath.Join(js.baseDir, filename)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %w", err)
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("file is empty")
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// Write atomically writes a JSON file by writing to a temp file first, then renaming.
func (js *JSONStore) Write(filename string, data interface{}) error {
	path := filepath.Join(js.baseDir, filename)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to temp file
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Sync to ensure data is written
	file, err := os.OpenFile(tmpPath, os.O_RDWR, 0644)
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to open temp file for sync: %w", err)
	}
	if err := file.Sync(); err != nil {
		file.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	file.Close()

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// Exists checks if a file exists.
func (js *JSONStore) Exists(filename string) bool {
	path := filepath.Join(js.baseDir, filename)
	_, err := os.Stat(path)
	return err == nil
}
