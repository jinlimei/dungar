package dcli

import (
	"log"
	"sort"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/accord"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// BackfillDiscordHistory imports the last DiscordLearningLookbackDays days of
// message history from the configured Discord output/learning allowlist.
func BackfillDiscordHistory() {
	utils.LoadSettingsAndSecrets()
	db.ConnectToDatabase()

	allowed := utils.DiscordAllowedLearningChannelIDs()
	if len(allowed) == 0 {
		log.Fatal("Discord backfill requires at least one allowed_output_channel_ids entry")
	}

	channelIDs := make([]string, 0, len(allowed))
	for channelID := range allowed {
		channelIDs = append(channelIDs, channelID)
	}
	sort.Strings(channelIDs)

	con := accord.NewRealDiscordConnection()
	driver, err := accord.New(con)
	if err != nil {
		log.Fatalf("Failed to create discord driver for backfill: %v", err)
	}

	if err := con.Start(utils.DiscordAccessToken()); err != nil {
		log.Fatalf("Failed to start discord session for backfill: %v", err)
	}

	since := time.Now().UTC().AddDate(0, 0, -utils.DiscordLearningLookbackDays)
	log.Printf(
		"Discord history backfill: channelIDs=%v since=%s",
		channelIDs,
		since.Format(time.RFC3339),
	)

	stats, err := driver.BackfillHistory(channelIDs, since)
	if err != nil {
		log.Fatalf("Discord history backfill failed: %v", err)
	}

	log.Printf(
		"Discord history backfill complete: channels=%d resumed=%d failed=%d fetched=%d stored=%d duplicates=%d skipped_unusable=%d skipped_old=%d",
		stats.Channels,
		stats.ResumedChannels,
		stats.FailedChannels,
		stats.Fetched,
		stats.Stored,
		stats.SkippedDuplicate,
		stats.SkippedUnusable,
		stats.SkippedOld,
	)
}
