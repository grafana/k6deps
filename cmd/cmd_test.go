package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/k6deps/cmd"
	"github.com/stretchr/testify/require"
)

//nolint:forbidigo
func Test_New(t *testing.T) {
	t.Parallel()

	root := cmd.New()
	require.Equal(t, "k6deps [flags] [script-file]", root.Use)

	dir := t.TempDir()

	scriptfile := filepath.Join("testdata", "script.js")
	archive := filepath.Join("testdata", "archive.tar")

	out := filepath.Clean(filepath.Join(dir, "output"))

	root = cmd.New()
	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "-o", out, scriptfile})
	err := root.Execute()
	require.NoError(t, err)

	contents, err := os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `{"k6/x/faker":">v0.3.0","xk6-top":"*"}`+"\n", string(contents))

	root = cmd.New()
	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "-o", out, archive})
	err = root.Execute()
	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `{"k6":">0.54","k6/x/faker":">0.4.0","k6/x/sql":">=1.0.1","k6/x/sql/driver/ramsql":"*"}`+"\n", string(contents))

	root = cmd.New()
	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "--format", "json", "-o", out, scriptfile})
	err = root.Execute()
	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `{"k6/x/faker":">v0.3.0","xk6-top":"*"}`+"\n", string(contents))

	root = cmd.New()
	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "--format", "text", "-o", out, scriptfile})
	err = root.Execute()
	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `k6/x/faker>v0.3.0;xk6-top*`+"\n", string(contents))

	root = cmd.New()
	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "--format", "js", "-o", out, scriptfile})
	err = root.Execute()
	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `"use k6 with k6/x/faker>v0.3.0";
"use k6 with xk6-top*";
`, string(contents))

	root = cmd.New()
	stdin := os.Stdin
	os.Stdin, err = os.Open(scriptfile) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "--input", "js", "--format", "text", "-o", out})
	err = root.Execute()
	os.Stdin = stdin

	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `k6/x/faker>v0.3.0;xk6-top*`+"\n", string(contents))

	root = cmd.New()
	stdin = os.Stdin
	os.Stdin, err = os.Open(archive) //nolint:gosec
	if err != nil {
		t.Fatal(err)
	}
	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "--input", "tar", "--format", "json", "-o", out})
	err = root.Execute()
	os.Stdin = stdin

	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `{"k6":">0.54","k6/x/faker":">0.4.0","k6/x/sql":">=1.0.1","k6/x/sql/driver/ramsql":"*"}`+"\n", string(contents))
}
