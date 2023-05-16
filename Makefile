benchmark_directory = benchmarks

gnark_directory = gnark
circom_directory = circom
snarkjs_directory = snarkjs
bellman_directory = bellman
bellman_ce_directory = bellman
gnark_benchmarks_directory = $(benchmark_directory)/$(gnark_directory)
circom_benchmarks_directory = $(benchmark_directory)/$(circom_directory)
snarkjs_benchmarks_directory = $(benchmark_directory)/$(snarkjs_directory)
bellman_ce_benchmarks_directory = $(benchmark_directory)/$(bellman_ce_directory)
bellman_benchmarks_directory = $(benchmark_directory)/$(bellman_directory)
halo2_pse_directory = halo2_pse
halo2_pse_benchmarks_directory = $(benchmark_directory)/$(halo2_pse_directory)


all: benchmark-gnark-arithmetics benchmark-gnark-ec benchmark-gnark-circuits benchmark-snarkjs-arithmetics benchmark-snarkjs-ec benchmark-circom-circuits benchmark-halo2-pse-circuits

benchmark-bellman-circuits:
	$(info --------------------------------------------)
	$(info ------    BELLMAN CIRCUIT BENCHMARKS  ------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/bellman/config_circuits.json

benchmark-bellman-ce-circuits:
	$(info --------------------------------------------)
	$(info ------    BELLMAN_CE CIRCUIT BENCHMARKS ----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/bellman_ce/config_circuits.json

benchmark-halo2-pse-circuits:
	$(info --------------------------------------------)
	$(info ----- HALO-PSE ARITHMETICS BENCHMARKS ------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/halo2_pse/config_circuits.json

benchmark-bellman-circuits:
	$(info --------------------------------------------)
	$(info ------    BELLMAN CIRCUIT BENCHMARKS  ------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/bellman/config_circuits.json

benchmark-bellman-ce-circuits:
	$(info --------------------------------------------)
	$(info ------    BELLMAN_CE CIRCUIT BENCHMARKS ----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/bellman_ce/config_circuits.json

benchmark-halo2-pse-circuits:
	$(info --------------------------------------------)
	$(info ----- HALO-PSE ARITHMETICS BENCHMARKS ------)
	$(info --------------------------------------------)
	cd $(halo2_pse_directory); \
		cargo run --bin exponentiation --release -- --input ../_input/circuit/exponentiate/input_1.json --output ../benchmarks/halo2_pse/jsons/exponentiate_input_1.json; \
		INPUT_FILE=../_input/circuit/exponentiate/input_1.json cargo criterion --message-format=json --bench exponentiation_bench 1> ../benchmarks/halo2_pse/jsons/exponentiate_input_1_bench.json

benchmark-snarkjs-arithmetics:
	$(info --------------------------------------------)
	$(info ------ SNARKJS ARITHMETICS BENCHMARKS ------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/snarkjs/config_arithmetics.json  

benchmark-snarkjs-ec:
	$(info --------------------------------------------)
	$(info ---------- SNARKJS EC BENCHMARKS -----------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/snarkjs/config_ec.json  

benchmark-rapidsnark-arithmetics:
	$(info --------------------------------------------)
	$(info ---- RAPIDSNARK ARITHMETICS BENCHMARKS -----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/rapidsnark/config_arithmetics.json  

benchmark-rapidsnark-ec:
	$(info --------------------------------------------)
	$(info --------- RAPIDSNARK EC BENCHMARKS ---------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/rapidsnark/config_ec.json  

benchmark-toy-circom:
	$(info --------------------------------------------)
	$(info ---------- CIRCOM TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/circom/config_all_toy.json  

benchmark-circom-circuits: benchmark-toy-circom

benchmark-gnark-arithmetics:
	$(info --------------------------------------------)
	$(info ------- GNARK ARITHMETICS BENCHMARKS -------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_arithmetics.json  

benchmark-gnark-ec:
	$(info --------------------------------------------)
	$(info ------ GNARK EC BENCHMARKS -----------------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_ec.json  

benchmark-toy-gnark:
	$(info --------------------------------------------)
	$(info ----------- GNARK TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_all_toy.json  

benchmark-hash:
	python3 -m _scripts.reader --config _input/config/gnark/config_hash.json  

benchmark-gnark-circuits: benchmark-toy-gnark benchmark-hash

test-simple:
	python3 -m _scripts.reader --config _input/config/gnark/config_gnark_simple.json  

clean:
	rm -rf $(gnark_benchmarks_directory)/*  $(circom_benchmarks_directory)/* $(snarkjs_benchmarks_directory)/* $(halo2_pse_benchmarks_directory)/*
