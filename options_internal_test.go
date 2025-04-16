package k6deps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_loadScript(t *testing.T) {
	t.Parallel()

	opts := new(Options)

	err := opts.loadScript()

	require.NoError(t, err)
	require.Empty(t, opts.Script.Contents)
	require.Empty(t, opts.Script.Name)

	name := filepath.Join("testdata", "foo", "foo.js")

	opts.Script.Name = name
	opts.Script.Ignore = true

	err = opts.loadScript()

	require.NoError(t, err)
	require.Empty(t, opts.Script.Contents)
	require.Equal(t, name, opts.Script.Name)

	opts.Script.Ignore = false

	err = opts.loadScript()

	require.NoError(t, err)
	require.Contains(t, string(opts.Script.Contents), "var faker = __require(\"k6/x/faker\");")

	opts.Script.Name = filepath.Join("testdata", "bad.js")
	err = opts.loadScript()

	require.NoError(t, err)

	opts.Script.Contents = nil

	err = opts.loadScript()

	require.Error(t, err)

	opts.Script.Name = filepath.Join("testdata", "missing.js")

	err = opts.loadScript()

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

	name, found, err := findManifest("testdata/foo/bar/bar.js")
	require.NoError(t, err)
	require.True(t, found)

	aname, err := filepath.Abs(filepath.Join("testdata", "foo", "package.json"))

	require.NoError(t, err)
	require.Equal(t, aname, name)

	name, found, err = findManifest("testdata/foo/foo.js")

	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, aname, name)

	name, found, err = findManifest(string(filepath.Separator))
	require.NoError(t, err)
	require.False(t, found)
	require.Empty(t, name)
}

//nolint:forbidigo,paralleltest
func Test_findManifest_empty_arg(t *testing.T) {
	pwd, err := os.Getwd()
	defer func() {
		require.NoError(t, os.Chdir(pwd))
	}()

	require.NoError(t, err)

	require.NoError(t, os.Chdir(filepath.Join("testdata", "foo", "bar")))
	name, found, err := findManifest("")

	require.NoError(t, err)
	require.True(t, found)

	aname := filepath.Join(pwd, "testdata", "foo", "package.json")

	require.NoError(t, err)
	require.Equal(t, aname, name)
}
