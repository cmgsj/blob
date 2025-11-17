package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cmgsj/blob/pkg/util/files"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	var configFile string

	flag.StringVar(&configFile, "c", configFile, "config file")

	flag.Parse()

	var c config

	err := files.Decode(configFile, &c)
	if err != nil {
		return err
	}

	err = c.validate()
	if err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	if c.Version == "" || c.Version == "latest" {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/swagger-api/swagger-ui/releases/latest", nil)
		if err != nil {
			return fmt.Errorf("failed to create latest release request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send latest release request: %w", err)
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to get latest release: %s", resp.Status)
			}

			return fmt.Errorf("failed to get latest release: %s: %s", resp.Status, data)
		}

		var release struct {
			TagName string `json:"tag_name"`
		}

		err = json.NewDecoder(resp.Body).Decode(&release)
		if err != nil {
			return fmt.Errorf("failed to decode json latest release: %w", err)
		}

		c.Version = strings.TrimPrefix(release.TagName, "v")
	}

	cacheFile := filepath.Join(".cache", "swagger-ui", fmt.Sprintf("v%s.tar.gz", c.Version))

	cache, err := os.Open(cacheFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to open cache file: %w", err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://github.com/swagger-api/swagger-ui/archive/refs/tags/v%s.tar.gz", c.Version), nil)
		if err != nil {
			return fmt.Errorf("failed to create release asset request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send release asset request: %w", err)
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to get release asset: %s", resp.Status)
			}

			return fmt.Errorf("failed to get release asset: %s: %s", resp.Status, data)
		}

		err = os.MkdirAll(filepath.Dir(cacheFile), 0o700)
		if err != nil {
			return fmt.Errorf("failed to create cache dir: %w", err)
		}

		cache, err = os.Create(cacheFile)
		if err != nil {
			return fmt.Errorf("failed to create cache file: %w", err)
		}

		_, err = io.Copy(cache, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to copy to cache file: %w", err)
		}

		err = cache.Close()
		if err != nil {
			return fmt.Errorf("failed to close cache file: %w", err)
		}

		cache, err = os.Open(cacheFile)
		if err != nil {
			return fmt.Errorf("failed to open cache file: %w", err)
		}
	}

	defer func() { _ = cache.Close() }()

	gzip, err := gzip.NewReader(cache)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}

	defer func() { _ = gzip.Close() }()

	tar := tar.NewReader(gzip)

	srcDir := fmt.Sprintf("swagger-ui-%s/dist", c.Version)

	for {
		header, err := tar.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("failed to get tar file: %w", err)
		}

		if !strings.HasPrefix(header.Name, srcDir) {
			continue
		}

		dstPath := filepath.Join(c.Output, strings.TrimPrefix(header.Name, srcDir))

		if header.FileInfo().IsDir() {
			err = os.MkdirAll(dstPath, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create tar dir: %w", err)
			}

			continue
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("failed to create tar file: %w", err)
		}

		_, err = io.Copy(dstFile, tar)
		if err != nil {
			return fmt.Errorf("failed to copy tar file: %w", err)
		}

		err = dstFile.Close()
		if err != nil {
			return fmt.Errorf("failed to close tar file: %w", err)
		}
	}

	initializerFile := filepath.Join(c.Output, "swagger-initializer.js")

	err = os.WriteFile(initializerFile, []byte(c.Initializer), 0o600)
	if err != nil {
		return fmt.Errorf("failed to write initializer file: %w", err)
	}

	return nil
}

type config struct {
	Output      string `json:"output"            yaml:"output"`
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	Initializer string `json:"initializer"       yaml:"initializer"`
}

func (c config) validate() error {
	if c.Output == "" {
		return errors.New("output is required")
	}

	if c.Initializer == "" {
		return errors.New("initializer is required")
	}

	return nil
}
