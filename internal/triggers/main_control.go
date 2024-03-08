package triggers

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func mainControlHandler(msg *core2.IncomingMessage) []*core2.Response {
	if !isAdmin(msg.UserID) {
		return core2.EmptyRsp()
	}

	if len(msg.Contents) > 9 && (strings.HasPrefix(msg.Contents, "!disable ") ||
		strings.HasPrefix(msg.Contents, "!change ")) {

		pieces := strings.Split(msg.Contents, " ")

		var (
			cmd = pieces[0]
			val = pieces[1]
			nxt = ""

			prev   float64
			next   float64
			err    error
			exists bool
			valid  bool
		)

		if len(pieces) > 2 {
			nxt = pieces[2]
		}

		if cmd == "!change" && nxt != "" {
			next, err = strconv.ParseFloat(nxt, 64)

			if err != nil {
				return core2.MakeSingleRsp(fmt.Sprintf("Could not figure out next value '%s': %v",
					nxt, err))
			}

			valid = true
		} else if cmd == "!disable" {
			valid = true
			next = 0
		}

		if !valid {
			return core2.MakeSingleRsp("IDK What you want")
		}

		prev, exists = masterChanceList[val]

		if !exists {
			return core2.MakeSingleRsp(fmt.Sprintf("Could not find value ID '%s'", val))
		}

		masterChanceList[val] = next

		return core2.MakeSingleRsp(fmt.Sprintf(
			"Changed val '%s' to be from '%.4f' to '%0.4f'",
			val,
			prev,
			next,
		))
	}

	return core2.EmptyRsp()
}
