package taut

import (
	"fmt"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// PingUser will take the incoming core2.User and provide a valid ping syntax for a message
func (d *Driver) PingUser(user core2.User) string {
	return fmt.Sprintf("<@%s>", user.ID)
}
