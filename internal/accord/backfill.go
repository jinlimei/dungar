package accord

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// HistoryBackfillStats tracks the outcome of a one-shot Discord history import.
type HistoryBackfillStats struct {
	Channels         int
	ResumedChannels  int
	FailedChannels   int
	Fetched          int
	Stored           int
	SkippedDuplicate int
	SkippedUnusable  int
	SkippedOld       int
}

// BackfillHistory imports Discord history for the provided channels from the
// given timestamp onward. The imported message format matches the live
// conversion path used by handleMessageCreateEvent.
func (d *Driver) BackfillHistory(channelIDs []string, since time.Time) (HistoryBackfillStats, error) {
	var stats HistoryBackfillStats

	sess := d.Con.GetSession()
	if sess == nil {
		return stats, ErrConnectNotCalled
	}

	me, err := sess.User("@me")
	if err != nil {
		return stats, fmt.Errorf("fetch current bot user: %w", err)
	}

	d.SetBotUser(&core2.BotUser{
		ID:    me.ID,
		IsBot: true,
	})

	sort.Strings(channelIDs)

	for _, channelID := range channelIDs {
		if channelID == "" {
			continue
		}

		guildID, err := d.seedGuildCache(channelID)
		if err != nil {
			stats.FailedChannels++
			log.Printf("Discord history backfill failed to seed channelID='%s': %v", channelID, err)
			_ = db.UpsertDiscordBackfillCheckpoint(&db.DiscordBackfillCheckpoint{
				ChannelID:   channelID,
				GuildID:     "",
				Status:      "failed",
				Since:       since.UTC(),
				LastError:   err.Error(),
				StartedAt:   time.Now().UTC(),
				CompletedAt: time.Time{},
			})
			continue
		}

		channelName := d.GetChannelName(channelID, guildID)
		log.Printf(
			"Discord history backfill channel start: channelID='%s' channel='%s' guildID='%s' since=%s",
			channelID,
			channelName,
			guildID,
			since.Format(time.RFC3339),
		)

		perChannel, resumed, err := d.backfillChannelHistory(channelID, guildID, channelName, since, me.ID)
		if err != nil {
			stats.FailedChannels++
			log.Printf(
				"Discord history backfill channel failed: channelID='%s' channel='%s' guildID='%s' error=%v",
				channelID,
				channelName,
				guildID,
				err,
			)
			continue
		}

		if resumed {
			stats.ResumedChannels++
		}

		log.Printf(
			"Discord history backfill channel complete: channelID='%s' channel='%s' fetched=%d stored=%d duplicates=%d skipped_unusable=%d skipped_old=%d",
			channelID,
			channelName,
			perChannel.Fetched,
			perChannel.Stored,
			perChannel.SkippedDuplicate,
			perChannel.SkippedUnusable,
			perChannel.SkippedOld,
		)

		stats.Channels++
		stats.Fetched += perChannel.Fetched
		stats.Stored += perChannel.Stored
		stats.SkippedDuplicate += perChannel.SkippedDuplicate
		stats.SkippedUnusable += perChannel.SkippedUnusable
		stats.SkippedOld += perChannel.SkippedOld
	}

	return stats, nil
}

func (d *Driver) seedGuildCache(channelID string) (string, error) {
	sess := d.Con.GetSession()

	channel, err := sess.Channel(channelID)
	if err != nil {
		return "", fmt.Errorf("fetch channel %s: %w", channelID, err)
	}

	guild := d.getOrMakeGuild(channel.GuildID)
	guild.channelCache[channel.ID] = channel

	channels, err := sess.GuildChannels(channel.GuildID)
	if err != nil {
		return "", fmt.Errorf("fetch guild channels for %s: %w", channel.GuildID, err)
	}

	for _, chn := range channels {
		guild.channelCache[chn.ID] = chn
	}

	guildResp, err := sess.Guild(channel.GuildID)
	if err == nil {
		for _, role := range guildResp.Roles {
			guild.roleCache[role.ID] = role
		}
	}

	return channel.GuildID, nil
}

func (d *Driver) backfillChannelHistory(channelID, guildID, channelName string, since time.Time, selfUserID string) (HistoryBackfillStats, bool, error) {
	var stats HistoryBackfillStats

	sess := d.Con.GetSession()
	beforeID := ""
	page := 0
	since = since.UTC()
	resumed := false

	checkpoint, err := db.GetDiscordBackfillCheckpoint(channelID)
	if err != nil {
		return stats, resumed, fmt.Errorf("load checkpoint for %s: %w", channelID, err)
	}

	if checkpoint != nil && (checkpoint.Status == "in_progress" || checkpoint.Status == "failed") {
		resumed = true
		beforeID = checkpoint.BeforeMessageID
		stats.Fetched = checkpoint.Fetched
		stats.Stored = checkpoint.Stored
		stats.SkippedDuplicate = checkpoint.SkippedDuplicate
		stats.SkippedUnusable = checkpoint.SkippedUnusable
		stats.SkippedOld = checkpoint.SkippedOld

		checkpoint.GuildID = guildID
		checkpoint.ChannelName = channelName
		checkpoint.Status = "in_progress"
		checkpoint.LastError = ""
		checkpoint.Since = since
		checkpoint.CompletedAt = time.Time{}

		log.Printf(
			"Discord history backfill resume: channelID='%s' channel='%s' beforeID='%s' fetched=%d stored=%d",
			channelID,
			channelName,
			beforeID,
			stats.Fetched,
			stats.Stored,
		)
	} else {
		checkpoint = &db.DiscordBackfillCheckpoint{
			ChannelID:   channelID,
			GuildID:     guildID,
			ChannelName: channelName,
			Status:      "in_progress",
			Since:       since,
			StartedAt:   time.Now().UTC(),
		}
	}

	if err := db.UpsertDiscordBackfillCheckpoint(checkpoint); err != nil {
		return stats, resumed, fmt.Errorf("save checkpoint for %s: %w", channelID, err)
	}

	for {
		page++

		msgs, err := sess.ChannelMessages(channelID, 100, beforeID, "", "")
		if err != nil {
			checkpoint.Status = "failed"
			checkpoint.LastError = err.Error()
			if saveErr := db.UpsertDiscordBackfillCheckpoint(checkpoint); saveErr != nil {
				log.Printf("Discord history backfill checkpoint save failed channelID='%s': %v", channelID, saveErr)
			}

			return stats, resumed, fmt.Errorf("fetch channel messages for %s: %w", channelID, err)
		}

		if len(msgs) == 0 {
			checkpoint.Status = "completed"
			checkpoint.BeforeMessageID = ""
			checkpoint.LastError = ""
			checkpoint.CompletedAt = time.Now().UTC()
			checkpoint.Fetched = stats.Fetched
			checkpoint.Stored = stats.Stored
			checkpoint.SkippedDuplicate = stats.SkippedDuplicate
			checkpoint.SkippedUnusable = stats.SkippedUnusable
			checkpoint.SkippedOld = stats.SkippedOld

			if err := db.UpsertDiscordBackfillCheckpoint(checkpoint); err != nil {
				return stats, resumed, fmt.Errorf("finalize checkpoint for %s: %w", channelID, err)
			}

			log.Printf(
				"Discord history backfill page: channelID='%s' page=%d batch=0 fetched=%d stored=%d duplicates=%d skipped_unusable=%d skipped_old=%d done=true reason='empty_batch'",
				channelID,
				page,
				stats.Fetched,
				stats.Stored,
				stats.SkippedDuplicate,
				stats.SkippedUnusable,
				stats.SkippedOld,
			)
			return stats, resumed, nil
		}

		reachedOlderThanWindow := false
		newestInBatch := msgs[0].Timestamp.UTC().Format(time.RFC3339)
		oldestInBatch := msgs[len(msgs)-1].Timestamp.UTC().Format(time.RFC3339)

		for _, msg := range msgs {
			if msg == nil {
				continue
			}

			if msg.GuildID == "" {
				msg.GuildID = guildID
			}

			if msg.Timestamp.UTC().Before(since) {
				stats.SkippedOld++
				reachedOlderThanWindow = true
				continue
			}

			stats.Fetched++

			if !isConsumableHistoricalDiscordMessage(selfUserID, msg) {
				stats.SkippedUnusable++
				continue
			}

			d.seedMessageUsers(msg)

			converted, err := d.convertMessageCreate(msg)
			if err != nil {
				checkpoint.Status = "failed"
				checkpoint.LastError = err.Error()
				_ = db.UpsertDiscordBackfillCheckpoint(checkpoint)
				return stats, resumed, fmt.Errorf("convert message %s: %w", msg.ID, err)
			}

			if db.IsRawMessageUserBlocked("discord", converted.ServerID, converted.UserID) {
				stats.SkippedUnusable++
				continue
			}

			inserted := db.RecordRawDiscordMessageAt(
				converted.ID,
				converted.ServerID,
				converted.ChannelID,
				converted.UserID,
				converted.Contents,
				msg.Timestamp.UTC(),
			)

			if inserted {
				stats.Stored++
			} else {
				stats.SkippedDuplicate++
			}
		}

		beforeID = msgs[len(msgs)-1].ID
		checkpoint.BeforeMessageID = beforeID
		checkpoint.LastError = ""
		checkpoint.Fetched = stats.Fetched
		checkpoint.Stored = stats.Stored
		checkpoint.SkippedDuplicate = stats.SkippedDuplicate
		checkpoint.SkippedUnusable = stats.SkippedUnusable
		checkpoint.SkippedOld = stats.SkippedOld
		checkpoint.LastMessageTimestamp = msgs[len(msgs)-1].Timestamp.UTC()
		checkpoint.Status = "in_progress"
		checkpoint.CompletedAt = time.Time{}

		done := reachedOlderThanWindow || len(msgs) < 100
		reason := "continue"
		if reachedOlderThanWindow {
			reason = "hit_lookback_window"
		} else if len(msgs) < 100 {
			reason = "short_batch"
		}

		log.Printf(
			"Discord history backfill page: channelID='%s' page=%d batch=%d newest=%s oldest=%s fetched=%d stored=%d duplicates=%d skipped_unusable=%d skipped_old=%d done=%t reason='%s'",
			channelID,
			page,
			len(msgs),
			newestInBatch,
			oldestInBatch,
			stats.Fetched,
			stats.Stored,
			stats.SkippedDuplicate,
			stats.SkippedUnusable,
			stats.SkippedOld,
			done,
			reason,
		)

		if done {
			checkpoint.Status = "completed"
			checkpoint.BeforeMessageID = ""
			checkpoint.CompletedAt = time.Now().UTC()
		}

		if err := db.UpsertDiscordBackfillCheckpoint(checkpoint); err != nil {
			return stats, resumed, fmt.Errorf("update checkpoint for %s: %w", channelID, err)
		}

		if reachedOlderThanWindow || len(msgs) < 100 {
			return stats, resumed, nil
		}
	}
}

func (d *Driver) seedMessageUsers(msg *discordgo.Message) {
	if msg == nil || msg.Author == nil {
		return
	}

	guild := d.getOrMakeGuild(msg.GuildID)

	authorMember := &discordgo.Member{
		GuildID: msg.GuildID,
		User:    msg.Author,
	}

	if msg.Member != nil {
		authorMember.Nick = msg.Member.Nick
	}

	guild.memberCache[msg.Author.ID] = authorMember

	for _, user := range msg.Mentions {
		if user == nil {
			continue
		}

		if _, ok := guild.memberCache[user.ID]; ok {
			continue
		}

		guild.memberCache[user.ID] = &discordgo.Member{
			GuildID: msg.GuildID,
			User:    user,
		}
	}
}

func isConsumableHistoricalDiscordMessage(selfUserID string, msg *discordgo.Message) bool {
	if msg == nil || msg.Author == nil {
		return false
	}

	if selfUserID != "" && msg.Author.ID == selfUserID {
		return false
	}

	if msg.WebhookID != "" {
		return false
	}

	if strings.TrimSpace(msg.Content) == "" {
		return false
	}

	switch msg.Type {
	case discordgo.MessageTypeDefault, discordgo.MessageTypeReply, discordgo.MessageTypeChatInputCommand:
		return true
	default:
		return false
	}
}
