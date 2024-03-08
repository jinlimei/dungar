package accord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func parseDiscordMessage(msg *discordgo.Message) *core2.ParsedMessage {
	var (
		raw    = strings.TrimSpace(msg.Content)
		pieces = strings.Split(raw, " ")

		parsed = &core2.ParsedMessage{
			Tokens:    make([]core2.MessageToken, 0),
			Raw:       raw,
			Converted: raw,
		}

		pos  = 0
		pLen = len(pieces)
	)

	for ; pos < pLen; pos++ {
		more := parseDiscordPiece(pieces[pos])

		if len(more) > 0 {
			parsed.Tokens = append(parsed.Tokens, more...)
		}

		if pos+1 < pLen {
			parsed.AddToken(core2.MessageToken{
				Token: " ",
				Type:  core2.TokenSpace,
			})
		}
	}

	return parsed
}

func ptr(s string) *string {
	return &s
}

func parseDiscordPiece(piece string) []core2.MessageToken {
	//log.Printf("parseDiscordPiece='%s'", piece)
	var (
		//pLen  = len(piece)
		runes = []rune(piece)
		rLen  = len(runes)

		start = 0
		skip  = 0
		pos   = 0

		inBracket bool

		chr rune
		//nxt rune

		value     string
		tokenType core2.TokenType

		tokens = make([]core2.MessageToken, 0)
	)

	if rLen == 0 {
		return nil
	}

	get := func(loc int) rune {
		if loc >= rLen {
			return rune(0x0000)
		}

		return runes[loc]
	}

	find := func(str string, item rune, start int) int {
		runes := []rune(str)
		max := len(runes)

		for ; start < max; start++ {
			if runes[start] == item {
				return start
			}
		}

		return -1
	}

	for ; pos <= rLen; pos++ {
		chr = get(pos)
		//nxt = get(pos + 1)

		switch chr {
		case rune(0x0000):
			value = string(runes[start:pos])
			tokens = append(tokens, core2.MessageToken{
				Token: value,
				Type:  core2.TokenWord,
			})

		case '<':
			if pos > 0 {
				value = string(runes[0:pos])
				tokens = append(tokens, core2.MessageToken{
					Token: value,
					Type:  core2.TokenWord,
				})
			}

			inBracket = true
			start = pos
		case '>':
			if inBracket {
				inBracket = false
				value = string(runes[start : pos+1])
				//log.Printf("value='%s' value-1='%s'", value, string(value[1]))
				switch value[1] {
				case '@':
					tokenType = core2.TokenUserID
					skip = 2

					switch value[2] {
					case '!':
						skip = 3
						tokenType = core2.TokenUserID
					case '&':
						skip = 3
						tokenType = core2.TokenRoleID
					}

					tokens = append(tokens, core2.MessageToken{
						Token: value,
						Type:  tokenType,
						Value: ptr(value[skip : len(value)-1]),
					})
				case '#':
					tokens = append(tokens, core2.MessageToken{
						Token: value,
						Type:  core2.TokenChanID,
						Value: ptr(value[2 : len(value)-1]),
					})
				case ':':
					next := find(value, ':', 2)

					//log.Printf("next='%d'", next)

					tokens = append(tokens, core2.MessageToken{
						Token:    value,
						Type:     core2.TokenEmoticon,
						Override: ptr(value[1 : next+1]),
						Value:    ptr(value[1 : len(value)-1]),
					})
				}
			}

			start = pos + 1
			pos++
		}
	}

	return tokens
}

func (d *Driver) translateParsedMessage(serverID string, parsed *core2.ParsedMessage) string {
	var out string

	//session := d.Con.GetSession()

	for _, token := range parsed.Tokens {

		switch token.Type {
		case core2.TokenWord, core2.TokenSpace, core2.TokenURL:
			out += token.Token
		case core2.TokenUserID:
			name := d.GetUserName(*token.Value, serverID)

			if name != "" {
				out += "@" + name
			} else {
				out += "@" + *token.Value
			}

		case core2.TokenRoleID:
			name := d.GetRoleName(*token.Value, serverID)

			if name != "" {
				out += "&" + name
			} else {
				out += "&" + *token.Value
			}

		case core2.TokenChanID:
			name := d.GetChannelName(*token.Value, serverID)

			if name != "" {
				out += "#" + name
			} else {
				out += "#" + *token.Value
			}

		case core2.TokenEmoticon:
			out += *token.Override

		default:
			db.LogIssue(
				"unknown_discord_token",
				fmt.Sprintf("Unknown discord token '%s'", token.Type.String()),
				fmt.Sprintf("raw:\n%s\n", spew.Sdump(token)),
			)
		}
	}

	return out
}
