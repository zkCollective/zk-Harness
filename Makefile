benchmark_directory = benchmarks

gnark_directory = gnark
circom_directory = circom
snarkjs_directory = snarkjs
bellman_directory = bellman
bellman_ce_directory = bellman
halo2_pse_directory = halo2_pse
MACHINE := $(shell cat machine 2> /dev/null || echo DEFAULT)


all: benchmark-gnark-arithmetics benchmark-gnark-ec benchmark-gnark-circuits benchmark-snarkjs-arithmetics benchmark-snarkjs-ec benchmark-circom-circuits benchmark-halo2-pse-circuits

ready: benchmark-bellman-circuits benchmark-halo2-pse-circuits benchmark-circom-circuits

benchmark-bellman-circuits:
	$(info --------------------------------------------)
	$(info ------    BELLMAN CIRCUIT BENCHMARKS  ------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/bellman/config_circuits.json --machine $(MACHINE)

benchmark-bellman-ce-circuits:
	$(info --------------------------------------------)
	$(info ------    BELLMAN_CE CIRCUIT BENCHMARKS ----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/bellman_ce/config_circuits.json --machine $(MACHINE)

benchmark-halo2-pse-circuits:
	$(info --------------------------------------------)
	$(info ----- HALO-PSE ARITHMETICS BENCHMARKS ------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/halo2_pse/config_circuits.json --machine $(MACHINE)

benchmark-snarkjs-arithmetics:
	$(info --------------------------------------------)
	$(info ------ SNARKJS ARITHMETICS BENCHMARKS ------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/snarkjs/config_arithmetics.json --machine $(MACHINE) 

benchmark-snarkjs-ec:
	$(info --------------------------------------------)
	$(info ---------- SNARKJS EC BENCHMARKS -----------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/snarkjs/config_ec.json --machine $(MACHINE)

benchmark-rapidsnark-arithmetics:
	$(info --------------------------------------------)
	$(info ---- RAPIDSNARK ARITHMETICS BENCHMARKS -----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/rapidsnark/config_arithmetics.json --machine $(MACHINE)

benchmark-rapidsnark-ec:
	$(info --------------------------------------------)
	$(info --------- RAPIDSNARK EC BENCHMARKS ---------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/rapidsnark/config_ec.json --machine $(MACHINE)

benchmark-toy-circom:
	$(info --------------------------------------------)
	$(info ---------- CIRCOM TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/circom/config_all_toy.json --machine $(MACHINE)

benchmark-exponentiate-circom:
	$(info --------------------------------------------)
	$(info ----- CIRCOM EXPONENTIATE BENCHMARKS -------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/circom/config_exponentiate.json --machine $(MACHINE)

benchmark-sha-circom:
	$(info --------------------------------------------)
	$(info -------- CIRCOM SHA256 BENCHMARKS ----------)
	$(info --------------------------------------------)
	orig_dir=$(shell pwd)
	cd circom/circuits/benchmarks && if [ ! -d "circomlib" ]; then git clone https://github.com/iden3/circomlib.git; fi
	cd $(orig_dir)
	python3 -m _scripts.reader --config _input/config/circom/config_sha.json --machine $(MACHINE)
	rm -rf circom/circuits/benchmarks/circomlib

benchmark-circom-circuits: benchmark-exponentiate-circom benchmark-sha-circom

benchmark-gnark-circuits: 
	$(info --------------------------------------------)
	$(info -------- GNARK CIRCUITS BENCHMARKS ---------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_circuits.json --machine $(MACHINE)

benchmark-gnark-arithmetics:
	$(info --------------------------------------------)
	$(info ------- GNARK ARITHMETICS BENCHMARKS -------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_arithmetics.json --machine $(MACHINE)

benchmark-gnark-ec:
	$(info --------------------------------------------)
	$(info ------ GNARK EC BENCHMARKS -----------------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_ec.json --machine $(MACHINE)

benchmark-toy-gnark:
	$(info --------------------------------------------)
	$(info ----------- GNARK TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_all_toy.json --machine $(MACHINE)

benchmark-gnark-hash:
	python3 -m _scripts.reader --config _input/config/gnark/config_hash.json --machine $(MACHINE)

benchmark-gnark-recursion:
	$(info --------------------------------------------)
	$(info ----------- GNARK RECURSION BENCHMARKS -----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_recursion.json --machine $(MACHINE)

test-simple:
	python3 -m _scripts.reader --config _input/config/gnark/config_gnark_simple.json --machine $(MACHINE)

clean:
	rm -rf $(benchmark_directory)/*
