package utils

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

const (
	// EnvCI is our name for the environmental variable IN_CI_ENV
	EnvCI = "IN_CI_ENV"
	// EnvCD is our name for the environmental variable IN_CD_ENV
	EnvCD = "IN_CD_ENV"
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
