package k6deps

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/grafana/k6pack"
)

// EnvDependencies holds the name of the environment variable thet describes additional dependencies.
const EnvDependencies = "K6_DEPENDENCIES"

// Source describes a generic dependency source.
// Such a source can be the k6 script, the manifest file, or an environment variable (e.g. K6_DEPENDENCIES).
type Source struct {
	// Name contains the name of the source (file, environment variable, etc.).
	Name string
	// Contents contains the content of the source (e.g. script)
	Contents []byte
	// Ignore disables automatic search and processing of that source.
	Ignore bool
}

// Options contains the parameters of the dependency analysis.
type Options struct {
	// Script contains the properties of the k6 test script to be analyzed.
	Script Source
	// Archive contains the properties of the k6 archive to be analyzed.
	// If archive is specified, the other three sources will not be taken into account,
	// since the archive may contain them.
	// An archive is a tar file, which can be created using the k6 archive command, for example.
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
	// If empty, os.LookupEnv will be used.
	LookupEnv func(key string) (value string, ok bool)
	// FindManifest function is used to find manifest file for the given scriptfile
	// if the Contents of Manifest option is empty.
	// If the scriptfile parameter is empty, FindManifest starts searching
	// for the manifest file from the current directory
	// If missing, the closest manifest file will be used.
	FindManifest func(scriptfile string) (contents []byte, filename string, ok bool, err error)
}

func (opts *Options) findManifest(filename string) ([]byte, string, bool, error) {
	if opts.FindManifest != nil {
		return opts.FindManifest(filename)
	}

	return findManifest(filename)
}

func (opts *Options) lookupEnv(key string) (string, bool) {
	if opts.LookupEnv != nil {
		return opts.LookupEnv(key)
	}

	return os.LookupEnv(key) //nolint:forbidigo
}

func loadSources(opts *Options) error {
	if !opts.Archive.Ignore && (len(opts.Archive.Contents) > 0 || len(opts.Archive.Name) > 0) {
		return loadArchive(opts)
	}

	if err := loadManifest(opts); err != nil {
		return err
	}

	if err := loadScript(opts); err != nil {
		return err
	}

	loadEnv(opts)

	return nil
}

func loadManifest(opts *Options) error {
	if len(opts.Manifest.Name) == 0 && !opts.Manifest.Ignore {
		pkg, pkgfile, found, err := opts.findManifest(opts.Script.Name)
		if err != nil {
			return err
		}

		if found {
			opts.Manifest.Name = pkgfile
			opts.Manifest.Contents = pkg
		}

		return nil
	}

	if len(opts.Manifest.Name) == 0 || len(opts.Manifest.Contents) > 0 || opts.Manifest.Ignore {
		return nil
	}

	pkg, err := os.ReadFile(opts.Manifest.Name) //nolint:forbidigo
	if err != nil {
		return err
	}

	opts.Manifest.Contents = pkg

	return nil
}

func loadScript(opts *Options) error {
	if len(opts.Script.Name) == 0 || len(opts.Script.Contents) > 0 || opts.Script.Ignore {
		return nil
	}

	scriptfile, err := filepath.Abs(opts.Script.Name)
	if err != nil {
		return err
	}

	contents, err := os.ReadFile(scriptfile) //nolint:forbidigo,gosec
	if err != nil {
		return err
	}

	script, _, err := k6pack.Pack(string(contents), &k6pack.Options{Filename: scriptfile})
	if err != nil {
		return err
	}

	opts.Script.Name = scriptfile
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

func findManifest(filename string) ([]byte, string, bool, error) {
	if len(filename) == 0 {
		filename = "any_file"
	}

	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, "", false, err
	}

	for dir := filepath.Dir(abs); ; dir = filepath.Dir(dir) {
		filename := filepath.Clean(filepath.Join(dir, "package.json"))
		if _, err := os.Stat(filename); !errors.Is(err, os.ErrNotExist) { //nolint:forbidigo
			contents, err := os.ReadFile(filename) //nolint:forbidigo
			return contents, filename, err == nil, err
		}

		if dir[len(dir)-1] == filepath.Separator {
			break
		}
	}

	return nil, "", false, nil
}

//nolint:forbidigo
func loadArchive(opts *Options) error {
	if opts.Archive.Ignore || (len(opts.Archive.Name) == 0 && len(opts.Archive.Contents) == 0) {
		return nil
	}

	var reader io.Reader

	if len(opts.Archive.Contents) == 0 {
		archivefile, err := filepath.Abs(opts.Archive.Name)
		if err != nil {
			return err
		}

		file, err := os.Open(filepath.Clean(archivefile))
		if err != nil {
			return err
		}

		defer file.Close() //nolint:errcheck

		reader = file
	} else {
		reader = bytes.NewReader(opts.Archive.Contents)
	}

	dir, err := os.MkdirTemp("", "k6deps-*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir) //nolint:errcheck

	err = extractArchive(dir, reader)
	if err != nil {
		return err
	}

	// archive should be self contained
	opts.Script.Ignore = true
	opts.Archive.Ignore = true
	opts.Env.Ignore = true
	opts.Manifest.Ignore = true

	return loadMetadata(dir, opts)
}
