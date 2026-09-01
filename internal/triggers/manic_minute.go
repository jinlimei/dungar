package triggers

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"unicode"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

const (
	manicMinuteDuration        = 1 * time.Minute
	manicMinuteTick            = 5 * time.Second
	manicMinuteCooldownMinutes = 30
	manicMinuteStartChance     = 1.0
	manicMinuteManualMinRunes  = 2
	manicMinuteAutoMinRunes    = 4
)

var manicMinuteStopPhrases = []string{
	"shut up",
	"stfu",
	"be quiet",
	"stop talking",
	"stop posting",
}

type manicMinuteManager struct {
	lock sync.Mutex

	initialized      bool
	restoreAttempted bool
	active           bool

	triggerWord string

	pendingConclusion *manicMinuteScheduleSnapshot

	activeChannelID         string
	activeServerID          string
	activeEventID           int64
	activeTriggerMessageID  string
	activeTriggeredByUserID string
	endsAt                  time.Time
	cooldownUntil           time.Time
	messageCount            int
	startMissCount          int

	version uint64
}

type manicMinuteStatus struct {
	TriggerWord  string
	Active       bool
	ChannelID    string
	ServerID     string
	EndsAt       time.Time
	MessageCount int
	StartChance  float64
	StartMisses  int
}

type manicMinuteScheduleSnapshot struct {
	ChannelID string
	ServerID  string
	Version   uint64
	Concluded bool
}

var (
	manicMinuteState = &manicMinuteManager{}

	manicMinuteTriggerWordPicker = defaultManicMinuteTriggerWordPicker
)

// InitializeManicMinute selects the current trigger word for startup.
func InitializeManicMinute() {
	manicMinuteState.initialize()
}

func manicMinuteUsesPersistence() bool {
	return core != nil && core.DriverName() != "mock" && !utils.InTestEnv()
}

func (mm *manicMinuteManager) initialize() {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.ensureInitializedLocked("startup")
}

func (mm *manicMinuteManager) ensureInitializedLocked(reason string) {
	if !mm.initialized {
		mm.initialized = true
	}

	if mm.triggerWord == "" {
		mm.restoreRuntimeStateLocked()
	}

	if mm.triggerWord == "" {
		mm.selectNewTriggerWordLocked(reason)
	}
}

// restoreRuntimeStateLocked reloads the persisted trigger word and cooldown so
// a restart continues where the last run left off. An interrupted active
// minute is not resumed; the event is treated as concluded.
func (mm *manicMinuteManager) restoreRuntimeStateLocked() {
	if mm.restoreAttempted || !manicMinuteUsesPersistence() {
		return
	}

	mm.restoreAttempted = true

	state, err := db.GetManicMinuteRuntimeState(utils.ProtocolMode())
	if err != nil {
		log.Printf("[ManicMinute] failed to restore runtime state: %v", err)
		return
	}

	if state == nil {
		return
	}

	word := sanitizeManicMinuteTriggerWord(state.TriggerWord)
	if word == "" {
		return
	}

	mm.triggerWord = word
	mm.startMissCount = 0
	if state.HasCooldown && state.CooldownUntil.After(time.Now()) {
		mm.cooldownUntil = state.CooldownUntil
	}
	mm.version++
	mm.persistRuntimeStateLocked("startup-restore")

	log.Printf(
		"[ManicMinute] Restored trigger word='%s' from persisted runtime state",
		word,
	)
}

func (mm *manicMinuteManager) persistRuntimeStateLocked(reason string) {
	if !manicMinuteUsesPersistence() {
		return
	}

	err := db.UpsertManicMinuteRuntimeState(&db.ManicMinuteRuntimeState{
		ProtocolDriver:  utils.ProtocolMode(),
		TriggerWord:     mm.triggerWord,
		StartChance:     mm.currentStartChanceLocked(),
		Active:          mm.active,
		ActiveServerID:  mm.activeServerID,
		ActiveChannelID: mm.activeChannelID,
		HasCooldown:     !mm.cooldownUntil.IsZero(),
		CooldownUntil:   mm.cooldownUntil,
		UpdatedReason:   reason,
		UpdatedAt:       time.Now().UTC(),
	})
	if err != nil {
		log.Printf("[ManicMinute] failed to persist runtime state reason='%s': %v", reason, err)
	}
}

func (mm *manicMinuteManager) selectNewTriggerWordLocked(reason string) string {
	previous := mm.triggerWord
	fallback := ""
	next := ""

	for attempt := 0; attempt < 20; attempt++ {
		word := sanitizeManicMinuteAutoTriggerWord(manicMinuteTriggerWordPicker())
		if word == "" {
			continue
		}

		if fallback == "" {
			fallback = word
		}

		if !strings.EqualFold(word, previous) {
			next = word
			break
		}
	}

	if next == "" {
		next = fallback
	}

	if next == "" {
		next = sanitizeManicMinuteAutoTriggerWord(markovPickWord())
	}

	if next == "" {
		next = previous
	}

	mm.triggerWord = next
	mm.startMissCount = 0
	mm.version++
	mm.persistRuntimeStateLocked(reason)

	if next == "" {
		log.Printf("[ManicMinute] No trigger word available after '%s'", reason)
	} else {
		log.Printf(
			"[ManicMinute] Selected trigger word='%s' reason='%s' previous='%s'",
			next,
			reason,
			previous,
		)
	}

	return next
}

func (mm *manicMinuteManager) rotate(reason string) (string, string) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.ensureInitializedLocked(reason)

	previous := mm.triggerWord
	next := mm.selectNewTriggerWordLocked(reason)

	return previous, next
}

func (mm *manicMinuteManager) setTriggerWord(reason, word string) string {
	word = sanitizeManicMinuteTriggerWord(word)
	if word == "" {
		return ""
	}

	mm.lock.Lock()
	defer mm.lock.Unlock()

	previous := mm.triggerWord
	mm.triggerWord = word
	mm.startMissCount = 0
	mm.version++
	mm.persistRuntimeStateLocked(reason)

	log.Printf(
		"[ManicMinute] Forced trigger word='%s' reason='%s' previous='%s'",
		word,
		reason,
		previous,
	)

	return word
}

func (mm *manicMinuteManager) currentTriggerWord() string {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.ensureInitializedLocked("current-trigger")
	return mm.triggerWord
}

func (mm *manicMinuteManager) isActive() bool {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	return mm.active
}

func (mm *manicMinuteManager) statusSnapshot() manicMinuteStatus {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.ensureInitializedLocked("status")

	return manicMinuteStatus{
		TriggerWord:  mm.triggerWord,
		Active:       mm.active,
		ChannelID:    mm.activeChannelID,
		ServerID:     mm.activeServerID,
		EndsAt:       mm.endsAt,
		MessageCount: mm.messageCount,
		StartChance:  mm.currentStartChanceLocked(),
		StartMisses:  mm.startMissCount,
	}
}

func (mm *manicMinuteManager) currentStartChanceLocked() float64 {
	return manicMinuteStartChance
}

func pickManicMinuteCooldownDuration() time.Duration {
	return time.Duration(manicMinuteCooldownMinutes) * time.Minute
}

func (mm *manicMinuteManager) start(channelID, serverID, triggerMessageID, triggeredByUserID string, bypassCooldown bool) (string, bool) {
	if channelID == "" || serverID == "" {
		return "", false
	}

	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.ensureInitializedLocked("start")

	if mm.active || mm.triggerWord == "" {
		return "", false
	}

	if manicMinuteUsesPersistence() && !bypassCooldown {
		if db.IsManicMinuteChannelOnCooldown(serverID, channelID) {
			log.Printf(
				"[ManicMinute] channel cooldown blocked channelID='%s' trigger='%s'",
				channelID,
				mm.triggerWord,
			)
			return "", false
		}

		if db.IsManicMinuteWordOnCooldown(serverID, mm.triggerWord) {
			log.Printf(
				"[ManicMinute] word cooldown blocked trigger='%s' serverID='%s'; rotating",
				mm.triggerWord,
				serverID,
			)
			mm.selectNewTriggerWordLocked("word-cooldown")
			return "", false
		}
	}

	eventID := int64(0)
	startedAt := time.Now().UTC()
	cooldownDuration := pickManicMinuteCooldownDuration()
	cooldownUntil := startedAt.Add(cooldownDuration)
	if manicMinuteUsesPersistence() {
		startedEventID, err := db.StartManicMinuteEvent(
			serverID,
			channelID,
			mm.triggerWord,
			triggerMessageID,
			triggeredByUserID,
			startedAt,
			cooldownUntil,
		)
		if err != nil {
			log.Printf("[ManicMinute] failed to persist start trigger='%s': %v", mm.triggerWord, err)
		} else {
			eventID = startedEventID
		}
	}

	mm.active = true
	mm.activeChannelID = channelID
	mm.activeServerID = serverID
	mm.activeEventID = eventID
	mm.activeTriggerMessageID = triggerMessageID
	mm.activeTriggeredByUserID = triggeredByUserID
	mm.endsAt = startedAt.Add(manicMinuteDuration)
	mm.cooldownUntil = cooldownUntil
	mm.messageCount = 1
	mm.startMissCount = 0
	mm.version++
	mm.persistRuntimeStateLocked("start")

	log.Printf(
		"[ManicMinute] Started channelID='%s' serverID='%s' trigger='%s' endsAt=%s cooldown=%s bypassCooldown=%t",
		channelID,
		serverID,
		mm.triggerWord,
		mm.endsAt.Format(time.RFC3339),
		cooldownDuration,
		bypassCooldown,
	)

	return mm.triggerWord, true
}

func manicMinuteEventStatus(reason string) string {
	if strings.Contains(reason, "stop") {
		return "stopped"
	}

	return "completed"
}

type manicMinuteStopResult struct {
	PreviousWord string
	NextWord     string
	ChannelID    string
	ServerID     string
	Stopped      bool
}

func (mm *manicMinuteManager) stop(reason string) manicMinuteStopResult {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.ensureInitializedLocked("stop")

	if !mm.active {
		return manicMinuteStopResult{
			PreviousWord: mm.triggerWord,
			NextWord:     mm.triggerWord,
		}
	}

	previous := mm.triggerWord
	channelID := mm.activeChannelID
	serverID := mm.activeServerID
	eventID := mm.activeEventID
	messageCount := mm.messageCount

	mm.active = false
	mm.activeChannelID = ""
	mm.activeServerID = ""
	mm.activeEventID = 0
	mm.activeTriggerMessageID = ""
	mm.activeTriggeredByUserID = ""
	mm.endsAt = time.Time{}
	mm.messageCount = 0

	if manicMinuteUsesPersistence() && eventID > 0 {
		if err := db.CompleteManicMinuteEvent(
			eventID,
			manicMinuteEventStatus(reason),
			reason,
			time.Now().UTC(),
			messageCount,
		); err != nil {
			log.Printf("[ManicMinute] failed to persist stop id=%d reason='%s': %v", eventID, reason, err)
		}
	}

	next := mm.selectNewTriggerWordLocked(reason)

	log.Printf(
		"[ManicMinute] Stopped reason='%s' previous='%s' next='%s'",
		reason,
		previous,
		next,
	)

	return manicMinuteStopResult{
		PreviousWord: previous,
		NextWord:     next,
		ChannelID:    channelID,
		ServerID:     serverID,
		Stopped:      true,
	}
}

// queueConclusion asks the scheduler to deliver the conclusion lines to the
// given channel on its next tick. Used when a minute is stopped from a
// different channel than the one it is running in.
func (mm *manicMinuteManager) queueConclusion(channelID, serverID string) {
	if channelID == "" || serverID == "" {
		return
	}

	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.pendingConclusion = &manicMinuteScheduleSnapshot{
		ChannelID: channelID,
		ServerID:  serverID,
		Concluded: true,
	}
}

func (mm *manicMinuteManager) snapshotForScheduledMessage(now time.Time) (*manicMinuteScheduleSnapshot, bool) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.ensureInitializedLocked("schedule")

	if mm.pendingConclusion != nil {
		snapshot := mm.pendingConclusion
		mm.pendingConclusion = nil
		return snapshot, true
	}

	if !mm.active {
		return nil, false
	}

	if !mm.endsAt.IsZero() && now.After(mm.endsAt) {
		previous := mm.triggerWord
		channelID := mm.activeChannelID
		serverID := mm.activeServerID
		eventID := mm.activeEventID
		messageCount := mm.messageCount

		mm.active = false
		mm.activeChannelID = ""
		mm.activeServerID = ""
		mm.activeEventID = 0
		mm.activeTriggerMessageID = ""
		mm.activeTriggeredByUserID = ""
		mm.endsAt = time.Time{}
		mm.messageCount = 0

		if manicMinuteUsesPersistence() && eventID > 0 {
			if err := db.CompleteManicMinuteEvent(
				eventID,
				"completed",
				"minute-ended",
				now.UTC(),
				messageCount,
			); err != nil {
				log.Printf("[ManicMinute] failed to persist completion id=%d: %v", eventID, err)
			}
		}

		next := mm.selectNewTriggerWordLocked("minute-ended")

		log.Printf(
			"[ManicMinute] Finished naturally previous='%s' next='%s'",
			previous,
			next,
		)

		return &manicMinuteScheduleSnapshot{
			ChannelID: channelID,
			ServerID:  serverID,
			Concluded: true,
		}, true
	}

	return &manicMinuteScheduleSnapshot{
		ChannelID: mm.activeChannelID,
		ServerID:  mm.activeServerID,
		Version:   mm.version,
	}, true
}

func (mm *manicMinuteManager) recordOutput(version uint64, channelID, serverID string) bool {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	if !mm.active ||
		mm.version != version ||
		mm.activeChannelID != channelID ||
		mm.activeServerID != serverID {
		return false
	}

	mm.messageCount++
	return true
}

func manicMinuteTriggeringUserName(svc *core2.Service, msg *core2.IncomingMessage) string {
	if svc != nil && msg != nil {
		if name := strings.TrimSpace(svc.Driver().GetUserName(msg.UserID, msg.ServerID)); name != "" {
			return name
		}
	}

	if msg != nil && strings.TrimSpace(msg.UserID) != "" {
		return strings.TrimSpace(msg.UserID)
	}

	return "Someone"
}

func manicMinuteStartResponses(triggeredBy string, word string) []*core2.Response {
	responses := []*core2.Response{
		core2.MakeRsp(fmt.Sprintf("%s triggered a dungarmatic manic minute by saying %q.", triggeredBy, word)),
	}

	if contents := manicMinuteGenerate(); contents != "" {
		responses = append(responses, core2.MakeRsp(contents))
	}

	return responses
}

func manicMinuteConclusionContents() []string {
	return []string{
		"This concludes dungarmatic's manic minute.",
		"New trigger word set.",
	}
}

func manicMinuteConclusionResponses() []*core2.Response {
	lines := manicMinuteConclusionContents()
	responses := make([]*core2.Response, 0, len(lines))

	for _, line := range lines {
		responses = append(responses, core2.MakeRsp(line))
	}

	return responses
}

func manicMinuteConclusionScheduledMessages(channelID, serverID string, sentAt time.Time) []*core2.ScheduledMessage {
	lines := manicMinuteConclusionContents()
	messages := make([]*core2.ScheduledMessage, 0, len(lines))

	for _, line := range lines {
		messages = append(messages, &core2.ScheduledMessage{
			ChannelID: channelID,
			ServerID:  serverID,
			Contents:  line,
			SentAt:    sentAt,
		})
	}

	return messages
}

func manicMinuteHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if svc == nil || msg == nil || msg.UserID == "" {
		return core2.EmptyRsp()
	}

	if svc.GetBotUser().ID != "" && msg.UserID == svc.GetBotUser().ID {
		return core2.EmptyRsp()
	}

	InitializeManicMinute()

	if isManicMinuteStopMessage(svc, msg) {
		if !canUseDiscordAdminCommand(svc, msg) {
			return core2.EmptyRsp()
		}

		result := manicMinuteState.stop("directed-stop")
		if !result.Stopped {
			return core2.EmptyRsp()
		}

		log.Printf(
			"[ManicMinute] Directed stop by userID='%s' previous='%s' next='%s'",
			msg.UserID,
			result.PreviousWord,
			result.NextWord,
		)

		if result.ChannelID != msg.ChannelID {
			manicMinuteState.queueConclusion(result.ChannelID, result.ServerID)
			return core2.EmptyRsp()
		}

		return manicMinuteConclusionResponses()
	}

	if !canTriggerManicMinute(svc, msg) || manicMinuteState.isActive() {
		return core2.EmptyRsp()
	}

	triggerWord := manicMinuteState.currentTriggerWord()
	if triggerWord == "" || !messageContainsTriggerWord(msg.Contents, triggerWord) {
		return core2.EmptyRsp()
	}

	word, started := manicMinuteState.start(msg.ChannelID, msg.ServerID, msg.ID, msg.UserID, false)
	if !started {
		return core2.EmptyRsp()
	}

	log.Printf(
		"[ManicMinute] Triggered by userID='%s' channelID='%s' word='%s'",
		msg.UserID,
		msg.ChannelID,
		word,
	)

	return manicMinuteStartResponses(manicMinuteTriggeringUserName(svc, msg), word)
}

func manicMinuteScheduler(_ *core2.Service) []*core2.ScheduledMessage {
	now := time.Now()
	snapshot, ok := manicMinuteState.snapshotForScheduledMessage(now)
	if !ok {
		return nil
	}

	if snapshot.Concluded {
		return manicMinuteConclusionScheduledMessages(snapshot.ChannelID, snapshot.ServerID, now)
	}

	contents := manicMinuteGenerate()
	if contents == "" {
		return nil
	}

	if !manicMinuteState.recordOutput(snapshot.Version, snapshot.ChannelID, snapshot.ServerID) {
		return nil
	}

	return []*core2.ScheduledMessage{
		{
			ChannelID: snapshot.ChannelID,
			ServerID:  snapshot.ServerID,
			Contents:  contents,
			SentAt:    now,
		},
	}
}

func defaultManicMinuteTriggerWordPicker() string {
	if !manicMinuteUsesPersistence() {
		return markovPickWord()
	}

	source := db.RandomRawMessageForProtocol(
		utils.ProtocolMode(),
		utils.DiscordAllowedLearningChannelIDs(),
	)

	if source != "" {
		words := buildManicMinuteAutomaticCandidateWords(source)
		if len(words) > 0 {
			return random.PickString(words)
		}
	}

	return markovPickWord()
}

func buildManicMinuteAutomaticCandidateWords(str string) []string {
	return buildManicMinuteCandidateWords(str, sanitizeManicMinuteAutoTriggerWord)
}

func buildManicMinuteMessageCandidateWords(str string) []string {
	return buildManicMinuteCandidateWords(str, sanitizeManicMinuteTriggerWord)
}

func buildManicMinuteCandidateWords(str string, sanitizer func(string) string) []string {
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

func sanitizeManicMinuteTriggerWord(word string) string {
	return sanitizeManicMinuteTriggerWordWithMinLen(word, manicMinuteManualMinRunes)
}

func sanitizeManicMinuteAutoTriggerWord(word string) string {
	return sanitizeManicMinuteTriggerWordWithMinLen(word, manicMinuteAutoMinRunes)
}

func sanitizeManicMinuteTriggerWordWithMinLen(word string, minRunes int) string {
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

func messageContainsTriggerWord(contents, triggerWord string) bool {
	triggerWord = sanitizeManicMinuteTriggerWord(triggerWord)
	if triggerWord == "" {
		return false
	}

	for _, word := range buildManicMinuteMessageCandidateWords(contents) {
		if word == triggerWord {
			return true
		}
	}

	return false
}

func canTriggerManicMinute(svc *core2.Service, msg *core2.IncomingMessage) bool {
	contents := strings.TrimSpace(msg.Contents)
	if contents == "" {
		return false
	}

	if strings.HasPrefix(contents, "!") {
		return false
	}

	if svc != nil && svc.DriverName() == "discord" && !shouldRecordDiscordRawMessage(msg.ChannelID) {
		return false
	}

	return true
}

func isManicMinuteStopMessage(svc *core2.Service, msg *core2.IncomingMessage) bool {
	if svc == nil || msg == nil || !isDirectedAtDungar(svc, msg) {
		return false
	}

	contents := strings.ToLower(normalizeDirectedContents(svc, msg.ServerID, msg.Contents, svc.GetBotUser()))
	contents = utils.CleanSpaces(contents)

	for _, phrase := range manicMinuteStopPhrases {
		if strings.Contains(contents, phrase) {
			return true
		}
	}

	return false
}

// manicMinuteGenerate produces a burst of random Markov output. The text is
// intentionally unrelated to the trigger word: the word only starts the
// minute, it does not theme it.
func manicMinuteGenerate() string {
	for attempt := 0; attempt < 10; attempt++ {
		seed := strings.TrimSpace(markovPickWord())
		if seed == "" {
			continue
		}

		if output := strings.TrimSpace(markovGenerate(seed)); output != "" {
			return output
		}
	}

	return ""
}
