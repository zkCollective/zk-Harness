all: sudoku cubic

sudoku:
	./scripts/run_circuit.sh circuits/toy/sudoku.circom inputs/toy/sudoku.input.json phase1/powersOfTau28_hez_final_20.ptau

cubic:
	./scripts/run_circuit.sh circuits/benchmarks/cubic/circuit.circom circuits/benchmarks/cubic/input.json phase1/powersOfTau28_hez_final_20.ptau

clean:
	rm -rf tmp 
