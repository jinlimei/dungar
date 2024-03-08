package db

import (
	"log"
)

var writeIssues = true

// FlagWritingIssues changes whether or not to write issues,
// and is built to disable writing issues in a test environment
func FlagWritingIssues(enabled bool) {
	writeIssues = enabled
}

// LogIssue logs an issue to the database for easier handling
// of it in the future.
func LogIssue(issueType, title, contents string) {
	if !writeIssues {
		return
	}

	log.Printf("[ISSUE ENCOUNTERED %s] %s, %s\n",
		issueType, title, contents)

	//utils.SentryMessage(fmt.Sprintf(
	//	"[LogIssue] Type='%s', Title='%s', Contents='%s'",
	//	issueType,
	//	title,
	//	contents,
	//), nil)

	sql := `
		INSERT INTO log_issues (issue_type, title, contents, created)
		VALUES($1, $2, $3, CURRENT_TIMESTAMP)
	`

	_, err := ConExec(sql, issueType, title, contents)

	if err != nil {
		log.Printf("[ISSUE EXPANDED] Also Failed to execute query: %v\n", err)
	}
}
