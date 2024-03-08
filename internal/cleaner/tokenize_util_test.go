package cleaner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRuneRange(t *testing.T) {
	a := []rune{'a', 'b', 'c'}
	assert.NotNil(t, getRuneRange(a, 0, 3))
	assert.Nil(t, getRuneRange(a, 0, 4))
}

func BenchmarkRuneRange(b *testing.B) {
	var working = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}

	for i := 0; i < b.N; i++ {
		getRuneRange(working, -1, 0)
		getRuneRange(working, 0, -1)
		getRuneRange(working, 2, 4)
	}
}

func TestPeekUntil(t *testing.T) {
	assert.Equal(t, 1, 1)

	data := []rune("hello world")

	v := peekUntil(data, 0, 'e')
	assert.Equal(t, 1, v)
	v = inlinePeekUntil(data, 0, 'e')
	assert.Equal(t, 1, v)

	v = peekUntil(data, 2, 'l')
	assert.Equal(t, 2, v)
	v = inlinePeekUntil(data, 2, 'l')
	assert.Equal(t, 2, v)

	v = peekUntil(data, 0, '<')
	assert.Equal(t, -1, v)
	v = inlinePeekUntil(data, 0, '<')
	assert.Equal(t, -1, v)
}

func BenchmarkPeekUntil(b *testing.B) {
	data := []rune("hello world")

	for k := 0; k < b.N; k++ {
		peekUntil(data, 0, ' ')
	}
}

func BenchmarkInlinePeekUntil(b *testing.B) {
	data := []rune("hello world")

	for k := 0; k < b.N; k++ {
		inlinePeekUntil(data, 0, ' ')
	}
}
