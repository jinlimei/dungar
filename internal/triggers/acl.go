package triggers

var admins = []string{
	"U3Q9ZPR32",
}

func isAdmin(userID string) bool {
	if userID == "" {
		return false
	}

	//log.Printf("isAdmin(userID='%s') compare against admins='%v'",
	//	userID, admins)

	for _, adminID := range admins {
		if adminID == userID {
			return true
		}
	}

	return false
}
