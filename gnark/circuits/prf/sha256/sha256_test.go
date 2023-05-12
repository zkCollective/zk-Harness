/*
Copyright © 2023 Jan Lauinger
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sha256

import (
	"encoding/hex"
	"testing"

	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
)

func TestSha256All(t *testing.T) {
	assert := test.NewAssert(t)

	input := "68656c6c6f20776f726c64"
	output := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

	// 'hello-world-hello-world-hello-world-hello-world-hello-world-12345' as hex
	// input := "68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d68656c6c6f2d776f726c642d3132333435"
	// output := "34caf9dcd6b137c56c59f81e071a4b77a11329f26c80d7023ac7dfc485dcd780"

	byteSlice, _ := hex.DecodeString(input)
	inputByteLen := len(byteSlice)

	byteSlice, _ = hex.DecodeString(output)
	outputByteLen := len(byteSlice)

	// witness definition
	preImageAssign := StrToIntSlice(input, true)
	outputAssign := StrToIntSlice(output, true)

	// witness values preparation
	//assignment := Sha256Circuit{
	//	PreImage:       make([]frontend.Variable, inputByteLen),
	//	ExpectedResult: [32]frontend.Variable{},
	//}

	assignment := Sha256Circuit{
		In:             make([]frontend.Variable, inputByteLen),
		ExpectedResult: [32]frontend.Variable{},
	}

	// assign values here because required to use make in assignment
	for i := 0; i < inputByteLen; i++ {
		assignment.In[i] = preImageAssign[i]
	}
	for i := 0; i < outputByteLen; i++ {
		assignment.ExpectedResult[i] = outputAssign[i]
	}

	circuit := Sha256Circuit{
		In: make([]frontend.Variable, inputByteLen),
	}

	// Currently, this version of SHA256 only works with groth16
	assert.SolvingSucceeded(&circuit, &assignment, test.WithBackends(backend.GROTH16))
}
