package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeIsZero(t *testing.T) {
	assert.Equal(t, 1, 1)
	ts := time.Unix(0, 0)
	assert.True(t, TimeIsZero(ts))
	assert.False(t, ts.IsZero())

	ts = ts.Add(1)
	assert.False(t, TimeIsZero(ts))
	assert.False(t, ts.IsZero())

	ts = time.Time{}
	assert.True(t, TimeIsZero(ts))
	assert.True(t, ts.IsZero())

	ts = ts.Add(1)
	assert.False(t, TimeIsZero(ts))
	assert.False(t, ts.IsZero())
}
