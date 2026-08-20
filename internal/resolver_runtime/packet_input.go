package resolver_runtime

func cloneInputPacket(packet []byte) []byte {
	if packet == nil {
		return nil
	}
	owned := make([]byte, len(packet))
	copy(owned, packet)
	return owned
}
