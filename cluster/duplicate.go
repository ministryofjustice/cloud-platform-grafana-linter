package cluster

func DuplicateUID(uid, dashboards string) bool {
	return uid == dashboards
}
