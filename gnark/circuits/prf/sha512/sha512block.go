package sha256

import (
	"github.com/consensys/gnark/frontend"
)

var _K = []xuint64{
	constUint64(0x428a2f98d728ae22),
	constUint64(0x7137449123ef65cd),
	constUint64(0xb5c0fbcfec4d3b2f),
	constUint64(0xe9b5dba58189dbbc),
	constUint64(0x3956c25bf348b538),
	constUint64(0x59f111f1b605d019),
	constUint64(0x923f82a4af194f9b),
	constUint64(0xab1c5ed5da6d8118),
	constUint64(0xd807aa98a3030242),
	constUint64(0x12835b0145706fbe),
	constUint64(0x243185be4ee4b28c),
	constUint64(0x550c7dc3d5ffb4e2),
	constUint64(0x72be5d74f27b896f),
	constUint64(0x80deb1fe3b1696b1),
	constUint64(0x9bdc06a725c71235),
	constUint64(0xc19bf174cf692694),
	constUint64(0xe49b69c19ef14ad2),
	constUint64(0xefbe4786384f25e3),
	constUint64(0x0fc19dc68b8cd5b5),
	constUint64(0x240ca1cc77ac9c65),
	constUint64(0x2de92c6f592b0275),
	constUint64(0x4a7484aa6ea6e483),
	constUint64(0x5cb0a9dcbd41fbd4),
	constUint64(0x76f988da831153b5),
	constUint64(0x983e5152ee66dfab),
	constUint64(0xa831c66d2db43210),
	constUint64(0xb00327c898fb213f),
	constUint64(0xbf597fc7beef0ee4),
	constUint64(0xc6e00bf33da88fc2),
	constUint64(0xd5a79147930aa725),
	constUint64(0x06ca6351e003826f),
	constUint64(0x142929670a0e6e70),
	constUint64(0x27b70a8546d22ffc),
	constUint64(0x2e1b21385c26c926),
	constUint64(0x4d2c6dfc5ac42aed),
	constUint64(0x53380d139d95b3df),
	constUint64(0x650a73548baf63de),
	constUint64(0x766a0abb3c77b2a8),
	constUint64(0x81c2c92e47edaee6),
	constUint64(0x92722c851482353b),
	constUint64(0xa2bfe8a14cf10364),
	constUint64(0xa81a664bbc423001),
	constUint64(0xc24b8b70d0f89791),
	constUint64(0xc76c51a30654be30),
	constUint64(0xd192e819d6ef5218),
	constUint64(0xd69906245565a910),
	constUint64(0xf40e35855771202a),
	constUint64(0x106aa07032bbd1b8),
	constUint64(0x19a4c116b8d2d0c8),
	constUint64(0x1e376c085141ab53),
	constUint64(0x2748774cdf8eeb99),
	constUint64(0x34b0bcb5e19b48a8),
	constUint64(0x391c0cb3c5c95a63),
	constUint64(0x4ed8aa4ae3418acb),
	constUint64(0x5b9cca4f7763e373),
	constUint64(0x682e6ff3d6b2b8a3),
	constUint64(0x748f82ee5defb2fc),
	constUint64(0x78a5636f43172f60),
	constUint64(0x84c87814a1f0ab72),
	constUint64(0x8cc702081a6439ec),
	constUint64(0x90befffa23631e28),
	constUint64(0xa4506cebde82bde9),
	constUint64(0xbef9a3f7b2c67915),
	constUint64(0xc67178f2e372532b),
	constUint64(0xca273eceea26619c),
	constUint64(0xd186b8c721c0c207),
	constUint64(0xeada7dd6cde0eb1e),
	constUint64(0xf57d4f7fee6ed178),
	constUint64(0x06f067aa72176fba),
	constUint64(0x0a637dc5a2c898a6),
	constUint64(0x113f9804bef90dae),
	constUint64(0x1b710b35131c471b),
	constUint64(0x28db77f523047d84),
	constUint64(0x32caab7b40c72493),
	constUint64(0x3c9ebe0a15c9bebc),
	constUint64(0x431d67c49c100d4c),
	constUint64(0x4cc5d4becb3e42b6),
	constUint64(0x597f299cfc657e2a),
	constUint64(0x5fcb6fab3ad6faec),
	constUint64(0x6c44198c4a475817),
}

func blockGeneric(dig *digest, p []xuint8) {
	var w []xuint64

	var uapi = newUint64API(dig.api)
	for i := 0; i < chunk; i++ {
		w = append(w, uapi.asUint64(frontend.Variable(0)))
	}

	h0, h1, h2, h3, h4, h5, h6, h7 := dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7]
	for len(p) >= chunk {
		// Can interlace the computation of w with the
		// rounds below if needed for speed.
		for i := 0; i < 16; i++ {
			j := i * 8

			o1 := uapi.lshift(p[j].toUint64(), 56)

			o2 := uapi.lshift(p[j+1].toUint64(), 48)
			o3 := uapi.lshift(p[j+2].toUint64(), 40)
			o4 := uapi.lshift(p[j+3].toUint64(), 32)
			o5 := uapi.lshift(p[j+4].toUint64(), 24)
			o6 := uapi.lshift(p[j+5].toUint64(), 16)
			o7 := uapi.lshift(p[j+6].toUint64(), 8)
			o8 := p[j+7].toUint64()

			w[i] = uapi.or(o1, o2, o3, o4, o5, o6, o7, o8)
		}

		for i := 16; i < 80; i++ {
			v1 := w[i-2]
			t1 := uapi.xor(uapi.lrot(v1, -19), uapi.lrot(v1, -61), uapi.rshift(v1, 6))
			v2 := w[i-15]
			t2 := uapi.xor(uapi.lrot(v2, -1), uapi.lrot(v2, -8), uapi.rshift(v2, 7))

			w[i] = uapi.add(t1, w[i-7], t2, w[i-16])
		}

		a, b, c, d, e, f, g, h := h0, h1, h2, h3, h4, h5, h6, h7

		for i := 0; i < 80; i++ {
			t1 := uapi.add(
				h,
				uapi.xor(uapi.lrot(e, -14), uapi.lrot(e, -18), uapi.lrot(e, -41)),
				uapi.xor(uapi.and(e, f), uapi.and(uapi.not(e), g)),
				_K[i],
				w[i],
			)
			t2 := uapi.add(
				uapi.xor(uapi.lrot(a, -28), uapi.lrot(a, -34), uapi.lrot(a, -39)),
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
