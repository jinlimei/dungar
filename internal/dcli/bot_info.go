package dcli

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// PrintBotInfo provides bot info to the CLI based off of settings inis
func PrintBotInfo() {
	utils.LoadSettingsAndSecrets()

	switch utils.ProtocolMode() {
	case "slack":
		printSlackBotInfo()
	case "discord":
		printDiscordBotInfo()
	default:
		log.Fatalf("unknown protocol mode for bot-info: %s", utils.ProtocolMode())
	}
}

func printSlackBotInfo() {
	api := slack.New(utils.SlackAccessToken())

	users, err := api.GetUsers()
	utils.HaltingError("slack get users failed", err)

	for _, user := range users {
		if user.IsBot && user.RealName == "Dungar" {
			bot, _ := api.GetBotInfo(slack.GetBotInfoParameters{
				Bot: user.Profile.BotID,
			})

			spew.Dump(user, bot)
		}
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("")
	fmt.Println("")

	//chans, _ := api.GetChannels(true)
	//spew.Dump(chans)

	output, err := api.GetUserGroups()

	if err != nil {
		panic(err)
	}

	for _, chn := range output {
		fmt.Printf("ID: %s, Name: %s, Is Group? %v, Preferences? %v\n", chn.ID, chn.Name, chn.IsUserGroup, chn.Prefs)
	}
}

func printDiscordBotInfo() {
	api, err := discordgo.New("Bot " + utils.DiscordAccessToken())
	utils.HaltingError("discord api init failed", err)

	user, err := api.User("@me")
	utils.HaltingError("discord get current user failed", err)

	spew.Dump(user)

	guildName := utils.DiscordGuildName()
	if guildName == "" {
		fmt.Println("Configured guild name: <unset>")
		return
	}

	fmt.Printf("Configured guild name: %s\n", guildName)
}
