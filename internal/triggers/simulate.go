package triggers

import (
	"fmt"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/markov3"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var userSimulationBuilder = defaultUserSimulationBuilder

func simulateUserHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if svc == nil || msg == nil || !isDirectedAtDungar(svc, msg) {
		return core2.EmptyRsp()
	}

	contents := strings.TrimSpace(normalizeDirectedContents(svc, msg.ServerID, msg.Contents, svc.GetBotUser()))
	if !strings.HasPrefix(strings.ToLower(contents), "simulate") {
		return core2.EmptyRsp()
	}

	if svc.DriverName() != "discord" && svc.DriverName() != "mock" && !utils.InTestEnv() {
		return core2.MakeSingleRsp("Simulation is only implemented for Discord right now.")
	}

	targetText := strings.TrimSpace(strings.TrimPrefix(contents, strings.Fields(contents)[0]))
	targetUser, err := resolveSimulationTargetUser(svc, msg, targetText)
	if err != nil {
		return core2.MakeSingleRsp(fmt.Sprintf("Could not resolve simulation target: %v", err))
	}

	if simulateSelfTarget(svc, targetUser) {
		output := strings.TrimSpace(simulateDungarmaticOutput())
		if output == "" {
			return core2.MakeSingleRsp("Dungarmatic does not know itself well enough yet.")
		}

		return core2.MakeSingleRsp(output)
	}

	output, learned, err := userSimulationBuilder(targetUser, msg.ServerID)
	if err != nil {
		return core2.MakeSingleRsp(fmt.Sprintf("Failed to simulate @%s: %v", targetUser.Name, err))
	}

	if learned <= 0 || strings.TrimSpace(output) == "" {
		return core2.MakeSingleRsp(fmt.Sprintf("I don't know enough about @%s yet.", targetUser.Name))
	}

	return core2.MakeSingleRsp(output)
}

func resolveSimulationTargetUser(svc *core2.Service, msg *core2.IncomingMessage, targetText string) (core2.User, error) {
	ignoreUserIDs := []string{svc.GetBotUser().ID}

	if userID := extractTargetUserID(msg, ignoreUserIDs...); userID != "" {
		return resolveUserByIDOrFallback(svc, userID, msg.ServerID, targetText), nil
	}

	fields := strings.Fields(strings.TrimSpace(targetText))
	if len(fields) == 0 {
		return core2.User{}, fmt.Errorf("usage: @%s simulate @user", svc.GetBotUser().Name)
	}

	targetText = fields[0]
	targetText = strings.TrimSpace(strings.TrimPrefix(targetText, "@"))
	if targetText == "" {
		return core2.User{}, fmt.Errorf("usage: @%s simulate @user", svc.GetBotUser().Name)
	}

	user, err := resolveUserByNickFromCache(svc, msg.ServerID, targetText)
	if err != nil {
		return core2.User{}, fmt.Errorf("could not find user '%s'", targetText)
	}

	return user, nil
}

func simulateSelfTarget(svc *core2.Service, target core2.User) bool {
	if svc == nil {
		return false
	}

	bot := svc.GetBotUser()
	return bot.ID != "" && target.ID == bot.ID
}

func simulateDungarmaticOutput() string {
	seed := strings.TrimSpace(markovPickWord())
	if seed == "" {
		return ""
	}

	return strings.TrimSpace(markovGenerate(seed))
}

func defaultUserSimulationBuilder(target core2.User, serverID string) (string, int, error) {
	if !strings.EqualFold(utils.ProtocolMode(), "discord") && !utils.InTestEnv() {
		return "", 0, fmt.Errorf("simulation is only implemented for Discord")
	}

	model := markov3.MakeMarkov(markov3.MarkovSpaceID("simulate-" + target.ID))
	learned, err := model.LearnFromDiscordRawMessagesByAuthor(serverID, target.ID, utils.DiscordAllowedLearningChannelIDs())
	if err != nil {
		return "", learned, err
	}

	if learned <= 0 {
		return "", learned, nil
	}

	return generateUserSimulation(model), learned, nil
}

func generateUserSimulation(model *markov3.Markov) string {
	if model == nil {
		return ""
	}

	candidates := make([]string, 0, len(model.RevWords))
	for _, word := range model.RevWords {
		word = strings.TrimSpace(word)
		if word == "" || strings.HasPrefix(word, "\x00") {
			continue
		}

		candidates = append(candidates, word)
	}

	for attempts := 0; attempts < 10 && len(candidates) > 0; attempts++ {
		seed := random.PickString(candidates)
		output := strings.TrimSpace(model.Generate(seed))
		if output != "" {
			return output
		}
	}

	return ""
}
