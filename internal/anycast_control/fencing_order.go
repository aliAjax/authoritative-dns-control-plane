package anycast_control

func isFencingAdvance(old, next Intent) bool {
	return next.Fencing > old.Fencing
}
