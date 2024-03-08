package accord

import (
	"net/url"
	"strings"
)

type parsedEmoji struct {
	id       string
	name     string
	original string
}

func (pe *parsedEmoji) toReaction() string {
	if pe.id != "" {
		return pe.name + ":" + pe.id
	}

	return pe.name
}

func (d *Driver) isEmoji(str string) bool {
	return str[0] == ':' && str[len(str)-1] == ':'
}

func (d *Driver) translateMessageEmoji(message, guildID string) string {
	words := strings.Split(message, " ")

	guild := d.getOrMakeGuild(guildID)

	for idx, word := range words {
		if !d.isEmoji(word) {
			continue
		}

		parsed := d.parseDiscordEmoji(word)
		// We now need to see if we can use this emoji or not. If not, we'll just use :pensive:

		emoji, ok := guild.emojiCache[parsed.name]
		if !ok {
			// lol
			words[idx] = ":pensive:"
		} else {
			words[idx] = emoji.MessageFormat()
		}
	}

	return strings.Join(words, " ")
}

func (d *Driver) translateReactionEmoji(emoji, guildID string) string {
	// Regardless of the emoji's ID incoming we'll need to only use our
	// local servers emoji, so we'll just have to deal.

	parsed := d.parseDiscordEmoji(emoji)
	guild := d.getOrMakeGuild(guildID)

	guildEmoji, ok := guild.emojiCache[parsed.name]

	if !ok {
		// Well... we couldn't find it! We'll try yolo'ing it anyway
		return url.QueryEscape(parsed.toReaction())
	}

	return url.QueryEscape(guildEmoji.APIName())
}

func (d *Driver) parseDiscordEmoji(emoji string) parsedEmoji {
	var (
		name string
		id   string

		runes = []rune(emoji)
		lim   = len(runes)
		chr   rune
		pos   int
	)

	pull := func(start, end int) string {
		out := make([]rune, end-start)

		for i, k := start, 0; i < end; i, k = i+1, k+1 {
			out[k] = runes[i]
		}

		return string(out)
	}

	for k := 0; k < lim; k++ {
		chr = runes[k]

		switch chr {
		case ':':
			// Skip if it's the very first : because we don't caaaaaaaaare
			if k == 0 {
				k++
				pos++
				break
			}

			name = pull(pos, k)
			pos = k + 1
		}
	}

	if name != "" {
		id = pull(pos, lim)
	} else {
		name = pull(pos, lim)
	}

	//fmt.Printf("original=%s, id=%s, name=%s\n", emoji, id, name)

	return parsedEmoji{
		id:       id,
		name:     name,
		original: emoji,
	}
}
