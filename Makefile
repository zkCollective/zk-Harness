SHELL = zsh

benchmark_directory = benchmarks

gnark_directory = gnark
gnark_benchmarks_directory = $(benchmark_directory)/$(gnark_directory)

circom_directory = circom
circom_benchmarks_directory = $(benchmark_directory)/$(circom_directory)
circom_script = $(circom_directory)/scripts/run_circuit.sh
circom_circuits = $(circom_directory)/circuits/benchmarks
circom_inputs = $(circom_directory)/inputs/benchmarks
circom_ptau = $(circom_directory)/phase1/powersOfTau28_hez_final_16.ptau 


all: benchmark-toy

benchmark-toy-circom:
	$(info --------------------------------------------)
	$(info ---------- CIRCOM TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	$(circom_script) $(circom_circuits)/cubic/circuit.circom $(circom_inputs)/cubic/input.json $(circom_ptau) $(circom_benchmarks_directory)/circom_cubic.csv; rm -rf tmp

benchmark-toy-gnark:
	$(info --------------------------------------------)
	$(info ----------- GNARK TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_all_toy.json  

benchmark-prf:
	python3 -m _scripts.reader --config _input/config/gnark/config_prf.json  

clean:
	rm -r $(gnark_benchmarks_directory)/*
