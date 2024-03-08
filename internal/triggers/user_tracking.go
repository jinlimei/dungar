package triggers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// Handles last seen related functionality:
// - recording last seen users
// - !seen <users>
// - !top
func userTrackingHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if !msg.HasValidUser() {
		return core2.EmptyRsp()
	}

	user, err := svc.GetUser(msg.UserID, msg.ServerID)

	// Always record last seen, but do it after the command is processed.
	if err == nil {
		defer db.UpsertLastSeen(user, svc.Driver().DriverName())
	} else {
		log.Printf("failed to get user '%s' for server '%s': %v", msg.UserID, msg.ServerID, err)
	}

	txt := msg.Contents

	// Handle last seen
	if strings.HasPrefix(txt, "!seen ") && len(txt) > 6 {
		split := strings.SplitN(txt, " ", 2)
		res, err := db.GetLastSeen(split[1], msg.ServerID)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return core2.PrefixedSingleRsp("could not find user")
			}

			return core2.PrefixedSingleRsp("A random unknown error occurred. Babies, maybe.")
		}

		if res == nil {
			return core2.PrefixedSingleRsp("could not find user")
		}

		return core2.PrefixedSingleRsp(fmt.Sprintf(
			"%s was last seen: %s",
			split[1],
			res.LastSeen.Format("2006-01-02 03:04:05 pm MST"),
		))
	}

	if strings.HasPrefix(txt, "!top") {
		posters := db.GetTopPosters(msg.ServerID)

		if len(posters) == 0 {
			return core2.MakeSingleRsp("No top posters, lol y'all suck")
		}

		out := ""
		for pos := 0; pos < len(posters); pos++ {
			poster := posters[pos]
			out += fmt.Sprintf("%02d: %s with %d posts\n", pos+1, poster.Nick, poster.MessageCount)
		}

		return core2.MakeSingleRsp(out)
	}

	return core2.EmptyRsp()
}
