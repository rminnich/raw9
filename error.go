package raw9

// size[4] Rerror tag[2] ename[s]

// assumption: we have a byte slice big enough to hold what
// we will produce.
func xrerror(b []byte, tag uint16, error string) int {
	s := 4 + 1 + 2 + 2 + len(error)
	b[0], b[1], b[2], b[3] = uint8(s), uint8(s>>8), uint8(s>>16), uint8(s>>24)
	b[4] = byte(Rerror)
	b[5], b[6] = uint8(tag), uint8(tag>>8)
	b[7], b[8] = uint8(len(error)), uint8(len(error)>>8)
	copy(b[9:], error)
	return s
}

// Assumption: something read the first 4 bytes of a message and returned
// a byte slice with that many bytes. For consistency, the message will include
// that length field.
func rrerror(b []byte) string {
	l := uint32(b[7]) | uint32(b[8])<<8
	e := string(b[9 : 9+l])
	return e
}
