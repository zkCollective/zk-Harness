package exponentiate

import (
	"testing"

	"github.com/consensys/gnark/test"
)

func TestExponentiateCircuit(t *testing.T) {

	assert := test.NewAssert(t)

	var exponentiateCircuit ExponentiateCircuit

	assert.ProverFailed(&exponentiateCircuit, &ExponentiateCircuit{
		X: 2,
		E: 12,
		Y: 4095,
	})

	assert.ProverSucceeded(&exponentiateCircuit, &ExponentiateCircuit{
		X: 1,
		E: 10000,
		Y: 1,
	})

}
