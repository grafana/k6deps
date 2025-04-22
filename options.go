package k6deps

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/grafana/k6deps/internal/pack"
)

const (
	// EnvDependencies holds the name of the environment variable that describes additional dependencies.
	EnvDependencies = "K6_DEPENDENCIES"

	// ManifestFileName is the name of the manifest file
	ManifestFileName = "package.json"
)

// Source describes a generic dependency source.
// Such a source can be the k6 script, the manifest file, or an environment variable (e.g. K6_DEPENDENCIES).
type Source struct {
	// Name contains the name of the source (file, environment variable, etc.).
	Name string
	// Reader provides streaming access to the source content as an alternative to Contents.
	Reader io.Reader
	// Contents contains the content of the source (e.g. script)
	Contents []byte
	// Ignore disables automatic search and processing of that source.
	Ignore bool
}

// IsEmpty returns true if the source is empty.
func (s *Source) IsEmpty() bool {
	return len(s.Contents) == 0 && s.Reader == nil && len(s.Name) == 0
}

func nopCloser() error {
	return nil
}

// contentReader returns a reader for the source content.
func (s *Source) contentReader() (io.Reader, func() error, error) {
	if s.Reader != nil {
		return s.Reader, nopCloser, nil
	}

	if len(s.Contents) > 0 {
		return bytes.NewReader(s.Contents), nopCloser, nil
	}

	fileName, err := filepath.Abs(s.Name)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Open(filepath.Clean(fileName)) //nolint:forbidigo
	if err != nil {
		return nil, nil, err
	}

	return file, file.Close, nil
}

// Options contains the parameters of the dependency analysis.
type Options struct {
	// Script contains the properties of the k6 test script to be analyzed.
	// If the name is specified, but no content is provided, the script is read from the file.
	// Any script file referenced will be recursively loaded and the dependencies merged.
	// If the Ignore property is set, the script will not be analyzed.
	Script Source
	// Archive contains the properties of the k6 archive to be analyzed.
	// If archive is specified, the other three sources will not be taken into account,
	// since the archive may contain them.
	// It is assumed that the script and all dependencies are in the archive. No external dependencies are analyzed.
	// An archive is a tar file, which can be created using the k6 archive command.
	Archive Source
	// Manifest contains the properties of the manifest file to be analyzed.
	// If the Ignore property is not set and no manifest file is specified,
	// the package.json file closest to the script is searched for.
	Manifest Source
	// Env contains the properties of the environment variable to be analyzed.
	// If the Ignore property is not set and no variable is specified,
	// the value of the variable named K6_DEPENDENCIES is read.
	Env Source
	// LookupEnv function is used to query the value of the environment variable
	// specified in the Env option Name if the Contents of the Env option is empty.
	// If not provided, os.LookupEnv will be used.
	LookupEnv func(key string) (value string, ok bool)
	// FindManifest function is used to find manifest if the path is not set explicitly in the options.
	// If not provided a default function is used. This function starts at the path to the script and traverses it
	// upwards until a manifest is found.
	// If the scriptPath parameter is empty, it starts searching from the manifest file from the current directory
	FindManifest func(scriptPath string) (filename string, ok bool, err error)
	// Fs is the file system to use for accessing files. If not provided, os file system is used
	Fs fs.FS
}

func (opts *Options) lookupEnv(key string) (string, bool) {
	if opts.LookupEnv != nil {
		return opts.LookupEnv(key)
	}

	return os.LookupEnv(key) //nolint:forbidigo
}

// returns the FS to use with this options
func (opts *Options) fs() fs.FS {
	if opts.Fs != nil {
		return opts.Fs
	}

	return os.DirFS(".") //nolint:forbidigo
}

func loadScript(opts *Options) error {
	if len(opts.Script.Name) == 0 || len(opts.Script.Contents) > 0 || opts.Script.Ignore {
		return nil
	}

	contents, err := fs.ReadFile(opts.fs(), opts.Script.Name)
	if err != nil {
		return err
	}

	script, _, err := pack.Pack(string(contents), &pack.Options{Filename: opts.Script.Name})
	if err != nil {
		return err
	}

	opts.Script.Contents = script

	return nil
}

func loadEnv(opts *Options) {
	if len(opts.Env.Contents) > 0 || opts.Env.Ignore {
		return
	}

	key := opts.Env.Name
	if len(key) == 0 {
		key = EnvDependencies
	}

	value, found := opts.lookupEnv(key)
	if !found || len(value) == 0 {
		return
	}

	opts.Env.Name = key
	opts.Env.Contents = []byte(value)
}

func (opts *Options) setManifest() error {
	// if the manifest is not provided, we try to find it
	// starting from the location of the script
	findManifest := opts.FindManifest
	if findManifest == nil {
		findManifest = opts.findManifest
	}

	path, found, err := findManifest(filepath.Dir(opts.Script.Name))
	if err != nil {
		return err
	}
	if found {
		opts.Manifest.Name = path
	}

	return nil
}

// looks for a package.json file
func (opts *Options) findManifest(basePath string) (string, bool, error) {
	manifestPath := ManifestFileName

	for {
		searchPath := filepath.Join(basePath, manifestPath)

		_, err := fs.Stat(opts.fs(), searchPath)
		if err == nil {
			return searchPath, true, nil
		}
		if errors.Is(err, fs.ErrNotExist) {
			manifestPath = filepath.Join("../", manifestPath)
			continue
		}
		if errors.Is(err, fs.ErrInvalid) {
			return "", false, nil
		}
		return "", false, err
	}
}
