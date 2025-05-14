package k6deps_test

import (
	"testing"

	"github.com/grafana/k6deps"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeContents(t *testing.T) {
	t.Parallel()

	opts := &k6deps.Options{
		Script: k6deps.Source{
			Name:     "script.js",
			Contents: []byte(`"use k6 with k6/x/bar>=2.0";`),
		},
		Manifest: k6deps.Source{
			Name:     "package.json",
			Contents: []byte(`{"dependencies":{"k6/x/foo":">1.0"}}`),
		},
		Env: k6deps.Source{
			Name:     "DEPS",
			Contents: []byte(`k6/x/yaml>v0.1.0`),
		},
	}

	deps, err := k6deps.Analyze(opts)

	require.NoError(t, err)
	require.Len(t, deps, 3)
	require.Equal(t, deps["k6/x/bar"].Constraints.String(), ">=2.0")
	require.Equal(t, deps["k6/x/foo"].Constraints.String(), ">1.0")
	require.Equal(t, deps["k6/x/yaml"].Constraints.String(), ">v0.1.0")

	opts.Script.Contents = nil
	opts.Script.Name = "__NO__SUCH__SCRIPT__"

	_, err = k6deps.Analyze(opts)
	require.Error(t, err)

	opts.Script.Contents = []byte(`"use k6 with k6/x/faker>>1.0";`)
	_, err = k6deps.Analyze(opts)
	require.Error(t, err)
}
