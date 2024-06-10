package k6deps

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_loadSources(t *testing.T) {
	t.Setenv(EnvDependencies, "k6>0.49")

	opts := new(Options)

	dir := filepath.Join("testdata", "foo")
	adir, err := filepath.Abs(dir)

	require.NoError(t, err)

	opts.Script.Name = filepath.Join(dir, "foo.js")

	err = loadSources(opts)

	require.NoError(t, err)

	require.Equal(t, EnvDependencies, opts.Env.Name)
	require.Equal(t, "k6>0.49", string(opts.Env.Contents))

	require.Equal(t, filepath.Join(adir, "package.json"), opts.Manifest.Name)
	require.Contains(t, string(opts.Script.Contents), "var faker = require(\"k6/x/faker\");")

	opts = new(Options)

	opts.Script.Name = filepath.Join("testdata", "bad-package", "script.js")

	err = loadSources(opts)

	require.Error(t, err)

	opts = new(Options)

	opts.Script.Name = filepath.Join("testdata", "bad.js")

	err = loadSources(opts)

	require.Error(t, err)
}

func Test_loadManifest(t *testing.T) {
	t.Parallel()

	opts := new(Options)

	err := loadManifest(opts)

	require.NoError(t, err)

	require.Empty(t, opts.Manifest.Contents)
	require.Empty(t, opts.Manifest.Name)

	name := filepath.Join("testdata", "foo", "package.json")

	opts.Manifest.Ignore = true
	opts.Manifest.Name = name

	err = loadManifest(opts)

	require.NoError(t, err)
	require.Empty(t, opts.Manifest.Contents)
	require.Equal(t, name, opts.Manifest.Name)

	opts.Manifest.Ignore = false

	err = loadManifest(opts)

	require.NoError(t, err)
	require.Contains(t, string(opts.Manifest.Contents), "{\"dependencies\":{}}")

	opts.Manifest.Name = "no such file"

	err = loadManifest(opts)
	require.NoError(t, err)

	opts.Manifest.Contents = nil

	err = loadManifest(opts)

	require.Error(t, err)

	opts.Manifest.Name = ""
	opts.Script.Name = filepath.Join("testdata", "foo", "foo.js")

	err = loadManifest(opts)
	require.NoError(t, err)
	require.Contains(t, string(opts.Manifest.Contents), "{\"dependencies\":{}}")

	opts.Manifest.Name = ""
	opts.Manifest.Contents = nil
	opts.Script.Name = filepath.Join("testdata", "bad-package", "script.js")
	err = loadManifest(opts)

	require.Error(t, err)
}

func Test_loadScript(t *testing.T) {
	t.Parallel()

	opts := new(Options)

	err := loadScript(opts)

	require.NoError(t, err)
	require.Empty(t, opts.Script.Contents)
	require.Empty(t, opts.Script.Name)

	name := filepath.Join("testdata", "foo", "foo.js")
	aname, err := filepath.Abs(name)
	require.NoError(t, err)

	opts.Script.Name = name
	opts.Script.Ignore = true

	err = loadScript(opts)

	require.NoError(t, err)
	require.Empty(t, opts.Script.Contents)
	require.Equal(t, name, opts.Script.Name)

	opts.Script.Ignore = false

	err = loadScript(opts)

	require.NoError(t, err)
	require.Contains(t, string(opts.Script.Contents), "var faker = require(\"k6/x/faker\");")
	require.Equal(t, aname, opts.Script.Name)

	opts.Script.Name = filepath.Join("testdata", "bad.js")
	err = loadScript(opts)

	require.NoError(t, err)

	opts.Script.Contents = nil

	err = loadScript(opts)

	require.Error(t, err)

	opts.Script.Name = filepath.Join("testdata", "missing.js")

	err = loadScript(opts)

	require.Error(t, err)
}

func Test_loadEnv(t *testing.T) {
	opts := new(Options)

	t.Setenv(EnvDependencies, "")

	loadEnv(opts)
	require.Empty(t, opts.Env.Contents)
	require.Empty(t, opts.Env.Name)

	t.Setenv(EnvDependencies, "k6>0.49")

	opts.Env.Ignore = true

	loadEnv(opts)
	require.Empty(t, opts.Env.Contents)
	require.Empty(t, opts.Env.Name)

	opts.Env.Ignore = false

	loadEnv(opts)
	require.Equal(t, "k6>0.49", string(opts.Env.Contents))
	require.Equal(t, EnvDependencies, opts.Env.Name)
}

func Test_findManifest(t *testing.T) {
	t.Parallel()

	content, name, err := findManifest("testdata/foo/bar/bar.js")

	require.NoError(t, err)

	aname, err := filepath.Abs(filepath.Join("testdata", "foo", "package.json"))

	require.NoError(t, err)
	require.Equal(t, aname, name)
	require.Contains(t, string(content), "{\"dependencies\":{}}")

	content, name, err = findManifest("testdata/foo/foo.js")

	require.NoError(t, err)
	require.Equal(t, aname, name)
	require.Contains(t, string(content), "{\"dependencies\":{}}")

	content, name, err = findManifest(string(filepath.Separator))

	require.NoError(t, err)
	require.Nil(t, content)
	require.Empty(t, name)
}
