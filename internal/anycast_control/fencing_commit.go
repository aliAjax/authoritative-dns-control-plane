package anycast_control

func commitIntent(nodeID string, in Intent) {
	state.last[nodeID] = in
}
