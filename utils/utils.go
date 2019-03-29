package utils

//Joins two bytes together to form a 16 bit integer
func JoinBytes(hb, lb byte) uint16 {
	return (uint16(hb) << 8) | uint16(lb)
}

//Splits one 16 bit integer to two bytes
func SplitBytes(bytes uint16) (hb byte, lb byte) {
	hb = byte(bytes >> 8)
	lb = byte(bytes & 0x00FF)
	return
}

// BIT MANIPULATION HELPERS
func SetBit(b byte, pos byte) byte {
	return b | (1 << uint(pos))
}

func IsBitSet(b byte, pos byte) bool {
	return (b&(1<<pos) > 0)
}

func ClearBit(b byte, pos byte) byte {
	return b & ^(1 << pos)
}

func SwapNibbles(b byte) byte {
	return (b << 4) | (b >> 4)
}
