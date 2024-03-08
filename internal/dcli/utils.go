package dcli

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// DBConnect connects to the DB for a CLI command
func DBConnect() {
	if db.GetDatabase() != nil {
		return
	}

	utils.LoadSettingsAndSecrets()
	db.ConnectToDatabase()
}
