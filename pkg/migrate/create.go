package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CreateMigration creates a new migration file in the specified directory
func CreateMigration(path, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("migration name cannot be empty")
	}

	// Create directory if it does not exist
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("error to create directory %s: %w", path, err)
	}

	// Timestamp actual YYYYMMDDHHMMSS
	timestamp := time.Now().Format("20060102150405")

	fileName := fmt.Sprintf("%s_%s.sql", timestamp, name)
	filePath := filepath.Join(path, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error to create file %s: %w", filePath, err)
	}
	defer file.Close()

	header := fmt.Sprintf("-- Migration: %s\n-- Created at: %s\n\n", fileName, time.Now().Format(time.RFC3339))
	if _, err := file.WriteString(header); err != nil {
		return "", fmt.Errorf("error to write to file %s: %w", filePath, err)
	}

	return filePath, nil
}
