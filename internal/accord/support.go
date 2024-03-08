package accord

import (
	"fmt"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// PingUser handles the format for pinging users
func (d *Driver) PingUser(user core2.User) string {
	return fmt.Sprintf("<@%s> ", user.ID)
}
