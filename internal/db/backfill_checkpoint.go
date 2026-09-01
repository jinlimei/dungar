package db

import (
	"database/sql"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

type DiscordBackfillCheckpoint struct {
	ChannelID            string
	GuildID              string
	ChannelName          string
	Status               string
	BeforeMessageID      string
	Since                time.Time
	Fetched              int
	Stored               int
	SkippedDuplicate     int
	SkippedUnusable      int
	SkippedOld           int
	LastMessageTimestamp time.Time
	LastError            string
	StartedAt            time.Time
	UpdatedAt            time.Time
	CompletedAt          time.Time
}

func GetDiscordBackfillCheckpoint(channelID string) (*DiscordBackfillCheckpoint, error) {
	if utils.InTestEnv() || channelID == "" {
		return nil, nil
	}

	row := ConQueryRow(`
		SELECT
			channel_id,
			guild_id,
			COALESCE(channel_name, ''),
			status,
			COALESCE(before_message_id, ''),
			since_ts,
			fetched,
			stored,
			skipped_duplicate,
			skipped_unusable,
			skipped_old,
			COALESCE(last_message_timestamp, 'epoch'::timestamptz),
			COALESCE(last_error, ''),
			started_at,
			updated_at,
			COALESCE(completed_at, 'epoch'::timestamptz)
		FROM discord_backfill_checkpoints
		WHERE channel_id = $1
	`, channelID)

	if row == nil {
		return nil, nil
	}

	cp := &DiscordBackfillCheckpoint{}
	if err := row.Scan(
		&cp.ChannelID,
		&cp.GuildID,
		&cp.ChannelName,
		&cp.Status,
		&cp.BeforeMessageID,
		&cp.Since,
		&cp.Fetched,
		&cp.Stored,
		&cp.SkippedDuplicate,
		&cp.SkippedUnusable,
		&cp.SkippedOld,
		&cp.LastMessageTimestamp,
		&cp.LastError,
		&cp.StartedAt,
		&cp.UpdatedAt,
		&cp.CompletedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return cp, nil
}

func UpsertDiscordBackfillCheckpoint(cp *DiscordBackfillCheckpoint) error {
	if utils.InTestEnv() || cp == nil || cp.ChannelID == "" {
		return nil
	}

	now := time.Now().UTC()
	if cp.StartedAt.IsZero() {
		cp.StartedAt = now
	}

	cp.UpdatedAt = now

	_, err := ConExec(`
		INSERT INTO discord_backfill_checkpoints(
			channel_id,
			guild_id,
			channel_name,
			status,
			before_message_id,
			since_ts,
			fetched,
			stored,
			skipped_duplicate,
			skipped_unusable,
			skipped_old,
			last_message_timestamp,
			last_error,
			started_at,
			updated_at,
			completed_at
		)
		VALUES(
			$1, $2, NULLIF($3, ''), $4, NULLIF($5, ''), $6, $7, $8, $9, $10, $11,
			NULLIF(CAST($12 AS timestamptz), 'epoch'::timestamptz),
			NULLIF($13, ''), $14, $15,
			NULLIF(CAST($16 AS timestamptz), 'epoch'::timestamptz)
		)
		ON CONFLICT (channel_id) DO UPDATE
		SET guild_id = EXCLUDED.guild_id,
		    channel_name = EXCLUDED.channel_name,
		    status = EXCLUDED.status,
		    before_message_id = EXCLUDED.before_message_id,
		    since_ts = EXCLUDED.since_ts,
		    fetched = EXCLUDED.fetched,
		    stored = EXCLUDED.stored,
		    skipped_duplicate = EXCLUDED.skipped_duplicate,
		    skipped_unusable = EXCLUDED.skipped_unusable,
		    skipped_old = EXCLUDED.skipped_old,
		    last_message_timestamp = EXCLUDED.last_message_timestamp,
		    last_error = EXCLUDED.last_error,
		    started_at = EXCLUDED.started_at,
		    updated_at = EXCLUDED.updated_at,
		    completed_at = EXCLUDED.completed_at
	`,
		cp.ChannelID,
		cp.GuildID,
		cp.ChannelName,
		cp.Status,
		cp.BeforeMessageID,
		cp.Since.UTC(),
		cp.Fetched,
		cp.Stored,
		cp.SkippedDuplicate,
		cp.SkippedUnusable,
		cp.SkippedOld,
		zeroTimeIfEmpty(cp.LastMessageTimestamp),
		cp.LastError,
		cp.StartedAt.UTC(),
		cp.UpdatedAt.UTC(),
		zeroTimeIfEmpty(cp.CompletedAt),
	)

	return err
}

func zeroTimeIfEmpty(inp time.Time) time.Time {
	if inp.IsZero() {
		return time.Unix(0, 0).UTC()
	}

	return inp.UTC()
}
