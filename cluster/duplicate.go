package cluster

import "strings"

func DuplicateUID(uid, dashboards string) bool {
	// compare the UID with the dashboards
	return strings.Contains(dashboards, uid)
}
