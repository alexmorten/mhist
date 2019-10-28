package mhist

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_timestampsFromFileName(t *testing.T) {
	t.Run("gets timestamps from filename", func(t *testing.T) {
		fileName := "1234-56789.csv"
		info, err := timestampsFromFileName(fileName)
		require.Nil(t, err)
		assert.EqualValues(t, 1234, info.oldestTs)
		assert.EqualValues(t, 56789, info.latestTs, 56789)
	})
}
