package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

type ManicMinuteEvent struct {
	ID                int64
	ServerID          string
	ChannelID         string
	TriggerWord       string
	TriggerMessageID  string
	TriggeredByUserID string
	Status            string
	StopReason        string
	MessageCount      int
	StartedAt         time.Time
	CooldownUntil     time.Time
	EndedAt           time.Time
}

type ManicMinuteWordCount struct {
	Word  string
	Count int
}

type ManicMinuteStats struct {
	TotalEvents int
	Events24h   int
	Events7d    int
	MostRecent  *ManicMinuteEvent
	TopWords    []ManicMinuteWordCount
}

type ManicMinuteRuntimeState struct {
	ProtocolDriver  string
	TriggerWord     string
	StartChance     float64
	Active          bool
	ActiveServerID  string
	ActiveChannelID string
	HasCooldown     bool
	CooldownUntil   time.Time
	UpdatedReason   string
	UpdatedAt       time.Time
}

var (
	manicMinuteCooldownColumnCheck sync.Once
	manicMinuteHasCooldownColumn   bool
)

func StartManicMinuteEvent(serverID, channelID, triggerWord, triggerMessageID, triggeredByUserID string, startedAt, cooldownUntil time.Time) (int64, error) {
	if utils.InTestEnv() {
		return 0, nil
	}

	var row *sql.Row
	if manicMinuteEventHasCooldownColumn() {
		row = ConQueryRow(`
			INSERT INTO manic_minute_events(
				server_id,
				channel_id,
				trigger_word,
				trigger_message_id,
				triggered_by_user_id,
				status,
				started_at,
				cooldown_until
			)
			VALUES($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), 'active', $6, $7)
			RETURNING id
		`, serverID, channelID, triggerWord, triggerMessageID, triggeredByUserID, startedAt.UTC(), cooldownUntil.UTC())
	} else {
		row = ConQueryRow(`
			INSERT INTO manic_minute_events(
				server_id,
				channel_id,
				trigger_word,
				trigger_message_id,
				triggered_by_user_id,
				status,
				started_at
			)
			VALUES($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), 'active', $6)
			RETURNING id
		`, serverID, channelID, triggerWord, triggerMessageID, triggeredByUserID, startedAt.UTC())
	}

	if row == nil {
		return 0, nil
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func CompleteManicMinuteEvent(eventID int64, status, stopReason string, endedAt time.Time, messageCount int) error {
	if utils.InTestEnv() || eventID <= 0 {
		return nil
	}

	_, err := ConExec(`
		UPDATE manic_minute_events
		SET status = $2,
		    stop_reason = NULLIF($3, ''),
		    message_count = $4,
		    ended_at = $5
		WHERE id = $1
	`, eventID, status, stopReason, messageCount, endedAt.UTC())
	return err
}

func IsManicMinuteChannelOnCooldown(serverID, channelID string) bool {
	if utils.InTestEnv() || serverID == "" || channelID == "" {
		return false
	}

	var row *sql.Row
	if manicMinuteEventHasCooldownColumn() {
		row = ConQueryRow(`
			SELECT 1
			FROM manic_minute_events
			WHERE server_id = $1
			  AND channel_id = $2
			  AND cooldown_until > CURRENT_TIMESTAMP
			LIMIT 1
		`, serverID, channelID)
	} else {
		row = ConQueryRow(`
			SELECT 1
			FROM manic_minute_events
			WHERE server_id = $1
			  AND channel_id = $2
			  AND started_at >= CURRENT_TIMESTAMP - INTERVAL '120 minutes'
			LIMIT 1
		`, serverID, channelID)
	}

	if row == nil {
		return false
	}

	var exists int
	if err := row.Scan(&exists); err != nil {
		return false
	}

	return exists == 1
}

func IsManicMinuteWordOnCooldown(serverID, word string) bool {
	if utils.InTestEnv() || serverID == "" || word == "" {
		return false
	}

	var row *sql.Row
	if manicMinuteEventHasCooldownColumn() {
		row = ConQueryRow(`
			SELECT 1
			FROM manic_minute_events
			WHERE server_id = $1
			  AND trigger_word = $2
			  AND cooldown_until > CURRENT_TIMESTAMP
			LIMIT 1
		`, serverID, word)
	} else {
		row = ConQueryRow(`
			SELECT 1
			FROM manic_minute_events
			WHERE server_id = $1
			  AND trigger_word = $2
			  AND started_at >= CURRENT_TIMESTAMP - INTERVAL '120 minutes'
			LIMIT 1
		`, serverID, word)
	}

	if row == nil {
		return false
	}

	var exists int
	if err := row.Scan(&exists); err != nil {
		return false
	}

	return exists == 1
}

func GetMostRecentManicMinuteEvent(serverID string) (*ManicMinuteEvent, error) {
	if utils.InTestEnv() || serverID == "" {
		return nil, nil
	}

	var row *sql.Row
	if manicMinuteEventHasCooldownColumn() {
		row = ConQueryRow(`
			SELECT
				id,
				server_id,
				channel_id,
				trigger_word,
				COALESCE(trigger_message_id, ''),
				COALESCE(triggered_by_user_id, ''),
				status,
				COALESCE(stop_reason, ''),
				message_count,
				started_at,
				cooldown_until,
				COALESCE(ended_at, started_at)
			FROM manic_minute_events
			WHERE server_id = $1
			ORDER BY started_at DESC
			LIMIT 1
		`, serverID)
	} else {
		row = ConQueryRow(`
			SELECT
				id,
				server_id,
				channel_id,
				trigger_word,
				COALESCE(trigger_message_id, ''),
				COALESCE(triggered_by_user_id, ''),
				status,
				COALESCE(stop_reason, ''),
				message_count,
				started_at,
				started_at + INTERVAL '120 minutes',
				COALESCE(ended_at, started_at)
			FROM manic_minute_events
			WHERE server_id = $1
			ORDER BY started_at DESC
			LIMIT 1
		`, serverID)
	}

	if row == nil {
		return nil, nil
	}

	ev, err := scanManicMinuteEvent(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return ev, nil
}

func GetManicMinuteStats(serverID string) (ManicMinuteStats, error) {
	var stats ManicMinuteStats

	if utils.InTestEnv() || serverID == "" {
		return stats, nil
	}

	row := ConQueryRow(`
		SELECT
			COUNT(*) AS total_events,
			COUNT(*) FILTER (
				WHERE started_at >= CURRENT_TIMESTAMP - INTERVAL '24 hours'
			) AS events_24h,
			COUNT(*) FILTER (
				WHERE started_at >= CURRENT_TIMESTAMP - INTERVAL '7 days'
			) AS events_7d
		FROM manic_minute_events
		WHERE server_id = $1
	`, serverID)

	if row != nil {
		if err := row.Scan(&stats.TotalEvents, &stats.Events24h, &stats.Events7d); err != nil && err != sql.ErrNoRows {
			return stats, err
		}
	}

	mostRecent, err := GetMostRecentManicMinuteEvent(serverID)
	if err != nil {
		return stats, err
	}

	stats.MostRecent = mostRecent
	stats.TopWords = make([]ManicMinuteWordCount, 0, 5)

	rows, err := ConQuery(`
		SELECT trigger_word, COUNT(*) AS used_count
		FROM manic_minute_events
		WHERE server_id = $1
		GROUP BY trigger_word
		ORDER BY used_count DESC, trigger_word
		LIMIT 5
	`, serverID)
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry ManicMinuteWordCount
		if err := rows.Scan(&entry.Word, &entry.Count); err != nil {
			return stats, err
		}

		stats.TopWords = append(stats.TopWords, entry)
	}

	return stats, nil
}

func GetManicMinuteRuntimeState(protocolDriver string) (*ManicMinuteRuntimeState, error) {
	if utils.InTestEnv() || protocolDriver == "" {
		return nil, nil
	}

	row := ConQueryRow(`
		SELECT
			protocol_driver,
			trigger_word,
			start_chance,
			active,
			COALESCE(active_server_id, ''),
			COALESCE(active_channel_id, ''),
			cooldown_until,
			COALESCE(updated_reason, ''),
			updated_at
		FROM manic_minute_runtime_state
		WHERE protocol_driver = $1
	`, protocolDriver)

	if row == nil {
		return nil, nil
	}

	state := &ManicMinuteRuntimeState{}
	var cooldownUntil sql.NullTime
	if err := row.Scan(
		&state.ProtocolDriver,
		&state.TriggerWord,
		&state.StartChance,
		&state.Active,
		&state.ActiveServerID,
		&state.ActiveChannelID,
		&cooldownUntil,
		&state.UpdatedReason,
		&state.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	state.HasCooldown = cooldownUntil.Valid
	if cooldownUntil.Valid {
		state.CooldownUntil = cooldownUntil.Time
	}

	return state, nil
}

func UpsertManicMinuteRuntimeState(state *ManicMinuteRuntimeState) error {
	if utils.InTestEnv() || state == nil || state.ProtocolDriver == "" {
		return nil
	}

	if state.UpdatedAt.IsZero() {
		state.UpdatedAt = time.Now().UTC()
	}

	var cooldownUntil interface{}
	if state.HasCooldown && !state.CooldownUntil.IsZero() {
		cooldownUntil = state.CooldownUntil.UTC()
	}

	_, err := ConExec(`
		INSERT INTO manic_minute_runtime_state(
			protocol_driver,
			trigger_word,
			start_chance,
			active,
			active_server_id,
			active_channel_id,
			cooldown_until,
			updated_reason,
			updated_at
		)
		VALUES($1, $2, $3, $4, NULLIF($5, ''), NULLIF($6, ''), $7, NULLIF($8, ''), $9)
		ON CONFLICT (protocol_driver) DO UPDATE
		SET trigger_word = EXCLUDED.trigger_word,
		    start_chance = EXCLUDED.start_chance,
		    active = EXCLUDED.active,
		    active_server_id = EXCLUDED.active_server_id,
		    active_channel_id = EXCLUDED.active_channel_id,
		    cooldown_until = EXCLUDED.cooldown_until,
		    updated_reason = EXCLUDED.updated_reason,
		    updated_at = EXCLUDED.updated_at
	`,
		state.ProtocolDriver,
		state.TriggerWord,
		state.StartChance,
		state.Active,
		state.ActiveServerID,
		state.ActiveChannelID,
		cooldownUntil,
		state.UpdatedReason,
		state.UpdatedAt.UTC(),
	)

	return err
}

func scanManicMinuteEvent(scanner interface {
	Scan(dest ...interface{}) error
}) (*ManicMinuteEvent, error) {
	ev := &ManicMinuteEvent{}

	if err := scanner.Scan(
		&ev.ID,
		&ev.ServerID,
		&ev.ChannelID,
		&ev.TriggerWord,
		&ev.TriggerMessageID,
		&ev.TriggeredByUserID,
		&ev.Status,
		&ev.StopReason,
		&ev.MessageCount,
		&ev.StartedAt,
		&ev.CooldownUntil,
		&ev.EndedAt,
	); err != nil {
		return nil, err
	}

	return ev, nil
}

func (ev ManicMinuteEvent) String() string {
	return fmt.Sprintf(
		"id=%d server=%s channel=%s trigger=%s status=%s messages=%d cooldown_until=%s",
		ev.ID,
		ev.ServerID,
		ev.ChannelID,
		ev.TriggerWord,
		ev.Status,
		ev.MessageCount,
		ev.CooldownUntil.Format(time.RFC3339),
	)
}

func manicMinuteEventHasCooldownColumn() bool {
	if utils.InTestEnv() {
		return false
	}

	manicMinuteCooldownColumnCheck.Do(func() {
		row := ConQueryRow(`
			SELECT EXISTS(
				SELECT 1
				FROM information_schema.columns
				WHERE table_schema = CURRENT_SCHEMA()
				  AND table_name = 'manic_minute_events'
				  AND column_name = 'cooldown_until'
			)
		`)
		if row == nil {
			return
		}

		_ = row.Scan(&manicMinuteHasCooldownColumn)
	})

	return manicMinuteHasCooldownColumn
}
