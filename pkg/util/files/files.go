package files

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Encode(path string, v any) error {
	if path == "" {
		return errors.New("path is required")
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() { _ = file.Close() }()

	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		err = encodeJSON(file, v)
		if err != nil {
			return fmt.Errorf("failed to encode json file: %w", err)
		}

	case ".yml", ".yaml":
		err = encodeYAML(file, v)
		if err != nil {
			return fmt.Errorf("failed to encode yaml file: %w", err)
		}

	default:
		return fmt.Errorf("unsupported file format %q", ext)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	return nil
}

func Decode(path string, v any) error {
	if path == "" {
		return errors.New("path is required")
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer func() { _ = file.Close() }()

	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		err = decodeJSON(file, v)
		if err != nil {
			return fmt.Errorf("failed to decode json file: %w", err)
		}

	case ".yml", ".yaml":
		err = decodeYAML(file, v)
		if err != nil {
			return fmt.Errorf("failed to decode yaml file: %w", err)
		}

	default:
		return fmt.Errorf("unsupported file format %q", ext)
	}

	return nil
}

func encodeJSON(writer io.Writer, v any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")

	return encoder.Encode(v)
}

func decodeJSON(reader io.Reader, v any) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	return decoder.Decode(v)
}

func encodeYAML(writer io.Writer, v any) error {
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)

	return encoder.Encode(v)
}

func decodeYAML(reader io.Reader, v any) error {
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	return decoder.Decode(v)
}
