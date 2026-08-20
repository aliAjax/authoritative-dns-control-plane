package anycast_control

func isDuplicateIntent(old, next Intent) bool {
	return old.Fencing == next.Fencing
}
