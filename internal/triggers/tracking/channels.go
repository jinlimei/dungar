package tracking

import (
	"encoding/json"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"log"
	"time"
)

func (s *Service) getChannelByID(channelID, serverID string) (*core2.Channel, error) {
	qry := `
SELECT name, name_normalized, channel_type, topic, can_post, archetype, previous_names, created_at, updated_at
FROM channel_tracking
WHERE channel_id = $1 
  AND server_id = $2
  AND protocol_driver = $3
`

	row := db.ConQueryRow(qry, channelID, serverID, s.driverName)

	var (
		name           string
		nameNormalized string
		chanType       int
		topic          string
		canPost        bool
		archetype      string
		prevNamesRaw   []byte
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := row.Scan(
		&name,
		&nameNormalized,
		&chanType,
		&topic,
		&canPost,
		&archetype,
		&prevNamesRaw,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	var prevNames []string

	if err = json.Unmarshal(prevNamesRaw, &prevNames); err != nil {
		return nil, err
	}

	return &core2.Channel{
		ID:             channelID,
		Name:           name,
		NameNormalized: nameNormalized,
		PreviousNames:  prevNames,
		Topic:          topic,
		Type:           core2.ChannelType(chanType),
		Archetype:      archetype,
	}, nil
}

func (s *Service) createChannel(channel *core2.Channel) error {
	qry := `
INSERT INTO channel_tracking (channel_id, server_id, protocol_driver, name, name_normalized, 
														  topic, channel_type, can_post, archetype, previous_names, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
`

	bytes, _ := json.Marshal(channel.PreviousNames)

	_, err := db.ConExec(
		qry,
		channel.ID,             // 1
		channel.ServerID,       // 2
		s.driverName,           // 3
		channel.Name,           // 4
		channel.NameNormalized, // 5
		channel.Topic,          // 6
		channel.Type,           // 7
		channel.CanPost,        // 8
		channel.Archetype,      // 9
		bytes,                  // 10
	)

	return err
}

func (s *Service) updateChannel(channel *core2.Channel) error {
	qry := `
UPDATE channel_tracking
SET name = $4,
    name_normalized = $5,
    can_post = $6,
    archetype = $7,
    topic = $8,
    channel_type = $9,
    previous_names = $10,
    updated_at = CURRENT_TIMESTAMP
WHERE channel_id = $1
  AND server_id = $2
  AND protocol_driver = $3
`

	bytes, _ := json.Marshal(channel.PreviousNames)

	_, err := db.ConExec(
		qry,
		channel.ID,             // 1
		channel.ServerID,       // 2
		s.driverName,           // 3
		channel.Name,           // 4
		channel.NameNormalized, // 5
		channel.CanPost,        // 6
		channel.Archetype,      // 7
		channel.Topic,          // 8
		channel.Type,           // 9
		bytes,                  // 10
	)

	return err
}

// StoreChannel stores a channel!
func (s *Service) StoreChannel(channel *core2.Channel) error {
	_, err := s.getChannelByID(channel.ID, channel.ServerID)

	if err == nil {
		return s.updateChannel(channel)
	}

	return s.createChannel(channel)
}

// SetChannelArchetype accurately
func (s *Service) SetChannelArchetype(channelID, serverID, archetype string) error {
	channel, err := s.getChannelByID(channelID, serverID)

	if err != nil {
		return err
	}

	channel.Archetype = archetype
	return s.updateChannel(channel)
}

// GetChannelByID returns a specific channel by ID and serverID
func (s *Service) GetChannelByID(channelID, serverID string) (*core2.Channel, error) {
	return s.getChannelByID(channelID, serverID)
}

// GetChannelsByArchetype returns a list of channels by archetype for a specific serverID
func (s *Service) GetChannelsByArchetype(archetype, serverID string) ([]*core2.Channel, error) {
	// Lazily just fetch the channelID's to query afterward
	qry := `
SELECT channel_id
FROM channel_tracking
WHERE archetype = $1 
  AND server_id = $2
  AND protocol_driver = $3
`
	rows, err := db.ConQuery(qry, archetype, serverID, s.driverName)

	if err != nil {
		return nil, err
	}

	channels := make([]*core2.Channel, 0)

	var (
		channelID string
		channel   *core2.Channel
	)

	for rows.Next() {
		err = rows.Scan(&channelID)
		if err != nil {
			log.Printf("ERROR: failed to scan row for archetype '%s': %v", archetype, err)
			return nil, err
		}

		channel, err = s.getChannelByID(channelID, serverID)
		if err != nil {
			log.Printf("ERROR: failed to get channel by id '%s': %v", channelID, err)
			return nil, err
		}

		channels = append(channels, channel)
	}

	return channels, nil
}

// GetChannelsByType returns a list of channels by type from a specific serverID
func (s *Service) GetChannelsByType(serverID string, chanType core2.ChannelType) ([]*core2.Channel, error) {
	// Lazily just fetch the channelID's to query afterward
	qry := `
SELECT channel_id
FROM channel_tracking
WHERE channel_type = $1 
  AND server_id = $2
  AND protocol_driver = $3
`
	rows, err := db.ConQuery(qry, int(chanType), serverID, s.driverName)

	if err != nil {
		return nil, err
	}

	channels := make([]*core2.Channel, 0)

	var (
		channelID string
		channel   *core2.Channel
	)

	for rows.Next() {
		err = rows.Scan(&channelID)

		if err != nil {
			log.Printf("ERROR: failed to scan row for channel type '%s': %v", chanType.String(), err)
			return nil, err
		}

		channel, err = s.getChannelByID(channelID, serverID)

		if err != nil {
			log.Printf("ERROR: failed to get channel by id '%s': %v", channelID, err)
			return nil, err
		}

		channels = append(channels, channel)
	}

	return channels, nil
}
