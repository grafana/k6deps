package k6deps

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_loadScript(t *testing.T) {
	t.Parallel()

	opts := new(Options)

	err := loadScript(opts)

	require.NoError(t, err)
	require.Empty(t, opts.Script.Contents)
	require.Empty(t, opts.Script.Name)

	name := filepath.Join("testdata", "foo", "foo.js")
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
	require.Contains(t, string(opts.Script.Contents), "var faker = __require(\"k6/x/faker\");")

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

	opts := new(Options)

	expected := filepath.Join("testdata", "foo", "package.json")

	path, found, err := opts.findManifest(filepath.Join("testdata", "foo", "bar"))
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, expected, path)

	path, found, err = opts.findManifest(filepath.Join("testdata", "foo"))
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, expected, path)

	path, found, err = opts.findManifest("")
	require.NoError(t, err)
	require.False(t, found)
	require.Empty(t, path)
}
