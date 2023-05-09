package sha256

import (
	"github.com/consensys/gnark/frontend"
)

var _K = []xuint32{
	constUint32(0x428a2f98),
	constUint32(0x71374491),
	constUint32(0xb5c0fbcf),
	constUint32(0xe9b5dba5),
	constUint32(0x3956c25b),
	constUint32(0x59f111f1),
	constUint32(0x923f82a4),
	constUint32(0xab1c5ed5),
	constUint32(0xd807aa98),
	constUint32(0x12835b01),
	constUint32(0x243185be),
	constUint32(0x550c7dc3),
	constUint32(0x72be5d74),
	constUint32(0x80deb1fe),
	constUint32(0x9bdc06a7),
	constUint32(0xc19bf174),
	constUint32(0xe49b69c1),
	constUint32(0xefbe4786),
	constUint32(0x0fc19dc6),
	constUint32(0x240ca1cc),
	constUint32(0x2de92c6f),
	constUint32(0x4a7484aa),
	constUint32(0x5cb0a9dc),
	constUint32(0x76f988da),
	constUint32(0x983e5152),
	constUint32(0xa831c66d),
	constUint32(0xb00327c8),
	constUint32(0xbf597fc7),
	constUint32(0xc6e00bf3),
	constUint32(0xd5a79147),
	constUint32(0x06ca6351),
	constUint32(0x14292967),
	constUint32(0x27b70a85),
	constUint32(0x2e1b2138),
	constUint32(0x4d2c6dfc),
	constUint32(0x53380d13),
	constUint32(0x650a7354),
	constUint32(0x766a0abb),
	constUint32(0x81c2c92e),
	constUint32(0x92722c85),
	constUint32(0xa2bfe8a1),
	constUint32(0xa81a664b),
	constUint32(0xc24b8b70),
	constUint32(0xc76c51a3),
	constUint32(0xd192e819),
	constUint32(0xd6990624),
	constUint32(0xf40e3585),
	constUint32(0x106aa070),
	constUint32(0x19a4c116),
	constUint32(0x1e376c08),
	constUint32(0x2748774c),
	constUint32(0x34b0bcb5),
	constUint32(0x391c0cb3),
	constUint32(0x4ed8aa4a),
	constUint32(0x5b9cca4f),
	constUint32(0x682e6ff3),
	constUint32(0x748f82ee),
	constUint32(0x78a5636f),
	constUint32(0x84c87814),
	constUint32(0x8cc70208),
	constUint32(0x90befffa),
	constUint32(0xa4506ceb),
	constUint32(0xbef9a3f7),
	constUint32(0xc67178f2),
}

func blockGeneric(dig *digest, p []xuint8) {
	var w []xuint32

	var uapi = newUint32API(dig.api)
	for i := 0; i < 64; i++ {
		w = append(w, uapi.asUint32(frontend.Variable(0)))
	}

	h0, h1, h2, h3, h4, h5, h6, h7 := dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7]
	for len(p) >= chunk {
		// Can interlace the computation of w with the
		// rounds below if needed for speed.
		for i := 0; i < 16; i++ {
			j := i * 4

			o1 := uapi.lshift(p[j].toUint32(), 24)

			o2 := uapi.lshift(p[j+1].toUint32(), 16)
			o3 := uapi.lshift(p[j+2].toUint32(), 8)
			o4 := p[j+3].toUint32()

			w[i] = uapi.or(o1, o2, o3, o4)
		}

		for i := 16; i < 64; i++ {
			v1 := w[i-2]
			t1 := uapi.xor(uapi.lrot(v1, -17), uapi.lrot(v1, -19), uapi.rshift(v1, 10))
			v2 := w[i-15]
			t2 := uapi.xor(uapi.lrot(v2, -7), uapi.lrot(v2, -18), uapi.rshift(v2, 3))

			w[i] = uapi.add(t1, w[i-7], t2, w[i-16])
		}

		a, b, c, d, e, f, g, h := h0, h1, h2, h3, h4, h5, h6, h7

		for i := 0; i < 64; i++ {
			t1 := uapi.add(
				h,
				uapi.xor(uapi.lrot(e, -6), uapi.lrot(e, -11), uapi.lrot(e, -25)),
				uapi.xor(uapi.and(e, f), uapi.and(uapi.not(e), g)),
				_K[i],
				w[i],
			)
			t2 := uapi.add(
				uapi.xor(uapi.lrot(a, -2), uapi.lrot(a, -13), uapi.lrot(a, -22)),
				uapi.xor(uapi.and(a, b), uapi.and(a, c), uapi.and(b, c)),
			)

			h = g
			g = f
			f = e
			e = uapi.add(d, t1)
			d = c
			c = b
			b = a
			a = uapi.add(t1, t2)
		}

		h0 = uapi.add(h0, a)
		h1 = uapi.add(h1, b)
		h2 = uapi.add(h2, c)
		h3 = uapi.add(h3, d)
		h4 = uapi.add(h4, e)
		h5 = uapi.add(h5, f)
		h6 = uapi.add(h6, g)
		h7 = uapi.add(h7, h)

		p = p[chunk:]
	}

	dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7] = h0, h1, h2, h3, h4, h5, h6, h7
}
