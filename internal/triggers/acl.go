package triggers

import "slices"

var admins = []string{
	"U3Q9ZPR32",
}

func isAdmin(userID string) bool {
	if userID == "" {
		return false
	}

	//log.Printf("isAdmin(userID='%s') compare against admins='%v'",
	//	userID, admins)

	return slices.Contains(admins, userID)
}
