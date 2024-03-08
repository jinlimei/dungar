package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmoticon(t *testing.T) {
	assert.False(t, IsEmoticon("So:"))
	assert.False(t, IsEmoticon(":well then:"))
	assert.True(t, IsEmoticon(":well_then:"))
	// TODO should this really be valid?
	assert.True(t, IsEmoticon(":well::then:"))
}

func TestEndsWithPunctuation(t *testing.T) {
	assert.True(t, EndsWithPunctuation("!"))
	assert.True(t, EndsWithPunctuation("hey!"))
	assert.False(t, EndsWithPunctuation("hey"))
	assert.False(t, EndsWithPunctuation(""))
}

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice("foo", []string{"foo", "banana"}))
	assert.True(t, StringInSlice("aaa", []string{"b", "c", "d", "r", "w", "aaa"}))
	assert.False(t, StringInSlice("haha", []string{"foo", "banana"}))
	assert.False(t, StringInSlice("foo", []string{}))
	assert.False(t, StringInSlice("foo", nil))
}

func TestIsUrl(t *testing.T) {
	assert.False(t, IsURL("hello world"))
	assert.False(t, IsURL("hx"))
	assert.False(t, IsURL("1234567"))
	assert.False(t, IsURL("12345678"))

	assert.True(t, IsURL("https://google"))
	assert.True(t, IsURL("google.com"))
	assert.True(t, IsURL("world.environment"))
	assert.True(t, IsURL("https://chris.allan.codes"))
	assert.True(t, IsURL("https://is.url"))
	assert.True(t, IsURL("https://bit.ly/hello+world(this)"))
}

func TestCleanSpaces(t *testing.T) {
	assert.Equal(t, "hello world", CleanSpaces(" hello  world "))
	assert.Equal(t, "there is a wonderful amount of space",
		CleanSpaces("there         is           a            wonderful                amount of space"))
}

func TestNormalize(t *testing.T) {
	assert.Equal(t, "hello", Normalize("HELLo"))
	assert.Equal(t, "hello", Normalize("HELLo!!"))
	assert.Equal(t, "!", Normalize("!"))
	assert.Equal(t, "a", Normalize("a!"))
}

func TestTrimPunctuation(t *testing.T) {
	assert.Equal(t, "fred", TrimPunctuation("@fred"))
	assert.Equal(t, "fred", TrimPunctuation("fred!"))
}

func TestCoalesceStr(t *testing.T) {
	ret, ok := CoalesceStr("foo")

	assert.True(t, ok)
	assert.Equal(t, "foo", ret)

	ret, ok = CoalesceStr("", "foo")

	assert.True(t, ok)
	assert.Equal(t, "foo", ret)

	ret, ok = CoalesceStr("", "", "foo", "")

	assert.True(t, ok)
	assert.Equal(t, "foo", ret)

	ret, ok = CoalesceStr("", "", "biz", "baz", "buz")

	assert.True(t, ok)
	assert.Equal(t, "biz", ret)

	ret, ok = CoalesceStr("", "", "")

	assert.False(t, ok)
	assert.Equal(t, "", ret)

	ret, ok = CoalesceStr("")

	assert.False(t, ok)
	assert.Equal(t, "", ret)

	ret, ok = CoalesceStr()

	assert.False(t, ok)
	assert.Equal(t, "", ret)
}

func TestTitleCase(t *testing.T) {
	assert.Equal(t, "Banana", TitleCase("banana"))
	assert.Equal(t, "Banana BaNaNa", TitleCase("banana baNaNa"))
}
