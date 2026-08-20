package anycast_control

func rememberIntent(nodeID string, in Intent) bool {
	state.Lock()
	defer state.Unlock()
	if old, exists := state.last[nodeID]; exists {
		if !isFencingAdvance(old, in) {
			return false
		}
		if isDuplicateIntent(old, in) {
			return false
		}
	}
	commitIntent(nodeID, in)
	return true
}
