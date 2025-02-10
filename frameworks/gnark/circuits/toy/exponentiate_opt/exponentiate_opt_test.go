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

package exponentiate_opt

import (
	"testing"

	"github.com/consensys/gnark/test"
)

func TestExponentiateGroth16(t *testing.T) {

	assert := test.NewAssert(t)

	var expCircuit ExponentiateOptCircuit

	assert.ProverFailed(&expCircuit, &ExponentiateOptCircuit{
		X: 2,
		E: 12,
		Y: 4095,
	})

	assert.ProverSucceeded(&expCircuit, &ExponentiateOptCircuit{
		X: 2,
		E: 12,
		Y: 4096,
	})

}
