package cleaner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	messageList = []string{
		`so this one's like "yo gimme dis" https://i.arq.dev/U78GfVt.png`,

		`Hey friends :sun: what's going on? I just went to ` +
			`https://www.google.com or maybe <@thing> or maybe @thing might like it idk ` +
			`but also check out <https://www.reddit.com> yeah ok cool bye :sun:`,

		`@dungar dungar or executive order for a month as long as your regular pop like ten days`,
		"``` ```",
		"```  ```",
		"```   ```",
		"`0.9 6.3 8271900 1169016`",
		"idk how to get raw `` through tho. any of those in your message get encoded to `&lgt;` from slack. only making urls or @'ing people makes the actual chars",
		"2 of these in it <https://ark.intel.com/content/www/us/en/ark/products/120495/intel-xeon-gold-6154-processor-24-75m-cache-3-00-ghz.html>",
		":abc::",
		":hello::world:",
		":hello: :world:",
		"here :sun: and an example: foo",
		"the word 'example:' should not be considered an imaginary emote",
		"uh hello: :rrrr::",
		"uh hello: :rrrr:: ::: :",
		"Hey there <:durrhuREEE:546076030666801177>",
		"<@&766488443248836609> this is pinging the example role",
		"free market conservatives: we're going to de-regulate the markets and let them solve their own problems!\nalso free market conservatives: the market is fixing itself in a way I don't like! let's regulate them",
		":musical_note: we're gonna lose the house :musical_note:",
		"been a good morning... realized my previous landlord is in violation of rental laws and is withholding our security deposit past the duration he is legally allowed to which means we're entitled to like 3x the amount withheld, plus $100, plus attorney fees  :notbadobama:",
	}
)

func TestMessageList(t *testing.T) {
	assert.True(t, len(messageList) > 0)
}
