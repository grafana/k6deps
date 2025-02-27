package k6deps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_analyzeArchive(t *testing.T) {
	t.Parallel()

	opts := &Options{
		Archive: Source{Name: filepath.Join("testdata", "archive.tar")},
	}

	actual, err := Analyze(opts)

	require.NoError(t, err)

	opts = &Options{
		Script:   Source{Name: filepath.Join("testdata", "combined.js")},
		Manifest: Source{Ignore: true},
		Env:      Source{Ignore: true},
	}

	expected, err := Analyze(opts)

	require.NoError(t, err)

	require.Equal(t, expected.String(), actual.String())
}

func Test_analyzeArchive_Reader(t *testing.T) {
	t.Parallel()

	file, err := os.Open(filepath.Join("testdata", "archive.tar")) //nolint:forbidigo
	require.NoError(t, err)
	defer file.Close() //nolint:errcheck

	opts := &Options{
		Archive: Source{Reader: file},
	}

	actual, err := Analyze(opts)
	require.NoError(t, err)

	expected, err := Analyze(opts)
	require.NoError(t, err)

	require.Equal(t, expected.String(), actual.String())
}
