package exponentiate

import "github.com/consensys/gnark/frontend"

// simple circuit that checks X*X*X*X*X... == Y
type ExponentiateCircuit struct {
	X frontend.Variable
	Y frontend.Variable `gnark:",public"`
	E int
}

// Circuit defines a an exponentiation for a frontend variable with itself
func (circuit *ExponentiateCircuit) Define(api frontend.API) error {
	for i := 0; i < circuit.E; i++ {
		circuit.X = api.Mul(circuit.X, circuit.X)
	}
	api.AssertIsEqual(circuit.Y, circuit.X)
	return nil
}
