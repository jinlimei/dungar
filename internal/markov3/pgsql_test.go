package markov3

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

func zTestMarkovLearnFromRawMessages(t *testing.T) {
	assert.Equal(t, 1, 1)

	db.TestDatabaseConnect()
	random.UseTestSeed()

	n := time.Now()
	m := MakeMarkov("")
	m.LearnFromRawMessagesM1()

	for k := 0; k < 50; k++ {
		log.Printf("%02d: %s\n", k, m.Generate("genius"))
	}

	log.Printf("Finished learning, having %d words and %d fragments. Took %s\n",
		len(m.RevWords), len(m.Fragments), time.Now().Sub(n).String())
}

func TestBuildSlackRawMessageLearningPlan(t *testing.T) {
	plan := buildSlackRawMessageLearningPlan(42)

	assert.Equal(t, cleaner.VariantSlack, plan.variant)
	assert.Equal(t, "slack", plan.source)
	assert.Equal(t, []any{uint64(42)}, plan.args)
	assert.Contains(t, plan.query, "FROM raw_messages")
	assert.Contains(t, plan.query, "WHERE id > $1")
	assert.Contains(t, plan.query, "ORDER BY id")
}

func TestBuildDiscordRawMessageLearningPlan(t *testing.T) {
	lookbackClause := fmt.Sprintf(
		"created_at >= (CURRENT_TIMESTAMP - INTERVAL '%d days')",
		utils.DiscordLearningLookbackDays,
	)

	plan := buildDiscordRawMessageLearningPlan(42, map[string]struct{}{
		"471817271749115905": {},
		"362144221148479489": {},
	})

	assert.Equal(t, cleaner.VariantDiscord, plan.variant)
	assert.Equal(t, "discord", plan.source)
	assert.Equal(t, []any{
		uint64(42),
		"362144221148479489",
		"471817271749115905",
	}, plan.args)
	assert.Contains(t, plan.query, "FROM raw_messages_discord")
	assert.Contains(t, plan.query, lookbackClause)
	assert.Contains(t, plan.query, "FROM raw_message_user_blocklist")
	assert.Contains(t, plan.query, "protocol_driver = 'discord'")
	assert.Contains(t, plan.query, "id > $1")
	assert.Contains(t, plan.query, "channel_id IN ($2, $3)")
	assert.Contains(t, plan.query, "ORDER BY id")
}

func TestBuildDiscordRawMessageLearningPlanUnrestricted(t *testing.T) {
	lookbackClause := fmt.Sprintf(
		"created_at >= (CURRENT_TIMESTAMP - INTERVAL '%d days')",
		utils.DiscordLearningLookbackDays,
	)

	plan := buildDiscordRawMessageLearningPlan(0, map[string]struct{}{})

	assert.Equal(t, cleaner.VariantDiscord, plan.variant)
	assert.Equal(t, "discord", plan.source)
	assert.Empty(t, plan.args)
	assert.Contains(t, plan.query, "FROM raw_messages_discord")
	assert.Contains(t, plan.query, lookbackClause)
	assert.Contains(t, plan.query, "FROM raw_message_user_blocklist")
	assert.NotContains(t, plan.query, "channel_id IN")
	assert.Contains(t, plan.query, "ORDER BY id")
}
