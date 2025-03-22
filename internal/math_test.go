package internal_test

import (
	"dfxluna/go-summarize/internal"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDotProduct(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	v1 := []int{1, 2, 3}
	v2 := []int{0, 0, 0}

	dp, err := internal.DotProduct(v1, v2)
	require.NoError(err)
	assert.Equal(0, dp)

	v1 = []int{1, 2, 3}
	v2 = []int{1, 1, 1}

	dp, err = internal.DotProduct(v1, v2)
	require.NoError(err)
	assert.Equal(6, dp)

	v1 = []int{1, 2, 3}
	v2 = []int{1, 1, 1, 1}

	_, err = internal.DotProduct(v1, v2)
	require.Error(err)
}
