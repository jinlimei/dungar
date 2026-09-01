package accord

import "strings"

// GetRoleName attempts to retrieve a role name from the discord API or from the local cache.
// this will return an empty string if it failed for some reason.
func (d *Driver) GetRoleName(roleID, serverID string) string {
	guild := d.getOrMakeGuild(serverID)

	role, ok := guild.roleCache[roleID]

	if !ok {
		return ""
	}

	return role.Name
}

// HasRole checks whether a user currently has a named role in the given guild.
func (d *Driver) HasRole(userID, serverID, roleName string) bool {
	if userID == "" || serverID == "" || strings.TrimSpace(roleName) == "" {
		return false
	}

	guild := d.getOrMakeGuild(serverID)
	member, ok := guild.memberCache[userID]

	if !ok {
		retrieved, err := d.Con.GetSession().GuildMember(serverID, userID)
		if err != nil {
			return false
		}

		member = retrieved
		guild.memberCache[userID] = member
	}

	if member == nil || len(member.Roles) == 0 {
		return false
	}

	roleName = strings.ToLower(strings.TrimSpace(roleName))

	for _, roleID := range member.Roles {
		role, ok := guild.roleCache[roleID]
		if ok && strings.ToLower(role.Name) == roleName {
			return true
		}
	}

	return false
}
