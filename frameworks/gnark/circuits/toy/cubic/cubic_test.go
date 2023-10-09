// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cubic

import (
	"testing"

	"github.com/consensys/gnark/test"
)

func TestCubicEquation(t *testing.T) {
	assert := test.NewAssert(t)

	var cubicCircuit CubicCircuit

	// ProverFailed fails the test if any of the following step errored:
	//
	// 1. compiles the circuit (or fetch it from the cache)
	// 2. using the test execution engine, executes the circuit with provided witness (must fail)
	// 3. run Setup / Prove / Verify with the backend (must fail)
	//
	// By default, this tests on all curves and proving schemes supported by gnark. See available TestingOption.

	assert.ProverFailed(&cubicCircuit, &CubicCircuit{
		X: 42,
		Y: 41,
	})

	// ProverSucceeded fails the test if any of the following step errored:
	//
	// 1. compiles the circuit (or fetch it from the cache)
	// 2. using the test execution engine, executes the circuit with provided witness
	// 3. run Setup / Prove / Verify with the backend
	// 4. if set, (de)serializes the witness and call ReadAndProve and ReadAndVerify on the backend
	//
	// By default, this tests on all curves and proving schemes supported by gnark. See available TestingOption.

	assert.ProverSucceeded(&cubicCircuit, &CubicCircuit{
		X: 3,
		Y: 35,
	})

}
