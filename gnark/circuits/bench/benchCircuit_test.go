package bench

import (
	"testing"

	"github.com/consensys/gnark/test"
)

func TestBenchCircuit(t *testing.T) {

	assert := test.NewAssert(t)

	var benchCircuit BenchCircuit

	assert.ProverFailed(&benchCircuit, &BenchCircuit{
		X: 2,
		N: 12,
		Y: 4095,
	})

	assert.ProverSucceeded(&benchCircuit, &BenchCircuit{
		X: 2,
		N: 10,
		Y: 1024,
	})

}
