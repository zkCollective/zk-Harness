package ed25519

import (
	"errors"
	"fmt"

	sha512 "github.com/zkCollective/zk-Harness/gnark/circuits/prf/sha512"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/logger"
	"github.com/consensys/gnark/std/math/emulated"
)

type Ed25519 struct {
	fp  *fp25519
	api frontend.API
}

// PublicKey stores an eddsa public key (to be used in gnark circuit)
type PublicKey struct {
	A Point
}

// Signature stores a signature  (to be used in gnark circuit)
// An EdDSA signature is a tuple (R,S) where R is a point on the twisted Edwards curve
// and S a scalar. Since the base field of the twisted Edwards is Fr, the number of points
// N on the Edwards is < r+1+2sqrt(r)+2 (since the curve has 2 points of multiplicity 2).
// The subgroup l used in eddsa is <1/2N, so the reduction
// mod l ensures S < r, therefore there is no risk of overflow.
type Signature struct {
	R Point
	S frontend.Variable
}

func NewEd25519(api frontend.API) (*Ed25519, error) {
	bf, err := emulated.NewField[Ed25519Fp](api)
	if err != nil {
		return nil, fmt.Errorf("new base api: %w", err)
	}
	sf, err := emulated.NewField[Ed25519Fr](api)
	if err != nil {
		return nil, fmt.Errorf("new scalar api: %w", err)
	}
	return &Ed25519{fp: newFp(bf, sf), api: api}, nil
}

func (e *Ed25519) bytes(v *Point) [32]frontend.Variable {
	outX := e.copyFieldElement(&v.X)
	outX0Bits := e.api.ToBinary(outX[0], 8)
	isNegative := outX0Bits[0]

	out := e.copyFieldElement(&v.Y)
	outY31Bits := e.api.ToBinary(out[31], 8)
	outY31Bits[7] = e.api.Or(outY31Bits[7], isNegative)
	out[31] = e.api.FromBinary(outY31Bits...)
	return out
}

func (e *Ed25519) copyFieldElement(v *emulated.Element[Ed25519Fp]) [32]frontend.Variable {
	bits := e.fp.baseField.ToBits(v)

	var buf [32]frontend.Variable
	for i := 0; i < 32; i++ {
		if i == 31 {
			buf[i] = e.api.FromBinary(bits[i*8:]...)
		} else {
			buf[i] = e.api.FromBinary(bits[i*8 : (i+1)*8]...)
		}
	}
	return buf
}

func (e *Ed25519) recomposeMsg(msg []frontend.Variable) []frontend.Variable {

	var msgBits []frontend.Variable
	for i := 4; i >= 0; i-- {
		if i == 4 {
			msgBits = append(msgBits, e.api.ToBinary(msg[i], 176)...)
		} else {
			msgI := e.api.ToBinary(msg[i], 200)
			msgBits = append(msgBits, msgI...)
		}
	}

	var buf [122]frontend.Variable
	for i := 0; i < 122; i++ {
		if i == 121 {
			buf[i] = e.api.FromBinary(msgBits[i*8:]...)
		} else {
			buf[i] = e.api.FromBinary(msgBits[i*8 : (i+1)*8]...)
		}
	}

	var result [122]frontend.Variable
	for i := 0; i < 122; i++ {
		result[i] = buf[len(result)-i-1]
	}

	return result[:]
}

func (e *Ed25519) recomposePbKey(vbs []frontend.Variable) []frontend.Variable {
	var bits []frontend.Variable

	for i := 0; i < 2; i++ {
		bits = append(bits, e.api.ToBinary(vbs[1-i], 128)...)
	}

	var buf [32]frontend.Variable
	for i := 0; i < 32; i++ {
		buf[i] = e.api.FromBinary(bits[i*8 : (i+1)*8]...)
	}

	var result [32]frontend.Variable
	for i := 0; i < 32; i++ {
		result[i] = buf[len(result)-i-1]
	}

	return result[:]
}

func (e *Ed25519) Verify(
	originPubs []frontend.Variable,
	sig Signature,
	msg [5]frontend.Variable,
	pubKey PublicKey) error {

	curve, err := newEdCurve(e.api)
	if err != nil {
		panic("new twistededwards curve failed")
	}

	hash := sha512.New(e.api)
	sigBytes := e.bytes(&sig.R)
	pubKeyBytes := e.bytes(&pubKey.A)

	recomposedPbs := e.recomposePbKey(originPubs)
	for i, b := range pubKeyBytes {
		e.api.AssertIsEqual(b, recomposedPbs[i])
	}

	hash.Write(sigBytes[:])
	hash.Write(pubKeyBytes[:])
	hash.Write(e.recomposeMsg(msg[:]))
	hRam := hash.Sum()

	var bits []frontend.Variable
	for i := 0; i < len(hRam); i++ {
		bits = append(bits, e.api.ToBinary(hRam[i], 8)...)
	}

	hRamScalar := e.fp.scalarField.FromBits(bits...)
	hRamScalar = e.fp.scalarField.Reduce(hRamScalar)

	baseX := emulated.ValueOf[Ed25519Fp](curve.Params().Base[0])
	baseY := emulated.ValueOf[Ed25519Fp](curve.Params().Base[1])

	base := Point{
		X: baseX,
		Y: baseY,
	}

	//[S]G-[H(R,A,M)]*A
	_A := curve.Neg(pubKey.A)

	Q := curve.DoubleBaseScalarMul(base, _A, sig.S, hRamScalar)
	curve.AssertIsOnCurve(Q)

	//[S]G-[H(R,A,M)]*A-R
	Q = curve.Add(curve.Neg(Q), sig.R)

	// [cofactor]*(lhs-rhs)
	log := logger.Logger()
	if !curve.Params().Cofactor.IsUint64() {
		err := errors.New("invalid cofactor")
		log.Err(err).Str("cofactor", curve.Params().Cofactor.String()).Send()
		return err
	}

	cofactor := curve.Params().Cofactor.Uint64()

	switch cofactor {
	case 4:
		Q = curve.Double(curve.Double(Q))
	case 8:
		Q = curve.Double(curve.Double(curve.Double(Q)))
	default:
		log.Warn().Str("cofactor", curve.Params().Cofactor.String()).Msg("curve cofactor is not implemented")
	}

	curve.AssertQ(Q)

	return nil
}
