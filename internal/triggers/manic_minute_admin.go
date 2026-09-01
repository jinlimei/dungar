package triggers

import (
	"fmt"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func manicMinuteAdminHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	txt := strings.TrimSpace(msg.Contents)
	if !strings.HasPrefix(strings.ToLower(txt), "!manic") {
		return core2.EmptyRsp()
	}

	if !canUseDiscordAdminCommand(svc, msg) {
		return core2.EmptyRsp()
	}

	pieces := strings.Fields(txt)
	if len(pieces) < 2 {
		return core2.MakeSingleRsp("Usage: !manic <status|stats|test|stop|rotate>")
	}

	switch strings.ToLower(pieces[1]) {
	case "status":
		return core2.MakeSingleRsp(renderManicMinuteStatus(svc, msg))
	case "stats":
		return core2.MakeSingleRsp(renderManicMinuteStats(msg.ServerID))
	case "rotate":
		if manicMinuteState.isActive() {
			return core2.MakeSingleRsp("Cannot rotate the trigger word while MANIC MINUTE is active.")
		}

		previous, next := manicMinuteState.rotate("admin-rotate")
		return core2.MakeSingleRsp(fmt.Sprintf("Rotated manic trigger from '%s' to '%s'.", previous, next))
	case "stop":
		result := manicMinuteState.stop("admin-stop")
		if !result.Stopped {
			return core2.MakeSingleRsp(fmt.Sprintf("MANIC MINUTE is not active. Current trigger is '%s'.", result.NextWord))
		}

		if result.ChannelID != msg.ChannelID {
			manicMinuteState.queueConclusion(result.ChannelID, result.ServerID)
		}

		return core2.MakeSingleRsp(fmt.Sprintf("Stopped MANIC MINUTE on '%s'. Next trigger is '%s'.", result.PreviousWord, result.NextWord))
	case "test":
		if manicMinuteState.isActive() {
			return core2.MakeSingleRsp("MANIC MINUTE is already active.")
		}

		word := ""
		if len(pieces) > 2 {
			word = strings.TrimSpace(strings.Join(pieces[2:], " "))
			if manicMinuteState.setTriggerWord("admin-test", word) == "" {
				return core2.MakeSingleRsp("Provide a normal word for !manic test, for example: !manic test bitcoin")
			}
		}

		_, started := manicMinuteState.start(msg.ChannelID, msg.ServerID, msg.ID, msg.UserID, true)
		if !started {
			if manicMinuteState.isActive() {
				return core2.MakeSingleRsp("MANIC MINUTE is already active.")
			}

			return core2.MakeSingleRsp("Could not start MANIC MINUTE right now.")
		}

		return core2.MakeSingleRsp(manicMinuteGenerate())
	default:
		return core2.MakeSingleRsp("Usage: !manic <status|stats|test|stop|rotate>")
	}
}

func renderManicMinuteStatus(svc *core2.Service, msg *core2.IncomingMessage) string {
	status := manicMinuteState.statusSnapshot()
	channelName := status.ChannelID
	if svc != nil && status.ChannelID != "" && status.ServerID != "" {
		if resolved := svc.Driver().GetChannelName(status.ChannelID, status.ServerID); resolved != "" {
			channelName = resolved
		}
	}

	channelCooldown := false
	wordCooldown := false
	if manicMinuteUsesPersistence() {
		channelCooldown = db.IsManicMinuteChannelOnCooldown(msg.ServerID, msg.ChannelID)
		wordCooldown = db.IsManicMinuteWordOnCooldown(msg.ServerID, status.TriggerWord)
	}

	parts := []string{
		fmt.Sprintf("active=%t", status.Active),
		fmt.Sprintf("cooldown=%dm", manicMinuteCooldownMinutes),
		fmt.Sprintf("channel_cooldown=%t", channelCooldown),
		fmt.Sprintf("word_cooldown=%t", wordCooldown),
	}

	if status.Active {
		remaining := time.Until(status.EndsAt).Round(time.Second)
		if remaining < 0 {
			remaining = 0
		}

		parts = append(parts,
			fmt.Sprintf("channel='%s'", channelName),
			fmt.Sprintf("remaining=%s", remaining),
			fmt.Sprintf("messages=%d", status.MessageCount),
		)
	}

	if manicMinuteUsesPersistence() {
		stats, err := db.GetManicMinuteStats(msg.ServerID)
		if err == nil {
			parts = append(parts,
				fmt.Sprintf("events_24h=%d", stats.Events24h),
				fmt.Sprintf("events_7d=%d", stats.Events7d),
			)
		}
	}

	return "MANIC MINUTE " + strings.Join(parts, " ")
}

func renderManicMinuteStats(serverID string) string {
	if !manicMinuteUsesPersistence() {
		return "MANIC MINUTE stats are only available with a live database-backed runtime."
	}

	stats, err := db.GetManicMinuteStats(serverID)
	if err != nil {
		return fmt.Sprintf("Failed to load MANIC MINUTE stats: %v", err)
	}

	parts := []string{
		fmt.Sprintf("total=%d", stats.TotalEvents),
		fmt.Sprintf("events_24h=%d", stats.Events24h),
		fmt.Sprintf("events_7d=%d", stats.Events7d),
	}

	if stats.MostRecent != nil {
		parts = append(parts, fmt.Sprintf("last='%s'/%s", stats.MostRecent.TriggerWord, stats.MostRecent.Status))
	}

	if len(stats.TopWords) > 0 {
		top := make([]string, 0, len(stats.TopWords))
		for _, word := range stats.TopWords {
			top = append(top, fmt.Sprintf("%s:%d", word.Word, word.Count))
		}

		parts = append(parts, "top="+strings.Join(top, ","))
	}

	return "MANIC MINUTE stats " + strings.Join(parts, " ")
}
