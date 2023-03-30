package util

import (
	"github.com/consensys/gnark-crypto/ecc"
	bls12377fr "github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	bls24315fr "github.com/consensys/gnark-crypto/ecc/bls24-315/fr"
	bn254fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	bw6633fr "github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
	bw6761fr "github.com/consensys/gnark-crypto/ecc/bw6-761/fr"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/frontend"
)

func PreCalcMIMC(curveID ecc.ID, preImage frontend.Variable) interface{} {

	switch curveID {
	case ecc.BN254:
		// compute expected Y
		var expectedY bn254fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BN254.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh

	case ecc.BLS12_377:
		// compute expected Y
		var expectedY bls12377fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BLS12_377.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh

	case ecc.BLS24_315:
		// compute expected Y
		var expectedY bls24315fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BLS24_315.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh

	case ecc.BW6_761:
		// compute expected Y
		var expectedY bw6761fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BW6_761.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh

	case ecc.BW6_633:
		// compute expected Y
		var expectedY bw6633fr.Element
		expectedY.SetInterface(preImage)

		// running MiMC (Go)
		goMimc := hash.MIMC_BW6_633.New()
		goMimc.Write(expectedY.Marshal())
		expectedh := goMimc.Sum(nil)
		return expectedh
	default:
		panic("not implemented")
	}
}
