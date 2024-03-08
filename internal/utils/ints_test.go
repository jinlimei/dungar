package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinInt64s(t *testing.T) {
	assert.Equal(t, []int64{1,2,3}, JoinInt64s([]int64{1},[]int64{2,3}))
}

func TestReverseInts(t *testing.T) {
	source := []int64{8, 6, 7, 5, 3, 0, 9}

	assert.Equal(
		t,
		[]int64{9, 0, 3, 5, 7, 6, 8},
		ReverseInt64s(source),
	)
}
