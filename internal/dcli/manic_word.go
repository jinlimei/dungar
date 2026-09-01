package dcli

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

const manicMinuteCLIBaseStartChance = 1.0
const (
	manicMinuteCLIManualMinRunes = 2
	manicMinuteCLIAutoMinRunes   = 4
)

// PrintManicMinuteWord prints the currently armed manic-minute trigger word for
// the configured protocol mode.
func PrintManicMinuteWord() {
	HandleManicWordCommand(nil)
}

// HandleManicWordCommand shows, rotates, or explicitly sets the persisted
// manic-minute trigger word for the configured protocol mode.
func HandleManicWordCommand(args []string) {
	utils.LoadSettingsAndSecrets()
	db.ConnectToDatabase()

	mode, err := normalizeManicMinuteProtocolMode(utils.ProtocolMode())
	if err != nil {
		log.Fatal(err)
	}

	switch len(args) {
	case 0:
		printManicMinuteWord(mode)
		return
	default:
		switch strings.ToLower(strings.TrimSpace(args[0])) {
		case "", "show":
			printManicMinuteWord(mode)
			return
		case "rotate":
			rotateManicMinuteWord(mode)
			return
		case "set":
			if len(args) <= 1 {
				log.Fatal("usage: dungar manic-word set <word>")
			}

			setManicMinuteWord(mode, strings.Join(args[1:], " "))
			return
		default:
			log.Fatalf("usage: dungar manic-word [show|rotate|set <word>]")
		}
	}
}

func printManicMinuteWord(mode string) {
	state, err := db.GetManicMinuteRuntimeState(mode)
	if err != nil {
		log.Fatal(err)
	}

	if state == nil || strings.TrimSpace(state.TriggerWord) == "" {
		log.Fatalf("no persisted manic minute trigger word for mode=%s; start dungar run first", mode)
	}

	fmt.Println(formatManicMinuteWord(state, time.Now().UTC()))
}

func rotateManicMinuteWord(mode string) {
	state, err := currentManicMinuteRuntimeState(mode)
	if err != nil {
		log.Fatal(err)
	}

	if state.Active {
		log.Fatal("cannot rotate manic-word while a manic minute is active; stop or restart dungar first")
	}

	next, err := nextRandomManicMinuteWord(mode, state.TriggerWord)
	if err != nil {
		log.Fatal(err)
	}

	previous := state.TriggerWord
	state.TriggerWord = next
	state.StartChance = manicMinuteCLIBaseStartChance
	state.UpdatedReason = "cli-rotate"
	state.UpdatedAt = time.Now().UTC()

	if err := db.UpsertManicMinuteRuntimeState(state); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rotated manic-word: %s -> %s\n", previous, next)
	fmt.Println("Restart dungar if it is already running to load the new trigger word.")
}

func setManicMinuteWord(mode, word string) {
	state, err := currentManicMinuteRuntimeState(mode)
	if err != nil {
		log.Fatal(err)
	}

	if state.Active {
		log.Fatal("cannot set manic-word while a manic minute is active; stop or restart dungar first")
	}

	word = sanitizeManicMinuteCLITriggerWord(word)
	if word == "" {
		log.Fatal("provided manic-word is invalid; use a normal word with at least 2 characters")
	}

	previous := state.TriggerWord
	state.TriggerWord = word
	state.StartChance = manicMinuteCLIBaseStartChance
	state.UpdatedReason = "cli-set"
	state.UpdatedAt = time.Now().UTC()

	if err := db.UpsertManicMinuteRuntimeState(state); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Set manic-word: %s -> %s\n", previous, word)
	fmt.Println("Restart dungar if it is already running to load the new trigger word.")
}

func currentManicMinuteRuntimeState(mode string) (*db.ManicMinuteRuntimeState, error) {
	state, err := db.GetManicMinuteRuntimeState(mode)
	if err != nil {
		return nil, err
	}

	if state == nil {
		return &db.ManicMinuteRuntimeState{
			ProtocolDriver: mode,
			StartChance:    manicMinuteCLIBaseStartChance,
		}, nil
	}

	if strings.TrimSpace(state.ProtocolDriver) == "" {
		state.ProtocolDriver = mode
	}

	return state, nil
}

func nextRandomManicMinuteWord(mode, current string) (string, error) {
	current = sanitizeManicMinuteCLITriggerWord(current)
	fallback := ""

	for attempt := 0; attempt < 50; attempt++ {
		source := db.RandomRawMessageForProtocol(mode, utils.DiscordAllowedLearningChannelIDs())
		if strings.TrimSpace(source) == "" {
			continue
		}

		words := buildManicMinuteCLIAutomaticCandidateWords(source)
		if len(words) == 0 {
			continue
		}

		if fallback == "" {
			fallback = random.PickString(words)
		}

		candidates := make([]string, 0, len(words))
		for _, word := range words {
			if current != "" && strings.EqualFold(word, current) {
				continue
			}

			candidates = append(candidates, word)
		}

		if len(candidates) > 0 {
			return random.PickString(candidates), nil
		}
	}

	if fallback != "" && !strings.EqualFold(fallback, current) {
		return fallback, nil
	}

	return "", fmt.Errorf("could not find a different manic-word from stored messages")
}

func formatManicMinuteWord(state *db.ManicMinuteRuntimeState, now time.Time) string {
	cooldownLabel := "none"
	if state != nil && state.HasCooldown && !state.CooldownUntil.IsZero() {
		remaining := state.CooldownUntil.Sub(now).Round(time.Second)
		if remaining > 0 {
			cooldownLabel = remaining.String()
		}
	}

	return fmt.Sprintf(
		"Current manic-word: %s (cooldown remaining: %s)",
		state.TriggerWord,
		cooldownLabel,
	)
}

func buildManicMinuteCLIAutomaticCandidateWords(str string) []string {
	return buildManicMinuteCLICandidateWords(str, sanitizeManicMinuteCLIAutoTriggerWord)
}

func buildManicMinuteCLICandidateWords(str string, sanitizer func(string) string) []string {
	words := utils.StringToWords(str, true)
	out := make([]string, 0, len(words))
	seen := make(map[string]struct{}, len(words))

	for _, word := range words {
		word = sanitizer(word)
		if word == "" {
			continue
		}

		if _, ok := seen[word]; ok {
			continue
		}

		seen[word] = struct{}{}
		out = append(out, word)
	}

	return out
}

func sanitizeManicMinuteCLITriggerWord(word string) string {
	return sanitizeManicMinuteCLITriggerWordWithMinLen(word, manicMinuteCLIManualMinRunes)
}

func sanitizeManicMinuteCLIAutoTriggerWord(word string) string {
	return sanitizeManicMinuteCLITriggerWordWithMinLen(word, manicMinuteCLIAutoMinRunes)
}

func sanitizeManicMinuteCLITriggerWordWithMinLen(word string, minRunes int) string {
	word = strings.TrimSpace(word)
	if word == "" {
		return ""
	}

	if utils.IsURL(word) {
		return ""
	}

	word = utils.Normalize(word)
	word = strings.TrimSpace(word)

	if len([]rune(word)) < minRunes {
		return ""
	}

	if strings.Contains(word, "@") || strings.Contains(word, "<") || strings.Contains(word, ">") {
		return ""
	}

	hasLetter := false
	for _, r := range word {
		if unicode.IsLetter(r) {
			hasLetter = true
			continue
		}

		if unicode.IsDigit(r) || r == '\'' || r == '-' || r == '_' {
			continue
		}

		return ""
	}

	if !hasLetter {
		return ""
	}

	return word
}

func normalizeManicMinuteProtocolMode(mode string) (string, error) {
	switch {
	case strings.EqualFold(mode, "discord"):
		return "discord", nil
	case strings.EqualFold(mode, "slack"):
		return "slack", nil
	default:
		return "", fmt.Errorf("unknown protocol mode for manic-word: %s", mode)
	}
}
