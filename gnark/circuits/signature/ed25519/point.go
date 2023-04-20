package ed25519

import (
	edwards25519 "github.com/zkCollective/zk-Harness/gnark/circuits/signature/ed25519/test"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
)

type curveF = emulated.Field[Ed25519Fp]
type curveFr = emulated.Field[Ed25519Fr]

type fp25519 struct {
	baseField   *curveF
	scalarField *curveFr
}

func newFp(baseField *curveF, scalarField *curveFr) *fp25519 {
	return &fp25519{baseField, scalarField}
}

// PutBigEndian p is a little-endian bytes
func PutBigEndian(p []byte) []byte {
	var result []byte

	for i := 0; i < len(p); i++ {
		result = append(result, p[len(p)-1-i])
	}

	return result
}

func NewEmulatedPoint(v *edwards25519.Point) *Point {
	affine := (&edwards25519.PointAffine{}).FromExtended(v)

	xbytes := affine.X.Bytes()
	for i := 0; i < len(xbytes)/2; i++ {
		xbytes[i], xbytes[len(xbytes)-i-1] = xbytes[len(xbytes)-i-1], xbytes[i]
	}

	ybytes := affine.Y.Bytes()
	for i := 0; i < len(ybytes)/2; i++ {
		ybytes[i], ybytes[len(ybytes)-i-1] = ybytes[len(ybytes)-i-1], ybytes[i]
	}

	return &Point{
		X: emulated.ValueOf[Ed25519Fp](xbytes),
		Y: emulated.ValueOf[Ed25519Fp](ybytes),
	}
}

func (c *curve) add(P1, P2 Point) Point {
	api := c.fp.baseField

	a := emulated.ValueOf[Ed25519Fp](c.Params().A)
	d := emulated.ValueOf[Ed25519Fp](c.Params().D)
	u1 := api.Mul(&P1.X, &a)
	u1 = api.Sub(&P1.Y, u1)
	u2 := api.Add(&P2.X, &P2.Y)
	u := api.Mul(u1, u2)

	// v0 = x1 * y2
	v0 := api.Mul(&P2.Y, &P1.X)

	// v1 = x2 * y1
	v1 := api.Mul(&P2.X, &P1.Y)

	// v2 = d * v0 * v1
	v0d := api.Mul(&d, v0)
	v2 := api.Mul(v0d, v1)

	// x = (v0 + v1) / (1 + v2)
	px := api.Add(v0, v1)
	px = api.Div(px, api.Add(api.One(), v2))

	// y = (u + a * v0 - v1) / (1 - v2)
	py := api.Mul(&a, v0)
	py = api.Sub(py, v1)
	py = api.Add(py, u)
	py = api.Div(py, api.Sub(api.One(), v2))

	return Point{
		X: *px,
		Y: *py,
	}
}

func (c *curve) neg(p1 *Point) Point {
	x := c.fp.baseField.Neg(&p1.X)

	return Point{
		X: *x,
		Y: p1.Y,
	}
}

func (c *curve) double(p1 *Point) Point {

	api := c.fp.baseField

	u := api.Mul(&p1.X, &p1.Y)
	v := api.Mul(&p1.X, &p1.X)
	w := api.Mul(&p1.Y, &p1.Y)

	two := emulated.ValueOf[Ed25519Fp](2)
	n1 := api.Mul(&two, u)

	a := emulated.ValueOf[Ed25519Fp](c.params.A)
	av := api.Mul(v, &a)
	n2 := api.Sub(w, av)
	d1 := api.Add(w, av)
	d2 := api.Sub(&two, d1)

	pX := api.Div(n1, d1)
	pY := api.Div(n2, d2)

	return Point{
		X: *pX,
		Y: *pY,
	}
}

func (c *curve) doubleBaseScalarMul(api frontend.API, p1, p2 *Point, s1 frontend.Variable, s2 *emulated.Element[Ed25519Fr]) Point {
	baseField := c.fp.baseField
	// first unpack the scalars

	b1 := api.ToBinary(s1)
	b2 := c.fp.scalarField.ToBits(s2)

	res := Point{}
	tmp := Point{}
	sum := Point{}

	sum = c.Add(*p1, *p2)

	n := len(b1)

	res.X = *baseField.Lookup2(b1[n-1], b2[n-1], baseField.Zero(), &p1.X, &p2.X, &sum.X)
	res.Y = *baseField.Lookup2(b1[n-1], b2[n-1], baseField.One(), &p1.Y, &p2.Y, &sum.Y)

	for i := n - 2; i >= 0; i-- {
		res = c.Double(res)
		tmp.X = *baseField.Lookup2(b1[i], b2[i], baseField.Zero(), &p1.X, &p2.X, &sum.X)
		tmp.Y = *baseField.Lookup2(b1[i], b2[i], baseField.One(), &p1.Y, &p2.Y, &sum.Y)
		res = c.Add(res, tmp)
	}

	return Point{
		X: res.X,
		Y: res.Y,
	}
}

func (c *curve) scalarMul(api frontend.API, p1 *Point, scalar frontend.Variable) Point {
	baseF := c.fp.baseField
	// first unpack the scalar
	b := api.ToBinary(scalar)

	res := Point{}
	tmp := Point{}
	B := Point{}

	A := c.Double(*p1)
	B = c.Add(A, *p1)

	n := len(b) - 1
	res.X = *baseF.Lookup2(b[n], b[n-1], baseF.Zero(), &A.X, &p1.X, &B.X)
	res.Y = *baseF.Lookup2(b[n], b[n-1], baseF.One(), &A.Y, &p1.Y, &B.Y)

	for i := n - 2; i >= 1; i -= 2 {

		res = c.Double(res)
		res = c.Double(res)

		tmp.X = *baseF.Lookup2(b[i], b[i-1], baseF.Zero(), &A.X, &p1.X, &B.X)
		tmp.Y = *baseF.Lookup2(b[i], b[i-1], baseF.One(), &A.Y, &p1.Y, &B.Y)

		res = c.Add(res, tmp)
	}

	if n%2 == 0 {
		res = c.Double(res)
		tmp = c.Add(res, *p1)
		res.X = *baseF.Select(b[0], &tmp.X, &res.X)
		res.Y = *baseF.Select(b[0], &tmp.Y, &res.Y)
	}

	return Point{
		X: res.X,
		Y: res.Y,
	}
}

func (c *curve) assertIsOnCurve(p Point) {
	api := c.fp.baseField
	xx := api.Mul(&p.X, &p.X)
	yy := api.Mul(&p.Y, &p.Y)
	a := emulated.ValueOf[Ed25519Fp](c.params.A)
	axx := api.Mul(xx, &a)
	lhs := api.Add(axx, yy)

	d := emulated.ValueOf[Ed25519Fp](c.params.D)
	dxx := api.Mul(xx, &d)
	dxxyy := api.Mul(dxx, yy)
	rhs := api.Add(dxxyy, api.One())

	api.AssertIsEqual(lhs, rhs)
}

func (c *curve) assertQ(p *Point) {
	c.fp.baseField.AssertIsEqual(&p.X, c.fp.baseField.Zero())
	c.fp.baseField.AssertIsEqual(&p.Y, c.fp.baseField.One())
}
