package utils

import (
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	// EnvCI is our name for the environmental variable IN_CI_ENV
	EnvCI = "IN_CI_ENV"
	// EnvCD is our name for the environmental variable IN_CD_ENV
	EnvCD = "IN_CD_ENV"
	// DiscordLearningLookbackDays is the fixed Discord history window used
	// for both learning queries and one-shot history backfills.
	DiscordLearningLookbackDays = 30
)

var (
	normalIniFile  *ini.File
	secretsIniFile *ini.File
)

// InTestEnv checks to see if we're in a test (non-db-backed) environment
func InTestEnv() bool {
	return MustUseEnvVars()
}

// MustUseEnvVars detects if we're in an environment which requires
// use of environment variables
func MustUseEnvVars() bool {
	return os.Getenv(EnvCI) != "" || os.Getenv(EnvCD) != ""
}

func doesFileExist(file string) bool {
	info, err := os.Stat(file)

	if os.IsNotExist(err) {
		return false
	}

	return info.Size() > 0
}

// IsConfiguredDebugMode returns whether or not Debug Mode is configured
func IsConfiguredDebugMode() bool {
	if MustUseEnvVars() {
		return true
	}

	return normalIniFile.Section("base").Key("debug").MustBool(false)
}

// IsSilentRunning returns whether or not Dungar should respond
// This places a hard limit on what dungar does in ~taut~ and ~accord~
func IsSilentRunning() bool {
	if MustUseEnvVars() {
		return true
	}

	return normalIniFile.Section("base").Key("silent_running").MustBool(false)
}

// LoadSettingsAndSecrets will grab settings ini and secrets ini from various locations.
// If MustUseEnvVars() is true, this does nothing
func LoadSettingsAndSecrets() {
	if MustUseEnvVars() {
		return
	}

	cfg, ok := TryIniLoad("settings.ini", "../settings.ini", "../../settings.ini")

	if !ok {
		log.Fatal("Could not find settings.ini in specified locations")
	}

	secret, ok := TryIniLoad("secrets.ini", "../secrets.ini", "../../secrets.ini")

	if !ok {
		log.Fatalf("Could not find secrets.ini in specified locations")
	}

	normalIniFile = cfg
	secretsIniFile = secret
}

// TryIniLoad will try to load an INI file from a list of locations
func TryIniLoad(paths ...string) (*ini.File, bool) {
	pickedFile := ""
	for _, file := range paths {
		if doesFileExist(file) {
			pickedFile = file
			break
		}
	}

	if pickedFile == "" {
		return nil, false
	}

	cfg, err := ini.Load(pickedFile)
	HaltingError("TryIniLoad "+pickedFile, err)
	return cfg, true
}

// SentryDSN will return the DSN for sentry
func SentryDSN() string {
	if MustUseEnvVars() {
		return os.Getenv("DUNGAR_SENTRY_DSN")
	}

	return secretsIniFile.Section("sentry").Key("dsn").String()
}

// PinsSqliteFile will return what the sqlite file is supposed to be
func PinsSqliteFile() string {
	if MustUseEnvVars() {
		return os.Getenv("DUNGAR_PINS_SQLITE_FILE")
	}

	return normalIniFile.Section("base").Key("pins_sqlite").String()
}

// ProtocolMode returns what protocol mode we should be using.
func ProtocolMode() string {
	if MustUseEnvVars() {
		return os.Getenv("DUNGAR_PROTOCOL_MODE")
	}

	return normalIniFile.Section("base").Key("mode").MustString("slack")
}

// DiscordGuildName returns the guild name of the discord we're working with
func DiscordGuildName() string {
	if MustUseEnvVars() {
		return os.Getenv("DUNGAR_GUILD_NAME")
	}

	return secretsIniFile.Section("discord").Key("guild_name").String()
}

// DiscordAccessToken returns the saved discord access token
func DiscordAccessToken() string {
	if MustUseEnvVars() {
		return os.Getenv("DUNGAR_USER_ACCESS_TOKEN")
	}

	return secretsIniFile.Section("discord").Key("token").String()
}

// DiscordAllowedOutputChannelIDs returns the channel IDs the bot is allowed
// to send messages or reactions to. An empty map means unrestricted output.
func DiscordAllowedOutputChannelIDs() map[string]struct{} {
	raw := ""

	if MustUseEnvVars() {
		raw = os.Getenv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS")
	} else if normalIniFile != nil {
		raw = normalIniFile.Section("discord").Key("allowed_output_channel_ids").String()
	}

	out := make(map[string]struct{})

	for _, channelID := range strings.Split(raw, ",") {
		channelID = strings.TrimSpace(channelID)

		if channelID == "" {
			continue
		}

		out[channelID] = struct{}{}
	}

	return out
}

// DiscordAllowedLearningChannelIDs returns the channel IDs the bot is allowed
// to learn from. Today this intentionally mirrors the outgoing allowlist so
// Discord speech and training stay scoped to the same channels.
func DiscordAllowedLearningChannelIDs() map[string]struct{} {
	return DiscordAllowedOutputChannelIDs()
}

// IsDiscordLearningChannelAllowed reports whether the bot may learn from the
// given channel. An empty allowlist means unrestricted learning.
func IsDiscordLearningChannelAllowed(channelID string) bool {
	allowed := DiscordAllowedLearningChannelIDs()

	if len(allowed) == 0 {
		return true
	}

	_, ok := allowed[channelID]
	return ok
}

// SlackAccessToken will return the slack access token
func SlackAccessToken() string {
	if MustUseEnvVars() {
		return os.Getenv("DUNGAR_USER_ACCESS_TOKEN")
	}

	return secretsIniFile.Section("slack").Key("bot_user_access_token").String()
}

// PinCredentials will return a map of the pin credentials
func PinCredentials() map[string]string {
	if MustUseEnvVars() {
		return map[string]string{
			"team": os.Getenv("DUNGAR_TEAM_ID"),
			"auth": os.Getenv("DUNGAR_PINS_AUTH"),
			"url":  os.Getenv("DUNGAR_PINS_URL"),
		}
	}

	sect := secretsIniFile.Section("pins")

	return map[string]string{
		"team": secretsIniFile.Section("slack").Key("team_id").String(),
		"auth": sect.Key("auth").String(),
		"url":  sect.Key("url").String(),
	}
}

// DatabaseCredentials will return a map of the database credentials
func DatabaseCredentials() map[string]string {
	if MustUseEnvVars() {
		return map[string]string{
			"user": os.Getenv("DUNGAR_DB_USER"),
			"pass": os.Getenv("DUNGAR_DB_PASS"),
			"host": os.Getenv("DUNGAR_DB_HOST"),
			"data": os.Getenv("DUNGAR_DB_DATA"),
		}
	}

	sect := secretsIniFile.Section("pgsql")

	return map[string]string{
		"user": sect.Key("user").String(),
		"pass": sect.Key("pass").String(),
		"host": sect.Key("host").String(),
		"data": sect.Key("data").String(),
	}
}
