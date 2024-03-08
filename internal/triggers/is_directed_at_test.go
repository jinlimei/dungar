package triggers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestExtractUsernameTarget(t *testing.T) {
	strs := []string{
		"@dungar h",
		"dungar: h",
	}

	for _, str := range strs {
		assert.Equal(t, "dungar", extractUserNameTarget(str, questionStartRegex))
	}
}

func makeIncMessage(msg, name string) *core2.IncomingMessage {
	contents := fmt.Sprintf(msg, name)

	return &core2.IncomingMessage{
		ID:             "",
		UserID:         "",
		ChannelID:      "",
		SubChannelID:   "",
		IsSubMessage:   false,
		Contents:       contents,
		LowerContents:  strings.ToLower(contents),
		ParsedContents: nil,
		Type:           0,
		Attachments:    nil,
	}
}
func TestIsDirectedAtDungar(t *testing.T) {
	svc := initMockServices()

	inf := svc.GetBotUser()

	msg := makeIncMessage("@%s sup", inf.Name)
	assert.Truef(t, isDirectedAtDungar(svc, msg), "got false for msg '%v'", msg)
	msg = makeIncMessage("@%s: sup", inf.Name)
	assert.Truef(t, isDirectedAtDungar(svc, msg), "got false for msg '%v'", msg)
	msg = makeIncMessage("%s sup", inf.Name)
	assert.Truef(t, isDirectedAtDungar(svc, msg), "got false for msg '%v'", msg)
	msg = makeIncMessage("%s: sup", inf.Name)
	assert.Truef(t, isDirectedAtDungar(svc, msg), "got false for msg '%v'", msg)

	msg = makeIncMessage("hello @%s", inf.Name)
	assert.Falsef(t, isDirectedAtDungar(svc, msg), "got true for msg '%v'", msg)
	msg = makeIncMessage("!%s sup", inf.Name)
	assert.Falsef(t, isDirectedAtDungar(svc, msg), "got true for msg '%v'", msg)
}

func TestIsDirectedAtDungarRandom(t *testing.T) {
	random.UseTestSeed()
	svc := initMockServices()

	inf := svc.GetBotUser()

	msg := makeIncMessage("@%s sup", inf.Name)
	assert.Truef(t, isDirectedAtDungarRandom(svc, msg, 0.10), "got false for msg '%v'", msg)

	msg = makeIncMessage("@%s get it together fam", "debra")
	k := 0
	anyTrue := false
	anyFalse := false
	for ; k < 100; k++ {
		if isDirectedAtDungarRandom(svc, msg, 0.10) {
			anyTrue = true
			break
		} else {
			anyFalse = true
		}
	}

	assert.Truef(t, anyTrue, "failed to find msg '%v' matching after %d iterations",
		msg, k)
	assert.Truef(t, anyFalse, "failed to find msg '%v' NOT matching dungar",
		msg)
}

func TestIsMentioningDungar(t *testing.T) {
	svc := initMockServices()
	inf := svc.GetBotUser()

	msg := makeIncMessage("anyway so why is %s so mad?", inf.Name)
	assert.Truef(t, isMentioningDungar(svc, msg), "msg '%v' is not mentioning dungar", msg)
	msg = makeIncMessage("whats up @%s?", inf.Name)
	assert.Truef(t, isMentioningDungar(svc, msg), "msg '%v' is not mentioning dungar", msg)

	msg = makeIncMessage("hey friends how are you doing?", "")
	assert.Falsef(t, isMentioningDungar(svc, msg), "msg '%v' is mentioning dungar", msg)
}

func TestIsMentioningDungarRandom(t *testing.T) {
	random.UseTestSeed()
	svc := initMockServices()
	inf := svc.GetBotUser()

	msg := makeIncMessage("anyway so why is %s so mad?", inf.Name)
	assert.Truef(t, isMentioningDungarRandom(svc, msg, 0.10), "msg '%v' is not mentioning dungar", msg)
	msg = makeIncMessage("whats up @%s?", inf.Name)
	assert.Truef(t, isMentioningDungarRandom(svc, msg, 0.10), "msg '%v' is not mentioning dungar", msg)
	assert.Truef(t, isMentioningDungarRandom(svc, msg, 0.00), "msg '%v' is not mentioning dungar", msg)

	msg = makeIncMessage("hey friends how are you doing?", "")
	assert.Truef(t, isMentioningDungarRandom(svc, msg, 1.00), "msg '%v' is not mentioning dungar", msg)

	var (
		pos      = 0
		anyFalse = false
		anyTrue  = false
	)

	msg = makeIncMessage("hey friends how are you doing?", "")
	for ; pos < 100; pos++ {
		if isMentioningDungarRandom(svc, msg, 0.10) {
			anyTrue = true
			break
		} else {
			anyFalse = true
		}
	}

	assert.Truef(t, anyTrue, "failed to find msg '%v' matching after %d iterations",
		msg, pos)
	assert.Truef(t, anyFalse, "failed to find msg '%v' NOT matching dungar after %d iterations",
		msg, pos)
}

func TestCheckDirectedPrefix(t *testing.T) {
	assert.False(t, checkDirectedPrefix("> banana", "dungar"))
	assert.False(t, checkDirectedPrefix(">banana", "dungar"))
	assert.False(t, checkDirectedPrefix(" >  banana", "dungar"))
}

func BenchmarkDirectedAtDungar(b *testing.B) {
	svc := initMockServices()
	inf := svc.GetBotUser()
	msg := makeIncMessage("@%s whats up friend?", inf.Name)

	for k := 0; k < b.N; k++ {
		isDirectedAtDungar(svc, msg)
	}
}

func BenchmarkMentioningDungar(b *testing.B) {
	svc := initMockServices()
	usr := svc.GetBotUser()
	msg := makeIncMessage("@%s whats up friend?", usr.Name)

	for k := 0; k < b.N; k++ {
		isMentioningDungar(svc, msg)
	}
}

func isAround(num int, around int, spread int) bool {
	return (around-spread) < num && (around+spread) > num
}
