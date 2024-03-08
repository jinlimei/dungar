package core2

import (
	"time"
)

// ScheduleHandler is our standardized func for handling the schedule
type ScheduleHandler func(svc *Service, currentOutgoingCount uint, lastSentMessage time.Time) ScheduledMessages

// HandleSchedule is responsible for handling some schedules.
func (s *Service) HandleSchedule(currentOutgoingCount uint) ScheduledMessages {
	var (
		messages = s.scheduler(s, currentOutgoingCount, s.lastSent)
		outgoing = make(ScheduledMessages, 0, len(messages))
	)

	// We need to populate channel IDs from their archetypes
	for _, message := range messages {
		if message.ChannelID != "" {
			outgoing = append(outgoing, message)
			continue
		}

		//switch message.ChannelArchetype {
		//case "any":
		//	// We'll pick a safe channel at random
		//	channels, err := s.infoDriver.GetChannelsByType(ChannelPublic)
		//
		//	if err != nil {
		//		log.Printf("ERROR: could not retrieve channels by type public: %v", err)
		//		continue
		//	}
		//
		//	if len(channels) == 0 {
		//		// No channels for the type, rip
		//		log.Println("WARNING: Could not find a channel for type public")
		//		continue
		//	}
		//
		//	channel := channels[random.Int(len(channels))]
		//	message.ChannelID = channel.ID
		//	message.ServerID = channel.ServerID
		//
		//	outgoing = append(outgoing, message)
		//
		//default:
		//	// We'll do a lookup specifically for a channel
		//	channels, err := s.infoDriver.GetChannelsByArchetype(message.ChannelArchetype)
		//
		//	if err != nil {
		//		log.Printf("ERROR: could not retrieve channels by archetype: %v", err)
		//		continue
		//	}
		//
		//	if len(channels) == 0 {
		//		// No channels for the archetype, rip
		//		log.Printf("WARNING: Could not find a channel for archetype '%s'", message.ChannelArchetype)
		//		continue
		//	}
		//
		//	channel := channels[random.Int(len(channels))]
		//	message.ChannelID = channel.ID
		//	outgoing = append(outgoing, message)
		//}

	}

	return outgoing
}
