SHELL = zsh

benchmark_directory = benchmarks
gnark_directory = gnark
gnark_benchmarks_directory = $(benchmark_directory)/$(gnark_directory)

all: benchmark-toy

benchmark-toy:
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

clean:
	rm -r $(gnark_benchmarks_directory)/*