package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/k6deps"
	"github.com/stretchr/testify/require"
)

//nolint:forbidigo,paralleltest
func Test_formatOptions(t *testing.T) {
	dir := t.TempDir()

	stderr := filepath.Clean(filepath.Join(dir, "output"))

	file, err := os.Create(stderr)
	require.NoError(t, err)

	saved := os.Stderr
	os.Stderr = file

	defer func() { os.Stderr = saved }()

	str := formatError(errors.ErrUnsupported)
	require.Equal(t, "✘ [ERROR] unsupported operation\n\n", str)

	str = formatError(k6deps.ErrConstraints)
	require.Equal(t, "✘ [ERROR] constraints error\n\n", str)

	require.NoError(t, file.Close())
}
