package ed25519

import (
	"math/big"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
)

// Curve methods implemented by a twisted edwards curve inside a circuit
type Curve interface {
	Params() *CurveParams
	Add(p1, p2 Point) Point
	Double(p1 Point) Point
	Neg(p1 Point) Point
	AssertIsOnCurve(p1 Point)
	AssertQ(p1 Point)
	ScalarMul(p1 Point, scalar frontend.Variable) Point
	DoubleBaseScalarMul(p1, p2 Point, s1 frontend.Variable, s2 *emulated.Element[Ed25519Fr]) Point
}

// Point represent a pair of X, Y coordinates inside a circuit
type Point struct {
	X, Y emulated.Element[Ed25519Fp]
}

// OPoint initial data present as []byte
type OPoint struct {
	X, Y []frontend.Variable
}

// CurveParams twisted edwards curve parameters ax^2 + y^2 = 1 + d*x^2*y^2
type CurveParams struct {
	A, D, Cofactor, Order *big.Int
	Base                  [2]*big.Int // base point coordinates
}

// EndoParams endomorphism parameters for the curve, if they exist
type EndoParams struct {
	Endo   [2]*big.Int
	Lambda *big.Int
}

func newEdCurve(api frontend.API) (Curve, error) {
	var curveParameter = &CurveParams{
		A:        new(big.Int),
		D:        new(big.Int),
		Cofactor: new(big.Int),
		Order:    new(big.Int),
		Base:     [2]*big.Int{new(big.Int), new(big.Int)},
	}

	curveParameter.A.SetString("-1", 10)
	curveParameter.D.SetString("37095705934669439343138083508754565189542113879843219016388785533085940283555", 10)
	curveParameter.Cofactor.SetString("8", 10)
	curveParameter.Order.SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)
	curveParameter.Base[0].SetString("15112221349535400772501151409588531511454012693041857206046113283949847762202", 10)
	curveParameter.Base[1].SetString("46316835694926478169428394003475163141307993866256225615783033603165251855960", 10)

	bf, err := emulated.NewField[Ed25519Fp](api)
	if err != nil {
		return nil, err
	}
	sf, err := emulated.NewField[Ed25519Fr](api)
	if err != nil {
		return nil, err
	}
	return &curve{api: api, fp: &fp25519{baseField: bf, scalarField: sf}, params: curveParameter}, nil
}
