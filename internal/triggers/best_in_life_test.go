package triggers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBestInLifeHandler(t *testing.T) {
	res := bestInLifeHandler("", "")

	assert.True(t, res != "")
}
