package markov3

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
)

type rawMessageLearningPlan struct {
	query   string
	args    []any
	variant cleaner.TokenVariant
	source  string
}

// LearnFromRawMessagesM1 is a method for learning from the legacy raw messages
// table built from the markov database of it
func (m *Markov) LearnFromRawMessagesM1() int {
	if utils.InTestEnv() {
		return -1
	}

	if m.learnedFromLegacy {
		return -1
	}

	db.EnsureDatabaseConnection()

	qry := `
		SELECT sentence_id, message
		FROM raw_messages_m1
	`

	var (
		res     *sql.Rows
		id      uint64
		msg     string
		learned = 0
	)

	res = db.ConMustQuery(qry)

	for res.Next() {
		if err := res.Scan(&id, &msg); err != nil {
			log.Printf("Failed to scan: %v\n", err)
			break
		}

		spc := spaceCount(msg)

		if spc > 0 {
			m.LearnString(msg, cleaner.VariantXMPP)
			learned++
		}
	}

	m.learnedFromLegacy = true
	return learned
}

func spaceCount(s string) int {
	var (
		rs = []rune(s)
		cn = 0
	)

	for k := 0; k < len(rs); k++ {
		if rs[k] == ' ' {
			cn++
		}
	}

	return cn
}

// LearnFromRawMessages will pull in the `raw_messages`
// table from PGSQL and push them into the internal
// markov state
func (m *Markov) LearnFromRawMessages() int {
	if utils.InTestEnv() {
		return -1
	}

	db.EnsureDatabaseConnection()
	plan := buildRawMessageLearningPlan(m.lastRawMessageID)

	log.Printf(
		"[LearnFromRawMessages] source=%s lastRawMessageID=%d args=%d",
		plan.source,
		m.lastRawMessageID,
		len(plan.args),
	)

	res := db.ConMustQuery(plan.query, plan.args...)

	var (
		id      uint64
		msg     string
		learned = 0
	)

	for res.Next() {
		if err := res.Scan(&id, &msg); err != nil {
			log.Printf("Failed to scan: %v\n", err)
			break
		}

		if len(msg) > 1 {
			m.LearnString(msg, plan.variant)
			m.lastRawMessageID = id
			learned++
		}
	}

	return learned
}

// LearnFromDiscordRawMessagesByAuthor builds a temporary model from Discord raw
// messages scoped to one author, using the same time/channel/blocklist filters
// as the normal Discord learning path.
func (m *Markov) LearnFromDiscordRawMessagesByAuthor(serverID, authorUserID string, allowed map[string]struct{}) (int, error) {
	if utils.InTestEnv() {
		return -1, nil
	}

	if strings.TrimSpace(serverID) == "" || strings.TrimSpace(authorUserID) == "" {
		return 0, nil
	}

	db.EnsureDatabaseConnection()
	plan := buildDiscordAuthorRawMessageLearningPlan(serverID, authorUserID, allowed)

	log.Printf(
		"[LearnFromDiscordRawMessagesByAuthor] serverID=%s authorUserID=%s args=%d",
		serverID,
		authorUserID,
		len(plan.args),
	)

	res, err := db.ConQuery(plan.query, plan.args...)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	var (
		id      uint64
		msg     string
		learned = 0
	)

	for res.Next() {
		if err := res.Scan(&id, &msg); err != nil {
			return learned, err
		}

		if len(msg) > 1 {
			m.LearnString(msg, plan.variant)
			m.lastRawMessageID = id
			learned++
		}
	}

	if err := res.Err(); err != nil {
		return learned, err
	}

	return learned, nil
}

func buildRawMessageLearningPlan(lastRawMessageID uint64) rawMessageLearningPlan {
	if strings.EqualFold(utils.ProtocolMode(), "discord") {
		return buildDiscordRawMessageLearningPlan(
			lastRawMessageID,
			utils.DiscordAllowedLearningChannelIDs(),
		)
	}

	return buildSlackRawMessageLearningPlan(lastRawMessageID)
}

func buildSlackRawMessageLearningPlan(lastRawMessageID uint64) rawMessageLearningPlan {
	plan := rawMessageLearningPlan{
		query: `
			SELECT id, message
			FROM raw_messages
		`,
		variant: cleaner.VariantSlack,
		source:  "slack",
	}

	if lastRawMessageID > 0 {
		plan.query += `
			WHERE id > $1
		`
		plan.args = append(plan.args, lastRawMessageID)
	}

	plan.query += `
			ORDER BY id
	`

	return plan
}

func buildDiscordRawMessageLearningPlan(lastRawMessageID uint64, allowed map[string]struct{}) rawMessageLearningPlan {
	plan := rawMessageLearningPlan{
		query: `
			SELECT id, message
			FROM raw_messages_discord
		`,
		variant: cleaner.VariantDiscord,
		source:  "discord",
	}

	where := make([]string, 0, 2)

	where = append(
		where,
		fmt.Sprintf("created_at >= (CURRENT_TIMESTAMP - INTERVAL '%d days')", utils.DiscordLearningLookbackDays),
	)

	where = append(where, `
		NOT EXISTS (
			SELECT 1
			FROM raw_message_user_blocklist
			WHERE protocol_driver = 'discord'
			  AND server_id = raw_messages_discord.server_id
			  AND user_id = raw_messages_discord.author_user_id
		)
	`)

	if lastRawMessageID > 0 {
		plan.args = append(plan.args, lastRawMessageID)
		where = append(where, fmt.Sprintf("id > $%d", len(plan.args)))
	}

	channelIDs := make([]string, 0, len(allowed))
	for channelID := range allowed {
		channelIDs = append(channelIDs, channelID)
	}

	sort.Strings(channelIDs)

	if len(channelIDs) > 0 {
		placeholders := make([]string, 0, len(channelIDs))

		for _, channelID := range channelIDs {
			plan.args = append(plan.args, channelID)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(plan.args)))
		}

		where = append(where, fmt.Sprintf("channel_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(where) > 0 {
		plan.query += `
			WHERE ` + strings.Join(where, " AND ") + `
		`
	}

	plan.query += `
			ORDER BY id
	`

	return plan
}

func buildDiscordAuthorRawMessageLearningPlan(serverID, authorUserID string, allowed map[string]struct{}) rawMessageLearningPlan {
	plan := rawMessageLearningPlan{
		query: `
			SELECT id, message
			FROM raw_messages_discord
		`,
		variant: cleaner.VariantDiscord,
		source:  "discord-author",
	}

	where := make([]string, 0, 5)

	where = append(where, "server_id = $1")
	plan.args = append(plan.args, serverID)

	where = append(where, "author_user_id = $2")
	plan.args = append(plan.args, authorUserID)

	where = append(
		where,
		fmt.Sprintf("created_at >= (CURRENT_TIMESTAMP - INTERVAL '%d days')", utils.DiscordLearningLookbackDays),
	)

	where = append(where, `
		NOT EXISTS (
			SELECT 1
			FROM raw_message_user_blocklist
			WHERE protocol_driver = 'discord'
			  AND server_id = raw_messages_discord.server_id
			  AND user_id = raw_messages_discord.author_user_id
		)
	`)

	channelIDs := make([]string, 0, len(allowed))
	for channelID := range allowed {
		channelIDs = append(channelIDs, channelID)
	}

	sort.Strings(channelIDs)

	if len(channelIDs) > 0 {
		placeholders := make([]string, 0, len(channelIDs))

		for _, channelID := range channelIDs {
			plan.args = append(plan.args, channelID)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(plan.args)))
		}

		where = append(where, fmt.Sprintf("channel_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	plan.query += `
			WHERE ` + strings.Join(where, " AND ") + `
			ORDER BY id
	`

	return plan
}

func keyToValues(m map[string]int) []string {
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
