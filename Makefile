benchmark_directory = benchmarks

gnark_directory = gnark
circom_directory = circom
snarkjs_directory = snarkjs
bellman_ce_directory = bellman_ce
gnark_benchmarks_directory = $(benchmark_directory)/$(gnark_directory)
circom_benchmarks_directory = $(benchmark_directory)/$(circom_directory)
snarkjs_benchmarks_directory = $(benchmark_directory)/$(snarkjs_directory)
bellman_ce_benchmarks_directory = $(benchmark_directory)/$(bellman_ce_directory)


all: benchmark-gnark-arithmetics benchmark-gnark-ec benchmark-gnark-circuits benchmark-snarkjs-arithmetics benchmark-snarkjs-ec benchmark-circom-circuits

bellman-ce-circuits:
	cd bellman_ce_circuits; cargo criterion --message-format=json 1> ../$(bellman_ce_benchmarks_directory)/bellman_ce_circuits.json

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
	rm -rf $(gnark_benchmarks_directory)/*  $(circom_benchmarks_directory)/* $(snarkjs_benchmarks_directory)/*
