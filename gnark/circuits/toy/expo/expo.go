package expo

import "github.com/consensys/gnark/frontend"

// benchCircuit is a simple circuit that checks X*X*X*X*X... == Y
type BenchCircuit struct {
	X frontend.Variable
	Y frontend.Variable `gnark:",public"`
	N int
}

func (circuit *BenchCircuit) Define(api frontend.API) error {
	for i := 0; i < circuit.N; i++ {
		circuit.X = api.Mul(circuit.X, circuit.X)
	}
	api.AssertIsEqual(circuit.X, circuit.Y)
	return nil
}
