SHELL = zsh

circom = circuits/toy/multiply.circom
directory = ./circuits/toy/tmp
r1cs = $(directory)/multiply.r1cs
wasm = $(directory)/multiply_js/multiply.wasm
wit_gen = $(directory)/multiply_js/generate_witness.js
compile_outputs = $(directory)/multiply_js/witness_calculator.js $(r1cs) $(wasm) $(wit_gen)

$(compile_outputs): $(circom)
	mkdir $(directory)
	circom $< --r1cs --wasm -o $(directory)

clean:
	rm -r $(compile_outputs)
	rmdir $(directory)/multiply_js
	rmdir $(directory)
