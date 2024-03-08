package db

import "gitlab.int.magneato.site/dungar/prototype/internal/utils"

// TestDatabaseConnect handles a database connection in a test environment
func TestDatabaseConnect() {
	//if utils.InTestEnv() {
	//	return
	//}

	//if GetDatabase() != nil {
	//	return
	//}

	utils.LoadSettingsAndSecrets()
	ConnectToDatabase()
}

// TruncateUserTracking truncates the user tracking table
func TruncateUserTracking() {
	//TestDatabaseConnect()
	//ConMustExec("TRUNCATE TABLE user_tracking")
}

// SeedUserTracking is a test function for seeding user tracking (either in here or in triggers)
func SeedUserTracking() {
	//TestDatabaseConnect()
	//ConMustExec(`
	//	INSERT INTO user_tracking (unique_id, nick, line_count, first_seen, last_seen)
	//	VALUES
	//  ('outdated', 'outdated', 100, '2019-06-01 12:00:00', '2019-07-01 12:00:00'),
	//	('updated',  'updated',   10, '2017-01-01 00:00:00', CURRENT_TIMESTAMP),
	//  ('ddddddd',  'ddddddd',  999, '2017-01-01 00:00:00', CURRENT_TIMESTAMP)
	//`)
}
