package accord

import "github.com/bwmarrin/discordgo"

func isMessageFromSelf(s *discordgo.Session, ev *discordgo.Message) bool {
	if ev.Author == nil {
		return false
	}

	if s.State == nil || s.State.User == nil {
		return false
	}

	return ev.Author.ID == s.State.User.ID
}

func isConsumableMessage(ev *discordgo.MessageCreate) bool {
	// by default we'll ignore webhook-created messages.
	if ev.WebhookID != "" {
		return false
	}

	// We don't want to deal with messages which have _zero_ content
	if ev.Content == "" {
		return false
	}

	if ev.Type == discordgo.MessageTypeDefault || ev.Type == discordgo.MessageTypeReply || ev.Type == discordgo.MessageTypeChatInputCommand {
		return true
	}

	return false
}
