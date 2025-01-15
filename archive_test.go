package k6deps

import (
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
