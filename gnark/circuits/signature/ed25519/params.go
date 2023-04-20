package ed25519

import "math/big"

var (
	qCurve25519, rCurve25519 *big.Int
)

func init() {
	qCurve25519, _ = new(big.Int).SetString("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed", 16)
	rCurve25519, _ = new(big.Int).SetString("1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed", 16)
}

type Ed25519Fp struct{}

func (fp Ed25519Fp) NbLimbs() uint     { return 5 }
func (fp Ed25519Fp) BitsPerLimb() uint { return 51 }
func (fp Ed25519Fp) IsPrime() bool     { return true }
func (fp Ed25519Fp) Modulus() *big.Int { return qCurve25519 }

type Ed25519Fr struct{}

func (fp Ed25519Fr) NbLimbs() uint     { return 5 }
func (fp Ed25519Fr) BitsPerLimb() uint { return 51 }
func (fp Ed25519Fr) IsPrime() bool     { return true }
func (fp Ed25519Fr) Modulus() *big.Int { return rCurve25519 }
