SHELL = zsh

gnark_directory = gnark

all: benchmark-toy

benchmark-toy:
	$(info --------------------------------------------)
	$(info ----------- GNARK TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	cd $(gnark_directory); go build; ./gnark groth16 --circuit expo
	cd $(gnark_directory); go build; ./gnark plonk --circuit expo