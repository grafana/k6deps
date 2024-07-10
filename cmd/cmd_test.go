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
	out := filepath.Clean(filepath.Join(dir, "output"))

	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "-o", out, scriptfile})

	err := root.Execute()
	require.NoError(t, err)

	contents, err := os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `{"k6/x/faker":">v0.3.0","xk6-top":"*"}`+"\n", string(contents))

	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "--format", "json", "-o", out, scriptfile})

	err = root.Execute()
	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `{"k6/x/faker":">v0.3.0","xk6-top":"*"}`+"\n", string(contents))

	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "--format", "text", "-o", out, scriptfile})

	err = root.Execute()
	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `k6/x/faker>v0.3.0;xk6-top*`+"\n", string(contents))

	root.SetArgs([]string{"--ingnore-env", "--ignore-manifest", "--format", "js", "-o", out, scriptfile})

	err = root.Execute()
	require.NoError(t, err)

	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `"use k6 with k6/x/faker>v0.3.0";
"use k6 with xk6-top*";
`, string(contents))
}
