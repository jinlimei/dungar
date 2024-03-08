package dcli

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gopkg.in/ini.v1"
)

// PrintBotInfo provides bot info to the CLI based off of settings inis
func PrintBotInfo() {
	cfg, err := ini.Load("settings.ini")
	utils.HaltingError("ini loading failed", err)

	secretCfg, err := ini.Load(cfg.Section("base").Key("secrets_file").String())
	utils.HaltingError("runner secrets file", err)

	api := slack.New(
		secretCfg.Section("slack").Key("bot_user_access_token").String(),
	)

	users, _ := api.GetUsers()

	for _, user := range users {
		if user.IsBot && user.RealName == "Dungar" {
			bot, _ := api.GetBotInfo(user.Profile.BotID)

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
