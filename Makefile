SHELL = zsh

benchmark_directory = benchmarks

gnark_directory = gnark
gnark_benchmarks_directory = $(benchmark_directory)/$(gnark_directory)

circom_directory = circom
circom_benchmarks_directory = $(benchmark_directory)/$(circom_directory)
circom_script = $(circom_directory)/scripts/run_circuit.sh
circom_circuits = $(circom_directory)/circuits/benchmarks
circom_inputs = $(circom_directory)/inputs/benchmarks
circom_ptau = $(circom_directory)/phase1/powersOfTau28_hez_final_20.ptau


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
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=setup 
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=prove 
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=verify 
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=setup --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=prove --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=verify --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=setup --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=prove --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=verify --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=setup --curve=bls24_315
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=prove --curve=bls24_315
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=verify --curve=bls24_315
	# cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=setup --curve=bw6_756
	# cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=prove --curve=bw6_756
	# cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=verify --curve=bw6_756
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=setup 
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=prove
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=verify 
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=setup --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=prove --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=verify --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=setup --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=prove --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=verify --curve=bw6_761
	# cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=setup --curve=bls12_378
	# cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=prove --curve=bls12_378
	# cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=verify --curve=bls12_378
	# cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=setup --curve=bw6_756
	# cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=prove --curve=bw6_756
	# cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=verify --curve=bw6_756

	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=setup 
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=prove 
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=verify 
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=setup --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=prove --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=verify --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=setup --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=prove --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=verify --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=setup --curve=bls24_315
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=prove --curve=bls24_315
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=verify --curve=bls24_315
	# cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=setup --curve=bw6_756
	# cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=prove --curve=bw6_756
	# cd $(gnark_directory); go build; ./gnark groth16 --circuit=cubic --algo=verify --curve=bw6_756
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=setup 
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=prove
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=verify 
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=setup --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=prove --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=verify --curve=bls12_377
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=setup --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=prove --curve=bw6_761
	cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=verify --curve=bw6_761
	# cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=setup --curve=bls12_378
	# cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=prove --curve=bls12_378
	# cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=verify --curve=bls12_378
	# cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=setup --curve=bw6_756
	# cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=prove --curve=bw6_756
	# cd $(gnark_directory); go build; ./gnark plonk --circuit cubic --algo=verify --curve=bw6_756

clean:
	rm -r $(gnark_benchmarks_directory)/*
