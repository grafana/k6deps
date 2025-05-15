// Package rootfs implements extensions to go's fs.FS to work around its limitations
package rootfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var ErrInvalidPath = errors.New("invalid path") //nolint:revive

// FS defines an interface that extends go's fs.FS with a mechanism for working with absolute paths
type FS interface {
	fs.FS
	// makes path relative to the FS's root
	Rel(path string) (string, error)
	// returns FS's root dir
	Root() string
}

type rootFS struct {
	fs.FS
	root string
}

// NewFromDir create a FS from a root directory. The root must be an absolute path
func NewFromDir(root string) (FS, error) {
	if !filepath.IsAbs(root) {
		return nil, fmt.Errorf("%w: %q is not absolute", ErrInvalidPath, root)
	}

	return &rootFS{
		FS:   os.DirFS(root), //nolint:forbidigo
		root: root,
	}, nil
}

func (f *rootFS) Rel(path string) (string, error) {
	return filepath.Rel(f.Root(), path)
}

func (f *rootFS) Root() string {
	return f.root
}

// Open overrides os.DirFS.Open to ensure we have a path compatible with fs, not an operating system specific one,
// and it is relative to the FS's root.
//
// We need to make some validations to return ErrNotExist instead of ErrInvalid if the path is outside the fs.
// For example, if root is '/path/to/root' and path is '../wrong/path'
// filepath.join(root, path) will return '/path/to/wrong/path'
// filepath.Rel(root, '/path/to/wrong/path') will return '../wrong/path' which is invalid (it starts with ..)
func (f *rootFS) Open(path string) (fs.File, error) {
	var err error
	if !filepath.IsAbs(path) {
		path = filepath.Join(f.root, path)
	}
	// check if the path is outside the root
	if !strings.HasPrefix(path, f.root) {
		return nil, &fs.PathError{Path: path, Err: fs.ErrNotExist}
	}

	path, err = filepath.Rel(f.root, path)
	if err != nil {
		return nil, err
	}

	return f.FS.Open(filepath.ToSlash(filepath.Clean(path)))
}

// NewFromFS return a FS from a FS
func NewFromFS(fs fs.FS) FS {
	root := "/"
	if runtime.GOOS == "windows" {
		root = "M:\\"
	}

	return &rootFS{root: root, FS: fs}
}
