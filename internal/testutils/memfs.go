package testutils

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

// Filemap defines a map of files, given their paths an content
type Filemap map[string][]byte

// MapFS return an afero.Fs from a Filemap. It creates the root directory and
// adds all files are under this directory. If the root directory is not absolute
// it is made absolute with respect of an OS specific root dir.
func MapFS(t *testing.T, root string, files Filemap) afero.Fs {
	t.Helper()

	if !filepath.IsAbs(root) {
		root = filepath.Join(OSRoot(), root)
	}

	memFS := afero.NewMemMapFs()
	for path, content := range files {
		file, err := memFS.Create(filepath.Join(root, path))
		if err != nil {
			t.Fatalf("creating file %v", err)
		}
		_, err = file.Write(content)
		if err != nil {
			t.Fatalf("writing content to %q %v", path, err)
		}
		err = file.Close()
		if err != nil {
			t.Fatalf("closing file %q %v", path, err)
		}
	}

	return memFS
}
