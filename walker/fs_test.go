package walker_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/sf9133/fanal/analyzer"
	"github.com/sf9133/fanal/walker"
)

func TestWalkDir(t *testing.T) {
	// happy path
	err := walker.WalkDir("testdata/fs", func(filePath string, info os.FileInfo, opener analyzer.Opener) error {
		if info.IsDir() {
			return nil
		}
		if filePath == "testdata/fs/bar" {
			b, err := opener()
			require.NoError(t, err)
			assert.Equal(t, "bar", string(b))
		} else {
			assert.Fail(t, "invalid file", filePath)
		}

		return nil
	})
	require.NoError(t, err, "happy path")

	// sad path
	err = walker.WalkDir("testdata/fs", func(filePath string, info os.FileInfo, opener analyzer.Opener) error {
		return errors.New("error")
	})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "failed to analyze file: error", "sad path")
}
