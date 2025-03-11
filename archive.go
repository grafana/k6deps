package k6deps

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"
	"slices"
)

type archiveMetadata struct {
	Filename string            `json:"filename"`
	Env      map[string]string `json:"env"`
}

const maxFileSize = 1024 * 1024 * 10 // 10M

func processArchive(input io.Reader) (analyzer, error) {
	reader := tar.NewReader(input)

	analyzers := make([]analyzer, 0)

	for {
		header, err := reader.Next()

		switch {
		case errors.Is(err, io.EOF):
			return mergeAnalyzers(analyzers...), nil
		case err != nil:
			return nil, err
		case header == nil:
			continue
		}

		if header.Typeflag != tar.TypeReg || !shouldProcess(header.Name) {
			continue
		}

		content := &bytes.Buffer{}
		if _, err := io.CopyN(content, reader, maxFileSize); err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		// if the file is metadata.json, we extract the dependencies from the env
		if header.Name == "metadata.json" {
			analyzer, err := analizeMetadata(content.Bytes())
			if err != nil {
				return nil, err
			}
			analyzers = append(analyzers, analyzer)
			continue
		}

		// analize the file content as an script
		target := filepath.Clean(filepath.FromSlash(header.Name))
		src := Source{
			Name:     target,
			Contents: content.Bytes(),
		}

		analyzers = append(analyzers, scriptAnalyzer(src))
	}
}

// indicates if the file should be processed during extraction
func shouldProcess(target string) bool {
	ext := filepath.Ext(target)
	return slices.Contains([]string{".js", ".ts"}, ext) || slices.Contains([]string{"metadata.json", "data"}, target)
}

// analizeMetadata extracts the dependencies from the metadata.json file
func analizeMetadata(content []byte) (analyzer, error) {
	metadata := archiveMetadata{}
	if err := json.Unmarshal(content, &metadata); err != nil {
		return nil, err
	}

	if value, found := metadata.Env[EnvDependencies]; found {
		src := Source{
			Name:     EnvDependencies,
			Contents: []byte(value),
		}
		return envAnalyzer(src), nil
	}

	return empty, nil
}
