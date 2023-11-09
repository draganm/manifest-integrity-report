package consistency_test

import (
	"os"
	"testing"

	"github.com/draganm/manifest-integrity-report/consistency"
	"github.com/stretchr/testify/require"
)

func TestLoadResources(t *testing.T) {
	require := require.New(t)
	fs := os.DirFS("_fixtures")
	resources, err := consistency.LoadResources(fs)
	require.NoError(err)
	require.Len(resources, 4)
}
