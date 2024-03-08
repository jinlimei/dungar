package cleaner

import "strings"

func (t *Token) parseWordForPlain() []Token {
	var (
		val = string(t.Value)
		out = make([]Token, 0)
	)

	if val == "&amp;" {
		out = append(out, Token{
			Type:     t.Type,
			StartPos: t.StartPos,
			EndPos:   t.EndPos,
			Value:    []rune("&"),
			Override: nil,
			hash:     nil,
		})
	} else if strings.HasSuffix(val, "%") {
		//fmt.Printf("'%s' HAS PERCENT SUFFIX WEW\n", val)
		// Let's check if this is a number w/ a percent
		isDigits := true
		for k := 0; k < len(t.Value)-1; k++ {
			if !isDigit(t.Value[k]) {
				isDigits = false
				break
			}
		}

		if isDigits {
			out = append(out, Token{
				Type:     TokenNumber,
				StartPos: t.StartPos,
				EndPos:   t.EndPos,
				Value:    t.Value,
				Override: nil,
				hash:     nil,
			})
		} else {
			out = append(out, *t)
		}

	} else {
		out = append(out, *t)
	}

	return out
}
