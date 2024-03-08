package accord

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
