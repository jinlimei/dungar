package db

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// UserTracking is the main struct for User Tracking
type UserTracking struct {
	UserID         string    `json:"user_id"`
	ServerID       string    `json:"server_id"`
	ProtocolDriver string    `json:"protocol_driver"`
	Nick           string    `json:"nick"`
	MessageCount   int       `json:"message_count"`
	CustomData     []byte    `json:"custom_data"`
	FirstSeen      time.Time `json:"first_seen"`
	LastSeen       time.Time `json:"last_seen"`
}

// CustomDataJSON returns a json.RawMessage for the UserTracking.CustomData member
func (ut *UserTracking) CustomDataJSON() json.RawMessage {
	if ut.CustomData == nil || len(ut.CustomData) == 0 {
		return []byte{}
	}

	return []byte(ut.CustomData)
}

// TopPoster is the struct for returning the top user trackerings.
type TopPoster struct {
	UserID       string
	ServerID     string
	Nick         string
	MessageCount int
}

// GetUserFromUserTracking attempts to grab a user from tracking data
// bool is whether the user was found.
func GetUserFromUserTracking(userID, serverID string) (UserTracking, bool) {
	if utils.InTestEnv() {
		return UserTracking{
			UserID:         userID,
			ServerID:       serverID,
			ProtocolDriver: "mock",
			Nick:           "bob",
			MessageCount:   100,
			CustomData:     nil,
			FirstSeen:      time.Unix(100000, 0),
			LastSeen:       time.Now(),
		}, true
	}

	qry := `
		SELECT protocol_driver, nick, message_count, first_seen, last_seen
		FROM user_tracking_new
		WHERE user_id = $1
		  AND server_id = $2
	`

	row := ConQueryRow(qry, userID, serverID)
	if row == nil {
		return UserTracking{}, false
	}

	var (
		protocolDriver string
		nick           string
		messageCount   int
		firstSeen      time.Time
		lastSeen       time.Time
	)

	err := row.Scan(
		&protocolDriver,
		&nick,
		&messageCount,
		&firstSeen,
		&lastSeen,
	)

	if err != nil {
		log.Printf("ERROR: failed to retrieve row from user tracking: %v", err)
		return UserTracking{}, false
	}

	return UserTracking{
		UserID:         userID,
		ServerID:       serverID,
		ProtocolDriver: protocolDriver,
		Nick:           nick,
		MessageCount:   messageCount,
		FirstSeen:      firstSeen,
		LastSeen:       lastSeen,
	}, true
}

// UpdateUserFromUserTracking will update the last seen, nick, and message count
// for users when an event is observed from them
func UpdateUserFromUserTracking(tracking UserTracking) error {
	if utils.InTestEnv() {
		return nil
	}

	qry := `
		UPDATE user_tracking_new
		SET nick = $3,
		    message_count = $4,
		    last_seen = CURRENT_TIMESTAMP
		WHERE user_id = $1
		  AND server_id = $2
	`

	//log.Printf("UpdateUserFromUserTracking: %+v", tracking)

	_, err := ConExec(
		qry,
		tracking.UserID,       // 1
		tracking.ServerID,     // 2
		tracking.Nick,         // 3
		tracking.MessageCount, // 4
	)

	return err
}

// CreateUserForUserTracking will do as labelled on the tin
func CreateUserForUserTracking(user core2.User, protocolDriver string) error {
	if utils.InTestEnv() {
		return nil
	}

	qry := `
		INSERT INTO user_tracking_new(user_id, server_id, protocol_driver, nick, custom_data, first_seen, last_seen)
		VALUES(CAST($1 AS text), CAST($2 as text), CAST($3 as text), CAST($4 as text), null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
  `

	_, err := ConExec(
		qry,
		user.ID,        // 1
		user.ServerID,  // 2
		protocolDriver, // 3
		user.Name,      // 4
	)

	return err
}

// UpsertLastSeen does an insert/on-conflict-update to handle line counts for a user based off of their unique_id
func UpsertLastSeen(user core2.User, protocolDriver string) {
	if utils.InTestEnv() {
		return
	}

	var (
		err    error
		action string
	)

	tracking, ok := GetUserFromUserTracking(user.ID, user.ServerID)
	if !ok {
		action = "insert"
		err = CreateUserForUserTracking(user, protocolDriver)
	} else {
		tracking.MessageCount++
		tracking.Nick = user.Name

		action = "update"
		err = UpdateUserFromUserTracking(tracking)
	}

	if err != nil {
		log.Printf("ERROR: Failed to action='%s' user '%+v': %v", action, user, err)
	}
}

// GetActiveUsers returns all active users (active means message in at least the last month)
func GetActiveUsers(serverID string) []UserTracking {
	if utils.InTestEnv() {
		return []UserTracking{
			{"", "", "", "", 0, nil, time.Now(), time.Now()},
			{"", "", "", "", 0, nil, time.Now(), time.Now()},
		}
	}

	query := `
		SELECT user_id, nick, message_count, custom_data, first_seen, last_seen
		FROM user_tracking_new
		WHERE last_seen > $1
		  AND server_id = $2
		ORDER BY last_seen DESC
	`

	prev := time.Now().AddDate(0, -1, 0)
	stmt := ConMustQuery(query, prev, serverID)

	defer stmt.Close()

	var (
		tmp    UserTracking
		output = make([]UserTracking, 0)
	)

	for stmt.Next() {
		tmp = UserTracking{}

		if err := stmt.Scan(&tmp.UserID, &tmp.Nick, &tmp.MessageCount, &tmp.FirstSeen, &tmp.LastSeen); err != nil {
			utils.HaltingError("GetActiveUsers", err)
		}

		output = append(output, tmp)
	}

	return output
}

// GetTopPosters returns the list of Top Posters
func GetTopPosters(serverID string) []TopPoster {
	if utils.InTestEnv() {
		return []TopPoster{
			{"", "", "foo", 2},
			{"", "", "abc", 0},
		}
	}

	query := `
		SELECT user_id, nick, message_count
		FROM user_tracking_new
		WHERE last_seen > $1
		  AND server_id = $2
		ORDER BY message_count DESC
		LIMIT 11
	`

	var (
		tmp    TopPoster
		output = make([]TopPoster, 0)
	)

	prev := time.Now().AddDate(0, -1, 0)
	stmt := ConMustQuery(query, prev, serverID)

	defer stmt.Close()

	for stmt.Next() {
		tmp = TopPoster{}

		if err := stmt.Scan(&tmp.UserID, &tmp.Nick, &tmp.MessageCount); err != nil {
			utils.HaltingError("GetTopPosters", err)
		}

		output = append(output, tmp)
	}

	return output
}

// GetLastSeen returns the last seen user as apart of a struct if they've existed, otherwise not at all.
func GetLastSeen(nick, serverID string) (*UserTracking, error) {
	if utils.InTestEnv() {
		if nick == "v" {
			return &UserTracking{
				UserID:       "",
				Nick:         nick,
				MessageCount: 1,
				FirstSeen:    time.Unix(0, 0),
				LastSeen:     time.Now(),
			}, nil
		}

		return nil, nil
	}

	nick = strings.ToLower(nick)

	query := `
		SELECT user_id, nick, message_count, last_seen, first_seen
		FROM user_tracking_new
		WHERE server_id = $1
		  AND LOWER(nick) = $2
	`

	row := ConQueryRow(query, serverID, nick)

	if row == nil {
		return nil, fmt.Errorf("row is nil in GetLastSeen")
	}

	out := &UserTracking{}

	if err := row.Scan(&out.UserID, &out.Nick, &out.MessageCount, &out.LastSeen, &out.FirstSeen); err != nil {
		return nil, err
	}

	return out, nil
}
