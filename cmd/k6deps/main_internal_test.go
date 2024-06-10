package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:forbidigo
func Test_main(t *testing.T) {
	dir := t.TempDir()

	t.Setenv("K6_DEPENDENCIES", "xk6-bar>1.0")

	scriptfile := filepath.Join("testdata", "script.js")
	out := filepath.Clean(filepath.Join(dir, "output"))

	os.Args = []string{appname, "-o", out, scriptfile}
	main()
	contents, err := os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `{"k6/x/faker":">v0.3.0","xk6-bar":">1.0","xk6-toml":">2.0","xk6-top":"*"}`+"\n", string(contents))

	os.Args = []string{appname, "--format", "json", "-o", out, scriptfile}
	main()
	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `{"k6/x/faker":">v0.3.0","xk6-bar":">1.0","xk6-toml":">2.0","xk6-top":"*"}`+"\n", string(contents))

	os.Args = []string{appname, "--format", "text", "-o", out, scriptfile}
	main()
	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `k6/x/faker>v0.3.0;xk6-bar>1.0;xk6-toml>2.0;xk6-top*`+"\n", string(contents))

	os.Args = []string{appname, "--format", "js", "-o", out, scriptfile}
	main()
	contents, err = os.ReadFile(out)
	require.NoError(t, err)
	require.Equal(t, `"use k6 with k6/x/faker>v0.3.0";
"use k6 with xk6-bar>1.0";
"use k6 with xk6-toml>2.0";
"use k6 with xk6-top*";
`, string(contents))
}
