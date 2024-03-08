package taut

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// Break down the message and look for @user and #channel references
// to convert them into a set of tokens
func parseSlackMessage(rawMsg string) *core2.ParsedMessage {
	rawMsg = strings.TrimSpace(rawMsg)

	parsed := &core2.ParsedMessage{
		Tokens:    make([]core2.MessageToken, 0),
		Raw:       rawMsg,
		Converted: rawMsg,
	}

	var (
		inBracket   = false
		inSecondary = false

		openPos      = 0
		bracketStart = 0
		splitStart   = 0

		pieceF = ""
		pieceL = ""

		msgRunes = []rune(rawMsg)
		pos      = 0
		mLen     = len(msgRunes)
	)

	get := func(pos int) rune {
		if pos >= mLen {
			return rune(0x00)
		}

		return msgRunes[pos]
	}

	getRunes := func(start, end int) string {
		out := make([]rune, 0)

		for start < end {
			out = append(out, msgRunes[start])
			start++
		}

		return string(out)
	}

	var chr rune

	for true {
		chr = get(pos)
		pos++

		switch chr {
		case rune(0x00):
			tmp := getRunes(openPos, pos-1)

			if tmp != "" {
				parsed.Tokens = append(parsed.Tokens, core2.MessageToken{
					Token:    tmp,
					Type:     core2.TokenWord,
					Value:    nil,
					Override: nil,
				})
			}

		case ' ':
			if inBracket {
				continue
			}

			tmp := getRunes(openPos, pos-1)

			if tmp != "" {
				parsed.Tokens = append(parsed.Tokens, core2.MessageToken{
					Token:    tmp,
					Type:     core2.TokenWord,
					Value:    nil,
					Override: nil,
				})
			}

			parsed.Tokens = append(parsed.Tokens, core2.MessageToken{
				Token:    " ",
				Type:     core2.TokenSpace,
				Value:    nil,
				Override: nil,
			})

			openPos = pos

		case '<':
			tmp := getRunes(openPos, pos-1)
			if tmp != "" {
				parsed.Tokens = append(parsed.Tokens, core2.MessageToken{
					Token:    tmp,
					Type:     core2.TokenWord,
					Value:    nil,
					Override: nil,
				})
			}

			bracketStart = pos
			splitStart = pos
			inBracket = true

		case '|':
			if inBracket && !inSecondary {
				pieceF = getRunes(bracketStart, pos-1)
				splitStart = pos
				inSecondary = true
				continue
			}
		case '>':
			inBracket = false
			inSecondary = false
			pieceL = getRunes(splitStart, pos-1)

			if pieceF == "" {
				pieceF = pieceL
				pieceL = ""
			}

			t := core2.TokenWord

			if utils.IsURL(pieceF) {
				t = core2.TokenURL
			} else if pieceF[0] == '@' {
				t = core2.TokenUserID
			} else if pieceF[0] == '#' {
				t = core2.TokenChanID
			}

			value := pieceF
			override := pieceL

			parsed.Tokens = append(parsed.Tokens, core2.MessageToken{
				Token:    getRunes(bracketStart-1, pos),
				Type:     t,
				Value:    &value,
				Override: &override,
			})

			openPos = pos
			pieceF = ""
			pieceL = ""
		}

		if chr == rune(0x00) {
			break
		}
	}

	return parsed
}

func getPtr(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// translateParsedMessage will bake in the "correct" stuff we're expecting from
// things like <@user> references or <#channel> references to raw "@user in #channel" text
func (d *Driver) translateParsedMessage(parsed *core2.ParsedMessage) string {
	var (
		out  string
		prev *core2.MessageToken
		next *core2.MessageToken
		tLen = len(parsed.Tokens)
	)

	for pos, token := range parsed.Tokens {
		if pos > 0 {
			prev = &parsed.Tokens[pos-1]
		} else {
			prev = nil
		}

		if (pos + 1) < tLen {
			next = &parsed.Tokens[pos+1]
		} else {
			next = nil
		}

		switch token.Type {

		case core2.TokenWord, core2.TokenSpace:
			out += token.Token

		case core2.TokenURL:
			if prev == nil && (next == nil || next.Type == core2.TokenSpace) {
				out += *token.Value
			} else if prev != nil && prev.Type == core2.TokenSpace &&
				(next == nil || next.Type == core2.TokenSpace) {
				out += *token.Value
			} else if token.Override == nil || *token.Override == "" {

				if prev != nil {
					out += " "
				}

				out += *token.Value

				if next != nil {
					out += " "
				}

			} else {
				out += *token.Override
			}
		case core2.TokenUserID:
			name := ""

			if d.Con != nil && d.Con.IsConnected() {
				name = d.GetUserName((*token.Value)[1:], d.teamID)
			}

			if name == "" && token.Override != nil && *token.Override != "" {
				name = *token.Override
			}

			if name != "" {
				out += "@" + name
			} else {
				out += *token.Value
			}

			//fmt.Printf(
			//	"NAME: '%s' FROM TOKEN token='%s', type='%s', value='%s', override='%s'\n",
			//	name, token.Token, token.Type.String(), getPtr(token.Value), getPtr(token.Override),
			//)

		case core2.TokenChanID:
			name := ""

			if d.Con != nil && d.Con.IsConnected() {
				name = d.GetChannelName((*token.Value)[1:], d.teamID)
			}

			if name == "" && token.Override != nil && *token.Override != "" {
				name = *token.Override
			}

			if name != "" {
				out += "#" + name
			} else {
				out += *token.Value
			}

			//fmt.Printf(
			//	"NAME: '%s' FROM TOKEN token='%s', type='%s', value='%s', override='%s'\n",
			//	name, token.Token, token.Type.String(), getPtr(token.Value), getPtr(token.Override),
			//)

		default:
			db.LogIssue("unknown_slack_token",
				fmt.Sprintf("unknown slack token '%v'", token.Type),
				fmt.Sprintf("raw: %s\n", spew.Sdump(token)))
		}
	}

	return out
}
