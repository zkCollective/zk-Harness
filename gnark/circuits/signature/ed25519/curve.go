package ed25519

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
)

// curve is the default twisted edwards companion curve (defined on api.Curve().Fr)
type curve struct {
	api    frontend.API
	fp     *fp25519
	params *CurveParams
}

func (c *curve) Params() *CurveParams {
	return c.params
}

func (c *curve) Add(p1, p2 Point) Point {
	return c.add(p1, p2)
}

func (c *curve) Double(p1 Point) Point {
	return c.double(&p1)
}
func (c *curve) Neg(p1 Point) Point {
	return c.neg(&p1)
}
func (c *curve) AssertIsOnCurve(p1 Point) {
	c.assertIsOnCurve(p1)
}

func (c *curve) AssertQ(p1 Point) {
	c.assertQ(&p1)
}

func (c *curve) ScalarMul(p1 Point, scalar frontend.Variable) Point {
	return c.scalarMul(c.api, &p1, scalar)
}
func (c *curve) DoubleBaseScalarMul(p1, p2 Point, s1 frontend.Variable, s2 *emulated.Element[Ed25519Fr]) Point {
	return c.doubleBaseScalarMul(c.api, &p1, &p2, s1, s2)
}
