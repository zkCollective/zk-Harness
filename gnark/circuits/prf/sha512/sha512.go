package sha256

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
)

const chunk = 128

var (
	init0 = constUint64(0x6a09e667f3bcc908)
	init1 = constUint64(0xbb67ae8584caa73b)
	init2 = constUint64(0x3c6ef372fe94f82b)
	init3 = constUint64(0xa54ff53a5f1d36f1)
	init4 = constUint64(0x510e527fade682d1)
	init5 = constUint64(0x9b05688c2b3e6c1f)
	init6 = constUint64(0x1f83d9abfb41bd6b)
	init7 = constUint64(0x5be0cd19137e2179)
)

type digest struct {
	h   [8]xuint64
	x   [chunk]xuint8 // 64 byte
	nx  int
	len uint64
	id  ecc.ID
	api frontend.API
}

func (d *digest) Reset() {
	d.h[0] = init0
	d.h[1] = init1
	d.h[2] = init2
	d.h[3] = init3
	d.h[4] = init4
	d.h[5] = init5
	d.h[6] = init6
	d.h[7] = init7

	d.nx = 0
	d.len = 0
}

func New(api frontend.API) digest {
	res := digest{}
	res.id = ecc.BN254
	res.api = api
	res.nx = 0
	res.len = 0
	res.Reset()
	return res
}

// p: byte array
func (d *digest) Write(p []frontend.Variable) (nn int, err error) {

	var in []xuint8
	for i := range p {
		in = append(in, newUint8API(d.api).asUint8(p[i]))
	}
	return d.write(in)

}

func (d *digest) write(p []xuint8) (nn int, err error) {
	nn = len(p)
	d.len += uint64(nn)

	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == chunk {
			blockGeneric(d, d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}

	if len(p) >= chunk {
		n := len(p) &^ (chunk - 1)
		blockGeneric(d, p[:n])
		p = p[n:]
	}

	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}

	return
}

func (d *digest) Sum() []frontend.Variable {

	d0 := *d
	hash := d0.checkSum()

	return hash[:]
}

func (d *digest) checkSum() []frontend.Variable {
	// Padding
	len := d.len
	var tmp [128]xuint8
	tmp[0] = constUint8(0x80)
	for i := 1; i < 128; i++ {
		tmp[i] = constUint8(0x0)
	}
	if len%128 < 112 {
		d.write(tmp[0 : 112-len%128])
	} else {
		d.write(tmp[0 : 128+112-len%128])
	}

	// fill length bit
	len <<= 3
	PutUint64(d.api, tmp[0:], newUint64API(d.api).asUint64(0))
	PutUint64(d.api, tmp[8:], newUint64API(d.api).asUint64(len))
	d.write(tmp[0:16])
	//fmt.Printf("block number:%d\n", d.len/64)

	if d.nx != 0 {
		panic("d.nx != 0")
	}

	var digest [64]xuint8

	// h[0]..h[7]
	PutUint64(d.api, digest[0:], d.h[0])
	PutUint64(d.api, digest[8:], d.h[1])
	PutUint64(d.api, digest[16:], d.h[2])
	PutUint64(d.api, digest[24:], d.h[3])
	PutUint64(d.api, digest[32:], d.h[4])
	PutUint64(d.api, digest[40:], d.h[5])
	PutUint64(d.api, digest[48:], d.h[6])
	PutUint64(d.api, digest[56:], d.h[7])

	var dv []frontend.Variable

	u8api := newUint8API(d.api)

	for i := 0; i < 64; i++ {
		dv = append(dv, u8api.fromUint8(digest[i]))
	}
	return dv
}
