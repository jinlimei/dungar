package core2

import "errors"

var (
	// ErrInvalidUserID is when a user id is invalid
	ErrInvalidUserID = errors.New("invalid user id")
	// ErrUserNotFound is when a user was unable to be found
	ErrUserNotFound = errors.New("could not find user with provided id")
)

// User is a struct for super basic information about a User
type User struct {
	ID       string
	ServerID string
	Name     string
	IsBot    bool
	IsAdmin  bool
}

// BotUser is the struct for the bot
type BotUser User

// GetUser will retrieve a User via their unique id
func (s *Service) GetUser(userID, serverID string) (User, error) {
	return s.driver.GetUser(userID, serverID)
}

// GetUsers will retrieve a list of users in a map by their unique ID
func (s *Service) GetUsers(serverID string) map[string]User {
	return s.driver.GetUsers(serverID)
}

// GetUserName will retrieve just whatever is considered their username by a given service
func (s *Service) GetUserName(userID, serverID string) string {
	return s.driver.GetUserName(userID, serverID)
}

// GetBotUser will retrieve the BotUser
func (s *Service) GetBotUser() BotUser {
	return s.driver.GetBotUser()
}
