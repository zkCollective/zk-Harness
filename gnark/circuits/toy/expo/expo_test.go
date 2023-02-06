package expo

import (
	"testing"

	"github.com/consensys/gnark/test"
)

func TestExpoGroth16(t *testing.T) {

	assert := test.NewAssert(t)

	var benchCircuit BenchCircuit

	assert.ProverFailed(&benchCircuit, &BenchCircuit{
		X: 2,
		Y: 5,
		N: 1,
	})

	assert.ProverSucceeded(&benchCircuit, &BenchCircuit{
		X: 2,
		Y: 4,
		N: 1,
	})

}
