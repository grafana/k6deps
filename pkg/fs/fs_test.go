package fs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestRootFSOpen(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, "path", "to", "file")

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("test setup %v", err)
	}
	if _, err := os.Create(path); err != nil {
		t.Fatalf("test setup %v", err)
	}
	rootFs, err := NewFromDir(root)
	if err != nil {
		t.Fatalf("unexpected setting up test %v", err)
	}

	testCases := []struct {
		title  string
		path   string
		expect error
	}{
		{
			title:  "valid relative path",
			path:   filepath.Join("path", "to", "file"),
			expect: nil,
		},
		{
			title:  "valid navigation",
			path:   filepath.Join(".", "path", "to", "file"),
			expect: nil,
		},
		{
			title:  "absolute path",
			path:   filepath.Join(root, "path", "to", "file"),
			expect: nil,
		},
		{
			title:  "file does not exists",
			path:   filepath.Join("path", "to", "nonexiting"),
			expect: fs.ErrNotExist,
		},
		{
			title:  "invalid navigation",
			path:   filepath.Join("..", "scape", "to", "other", "file"),
			expect: fs.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			t.Parallel()
			_, err = rootFs.Open(tc.path)
			if !errors.Is(err, tc.expect) {
				t.Errorf("expected %v got %v", tc.expect, err)
			}
		})
	}
}
