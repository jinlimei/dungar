package taut

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRubbieIssue(t *testing.T) {
	driver := initMockDriver()

	mockConnection.SetUser("U9LDWA6QL", "Dungar")

	parsed := parseSlackMessage("<@U9LDWA6QL|dungar> what do you think about butts?")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, "@Dungar what do you think about butts?", parsed.Converted)

	parsed = parseSlackMessage("<@U9999999|dungar> what do you think about butts?")
	parsed.Converted = driver.translateParsedMessage(parsed)
	assert.Equal(t, "@dungar what do you think about butts?", parsed.Converted)

	parsed = parseSlackMessage("<@U9999999> what do you think about butts?")
	parsed.Converted = driver.translateParsedMessage(parsed)
	assert.Equal(t, "@U9999999 what do you think about butts?", parsed.Converted)
}

func TestUrlParsing(t *testing.T) {
	driver := initMockDriver()

	parsed := parseSlackMessage("here is a test of <https://romeosquad.com|a url with text separate from the link>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, "here is a test of https://romeosquad.com", parsed.Converted)

	parsed = parseSlackMessage("here is|another test: <https://romeosquad.com/#|||#|a text| linke|e>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, "here is|another test: https://romeosquad.com/#", parsed.Converted)
}

func TestParseMessage(t *testing.T) {
	driver := initMockDriver()

	mockConnection.SetUser("U9LDWA6QL", "Dungar")

	parsed := parseSlackMessage("Hello, World")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, 3, len(parsed.Tokens))
	assert.Equal(t, "Hello, World", parsed.Converted)

	parsed = parseSlackMessage("Hello, <@user>")
	assert.Equal(t, 1, len(parsed.IDTokens()))

	parsed = parseSlackMessage("Hello, <@U9LDWA6QL|dungar>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	idTokens := parsed.IDTokens()
	assert.Equal(t, 1, len(idTokens))
	assert.Equal(t, 3, len(parsed.Tokens))
	assert.Equal(t, "@U9LDWA6QL", *idTokens[0].Value)
	assert.Equal(t, "dungar", *idTokens[0].Override)
	assert.Equal(t, "Hello, @Dungar", parsed.Converted)

	parsed = parseSlackMessage("Hello, <#thing|others>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	idTokens = parsed.IDTokens()
	assert.Equal(t, 1, len(idTokens))
	assert.Equal(t, "#thing", *idTokens[0].Value)
	assert.Equal(t, "others", *idTokens[0].Override)

	parsed = parseSlackMessage("Hello, <#thing>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	idTokens = parsed.IDTokens()
	assert.Equal(t, 1, len(idTokens))
	assert.Equal(t, "#thing", *idTokens[0].Value)
	assert.Equal(t, "", *idTokens[0].Override)

	parsed = parseSlackMessage("A neato URL-o: <https://www.google.com|Google URL>, check it out <@user> or maybe <#channel|Channel>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	idTokens = parsed.IDTokens()
	urlTokens := parsed.URLTokens()

	assert.Equal(t, 2, len(idTokens))
	assert.Equal(t, 1, len(urlTokens))

	assert.Equal(t, "https://www.google.com", *urlTokens[0].Value)
	assert.Equal(t, "Google URL", *urlTokens[0].Override)

	if len(idTokens) != 2 {
		assert.Fail(t, "idTokens must be 2")
	} else {
		assert.Equal(t, "@user", *idTokens[0].Value)
		assert.Equal(t, "#channel", *idTokens[1].Value)
	}
}

func TestNewParser(t *testing.T) {
	driver := initMockDriver()
	//clearMockConnection()

	parsed := parseSlackMessage(
		"Here is a complex URL: here<https://www.google.com/|isalink>togo<https://twitter.com|somewhere>" +
			"fun fun fun. <@U9LDWA6QL|dungar> how you doing. check out <#C7ZPDPR8W|programming> for more info",
	)

	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(
		t,
		"Here is a complex URL: hereisalinktogosomewherefun fun fun. @dungar how you doing. check out #programming for more info",
		parsed.Converted,
	)

	parsed = parseSlackMessage("except ❌✔🦀🤖📦 when it is, ofc, but")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, "except ❌✔🦀🤖📦 when it is, ofc, but", parsed.Converted)

	parsed = parseSlackMessage("<https://www.google.com/|a link>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, "https://www.google.com/", parsed.Converted)

	parsed = parseSlackMessage("a<https://www.google.com/|link>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, "alink", parsed.Converted)

	parsed = parseSlackMessage("<https://www.google.com/|foo>l")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, "fool", parsed.Converted)

	parsed = parseSlackMessage("a<https://www.google.com/>")
	parsed.Converted = driver.translateParsedMessage(parsed)

	assert.Equal(t, "a https://www.google.com/", parsed.Converted)
}
