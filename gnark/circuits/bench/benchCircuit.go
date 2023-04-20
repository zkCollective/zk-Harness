package bench

import "github.com/consensys/gnark/frontend"

// benchCircuit is a simple circuit that checks X*X*X*X*X... == Y
type BenchCircuit struct {
	X frontend.Variable
	Y frontend.Variable `gnark:",public"`
	N int
}

// Circuit defines a an exponentiation for a frontend variable with itself
func (circuit *BenchCircuit) Define(api frontend.API) error {
	for i := 0; i < circuit.N; i++ {
		circuit.X = api.Mul(circuit.X, circuit.X)
	}
	api.AssertIsEqual(circuit.Y, circuit.X)
	return nil
}
