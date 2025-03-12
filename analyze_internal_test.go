package k6deps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty(t *testing.T) {
	t.Parallel()

	inst, err := empty()
	require.NoError(t, err)
	require.Empty(t, inst)
}

func Test_filterInvalid(t *testing.T) {
	t.Parallel()

	deps := Dependencies{
		"k6":               &Dependency{Name: "k6"},
		"foo":              &Dependency{Name: "foo"},
		"bar":              &Dependency{Name: "bar"},
		"k6/x/faker":       &Dependency{Name: "k6/x/faker"},
		"xk6-foo":          &Dependency{Name: "xk6-foo"},
		"@grafana/xk6-bar": &Dependency{Name: "@grafana/xk6-bar"},
	}

	valid := filterInvalid(deps)

	require.Len(t, valid, 4)
	require.Contains(t, valid, "k6")
	require.Contains(t, valid, "k6/x/faker")
	require.Contains(t, valid, "xk6-foo")
	require.Contains(t, valid, "@grafana/xk6-bar")
}

func Test_manifestAnalyzer(t *testing.T) {
	t.Parallel()

	src := Source{Name: "package.json"}
	fn := manifestAnalyzer(src)
	deps, err := fn()

	require.NoError(t, err)
	require.Empty(t, deps)

	src.Contents = []byte(`{"dependencies":{"@grafana/xk6-faker":"*"}}`)
	fn = manifestAnalyzer(src)
	deps, err = fn()

	require.NoError(t, err)
	require.NotEmpty(t, deps)
	require.Len(t, deps, 1)
	require.Equal(t, deps["@grafana/xk6-faker"].Constraints.String(), "*")

	src.Contents = []byte(`{`)
	fn = manifestAnalyzer(src)
	_, err = fn()
	require.Error(t, err)
}

func Test_scriptAnalyzer(t *testing.T) {
	t.Parallel()

	src := Source{Name: "script.js"}
	fn := scriptAnalyzer(src)
	deps, err := fn()

	require.Error(t, err)
	require.Empty(t, deps)

	src.Contents = []byte(`"use k6 with @grafana/xk6-faker>v0.3.0";`)
	fn = scriptAnalyzer(src)
	deps, err = fn()

	require.NoError(t, err)
	require.NotEmpty(t, deps)
	require.Len(t, deps, 1)
	require.Equal(t, deps["@grafana/xk6-faker"].Constraints.String(), ">v0.3.0")

	src.Contents = []byte(`"use k6 with k6/x/faker>>1.0";`)
	fn = scriptAnalyzer(src)
	_, err = fn()
	require.Error(t, err)
}

func Test_envAnalyzer(t *testing.T) {
	t.Parallel()

	src := Source{Name: "DEPS"}
	fn := envAnalyzer(src)
	deps, err := fn()

	require.NoError(t, err)
	require.Empty(t, deps)

	src.Contents = []byte(`@grafana/xk6-faker>v0.3.0`)
	fn = envAnalyzer(src)
	deps, err = fn()

	require.NoError(t, err)
	require.NotEmpty(t, deps)
	require.Len(t, deps, 1)
	require.Equal(t, deps["@grafana/xk6-faker"].Constraints.String(), ">v0.3.0")

	src.Contents = []byte(`k6/x/faker>>1.0`)
	fn = envAnalyzer(src)
	_, err = fn()
	require.Error(t, err)
}

func Test_mergeAnalyzers_error(t *testing.T) {
	t.Parallel()

	src := Source{Name: "DEPS"}
	src.Contents = []byte(`@grafana/xk6-faker>v0.3.0`)
	fn := envAnalyzer(src)

	src.Contents = []byte(`@grafana/xk6-foo>>v0.4.0`)
	fn2 := envAnalyzer(src)

	_, err := mergeAnalyzers(fn, fn2)()
	require.Error(t, err)

	src.Contents = []byte(`@grafana/xk6-faker>v0.4.0`)
	fn2 = envAnalyzer(src)

	_, err = mergeAnalyzers(fn, fn2)()
	require.Error(t, err)
}
