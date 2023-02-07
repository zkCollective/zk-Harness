SHELL = zsh

gnark_directory = gnark

all: benchmark-toy

benchmark-toy:
	$(info --------------------------------------------)
	$(info ----------- GNARK TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=setup 
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=prove 
	cd $(gnark_directory); go build; ./gnark groth16 --circuit=expo --algo=verify 
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=setup 
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=prove
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo --algo=verify 