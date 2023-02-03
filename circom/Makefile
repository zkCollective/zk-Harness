SHELL = zsh

circom = circuits/toy/sudoku.circom
p_input = circuits/toy/sudoku.input.json
directory = ./circuits/toy/tmp
r1cs = $(directory)/sudoku.r1cs
wasm = $(directory)/sudoku_js/sudoku.wasm
wit_gen = $(directory)/sudoku_js/generate_witness.js
compile_outputs = $(directory)/sudoku_js/witness_calculator.js $(r1cs) $(wasm) $(wit_gen)
pk = $(directory)/sudoku.pk
vk = $(directory)/sudoku.vk
ptau = $(directory)/sudoku.ptau
keys = $(pk) $(vk)
wit = $(directory)/sudoku.wtns
pf = $(directory)/sudoku.pf.json
inst = $(directory)/sudoku.inst.json

$(compile_outputs): $(circom)
		mkdir $(directory)
		circom $< --r1cs --wasm -o $(directory)

.PHONY = setup prove verify clean

setup: $(r1cs) $(wasm) $(wit_gen)
	test -d $(directory)
	node $(wit_gen) $(wasm) $(p_input) $(wit)
	# Trusted Setup
	snarkjs powersoftau new bn128 12 $(directory)/tmp.ptau
	snarkjs powersoftau prepare phase2 $(directory)/tmp.ptau $(ptau)
	rm $(directory)/tmp.ptau
	# Key Gen
	snarkjs groth16 setup $(r1cs) $(ptau) $(pk)
	snarkjs zkey export verificationkey $(keys)

prove: $(wit) $(pk)
	snarkjs groth16 prove $(pk) $(wit) $(pf) $(inst)

verify: $(inst) $(pf)
	snarkjs groth16 verify $(vk) $(inst) $(pf)


clean:
	rm -r $(compile_outputs) $(ptau) $(keys) $(wit) $(pf) $(inst)
	rmdir $(directory)/sudoku_js
	rmdir $(directory)
