package sha256

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
)

const chunk = 64

var (
	init0 = constUint32(0x6A09E667)
	init1 = constUint32(0xBB67AE85)
	init2 = constUint32(0x3C6EF372)
	init3 = constUint32(0xA54FF53A)
	init4 = constUint32(0x510E527F)
	init5 = constUint32(0x9B05688C)
	init6 = constUint32(0x1F83D9AB)
	init7 = constUint32(0x5BE0CD19)
)

type digest struct {
	h   [8]xuint32
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
	var tmp [64]xuint8
	tmp[0] = constUint8(0x80)
	for i := 1; i < 64; i++ {
		tmp[i] = constUint8(0x0)
	}
	if len%64 < 56 {
		d.write(tmp[0 : 56-len%64])
	} else {
		d.write(tmp[0 : 64+56-len%64])
	}

	// fill length bit
	len <<= 3
	PutUint64(d.api, tmp[:], newUint64API(d.api).asUint64(len))
	d.write(tmp[0:8])
	fmt.Printf("block number:%d\n", d.len/64)

	if d.nx != 0 {
		panic("d.nx != 0")
	}

	var digest [32]xuint8

	// h[0]..h[7]
	PutUint32(d.api, digest[0:], d.h[0])
	PutUint32(d.api, digest[4:], d.h[1])
	PutUint32(d.api, digest[8:], d.h[2])
	PutUint32(d.api, digest[12:], d.h[3])
	PutUint32(d.api, digest[16:], d.h[4])
	PutUint32(d.api, digest[20:], d.h[5])
	PutUint32(d.api, digest[24:], d.h[6])
	PutUint32(d.api, digest[28:], d.h[7])

	var dv []frontend.Variable

	u8api := newUint8API(d.api)

	for i := 0; i < 32; i++ {
		dv = append(dv, u8api.fromUint8(digest[i]))
	}
	return dv
}
