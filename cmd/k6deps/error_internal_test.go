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
	require.Contains(t, str, "[ERROR] unsupported operation")

	str = formatError(k6deps.ErrConstraints)
	require.Contains(t, str, "[ERROR] constraints error")

	require.NoError(t, file.Close())
}
