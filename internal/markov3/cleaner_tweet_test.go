package markov3

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanTweet(t *testing.T) {
	tw := "Will be on @meetthepress this morning. Check times. @chucktodd  @NBCNews"
	cl := "Will be on @user this morning. Check times. @user @user"

	res := CleanTweet(tw, true)
	assert.Equal(t, strings.ToLower(cl), res.Cleaned)
	assert.Equal(t, []string{"@meetthepress", "@chucktodd", "@NBCNews"}, res.Mentions)
}
