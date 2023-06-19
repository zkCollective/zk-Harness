package sha256

import "github.com/consensys/gnark/frontend"

func PutUint64(api frontend.API, b []xuint8, v xuint64) {
	_ = b[7] // early bounds check to guarantee safety of writes below
	u64api := newUint64API(api)
	//b[0] = byte(v >> 56)
	b[0] = u64api.rshift(v, 56).toxUnit8()
	//b[1] = byte(v >> 48)
	b[1] = u64api.rshift(v, 48).toxUnit8()
	//b[2] = byte(v >> 40)
	b[2] = u64api.rshift(v, 40).toxUnit8()
	//b[3] = byte(v >> 32)
	b[3] = u64api.rshift(v, 32).toxUnit8()
	//b[4] = byte(v >> 24)
	b[4] = u64api.rshift(v, 24).toxUnit8()
	//b[5] = byte(v >> 16)
	b[5] = u64api.rshift(v, 16).toxUnit8()
	//b[6] = byte(v >> 8)
	b[6] = u64api.rshift(v, 8).toxUnit8()
	//b[7] = byte(v)
	b[7] = v.toxUnit8()
}

func PutUint32(api frontend.API, b []xuint8, v xuint32) {
	_ = b[3] // early bounds check to guarantee safety of writes below
	uint32api := newUint32API(api)
	// b[0] = byte(v >> 24)
	b[0] = uint32api.rshift(v, 24).toUnit8()
	// b[1] = byte(v >> 16)
	b[1] = uint32api.rshift(v, 16).toUnit8()
	// b[2] = byte(v >> 8)
	b[2] = uint32api.rshift(v, 8).toUnit8()
	// b[3] = byte(v)
	b[3] = v.toUnit8()
}
