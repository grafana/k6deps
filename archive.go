package k6deps

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/k6pack"
)

//nolint:forbidigo
func loadMetadata(dir string, opts *Options) error {
	var meta archiveMetadata

	data, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "metadata.json"))
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &meta); err != nil {
		return err
	}

	opts.Manifest.Ignore = true // no manifest (yet) in archive

	opts.Script.Name = filepath.Join(
		dir,
		"file",
		filepath.FromSlash(strings.TrimPrefix(meta.Filename, "file:///")),
	)

	if value, found := meta.Env[EnvDependencies]; found {
		opts.Env.Name = EnvDependencies
		opts.Env.Contents = []byte(value)
	} else {
		opts.Env.Ignore = true
	}

	contents, err := os.ReadFile(filepath.Join(filepath.Clean(dir), "data"))
	if err != nil {
		return err
	}

	script, _, err := k6pack.Pack(string(contents), &k6pack.Options{Filename: opts.Script.Name})
	if err != nil {
		return err
	}

	opts.Script.Contents = script

	return nil
}

type archiveMetadata struct {
	Filename string            `json:"filename"`
	Env      map[string]string `json:"env"`
}

//nolint:forbidigo
func extractArchive(dir string, input io.Reader) error {
	reader := tar.NewReader(input)

	const maxFileSize = 1024 * 1024 * 10 // 10M

	for {
		header, err := reader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dir, filepath.Clean(filepath.FromSlash(header.Name)))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o750); err != nil {
				return err
			}

		case tar.TypeReg:
			if ext := filepath.Ext(target); ext == ".csv" || (ext == ".json" && filepath.Base(target) != "metadata.json") {
				continue
			}

			file, err := os.OpenFile(filepath.Clean(target), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode)) //nolint:gosec
			if err != nil {
				return err
			}

			if _, err := io.CopyN(file, reader, maxFileSize); err != nil && !errors.Is(err, io.EOF) {
				return err
			}

			if err = file.Close(); err != nil {
				return err
			}
		}
	}
}
