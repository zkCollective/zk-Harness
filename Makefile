ROOT_DIR=$(pwd)
benchmark_directory = benchmarks

arkworks_directory = arkworks
blstrs_directory = blstrs
curve25519_dalek_directory = curve25519-dalek
pasta_curves_directory = pasta_curves
zkcrypto_directory = zkcrypto
gnark_directory = gnark
circom_directory = circom
snarkjs_directory = snarkjs
pairing_ce_directory = pairing_ce
halo2_curves_directory = halo2_curves

arkworks_benchmarks_directory = $(benchmark_directory)/$(arkworks_directory)
blstrs_benchmarks_directory = $(benchmark_directory)/$(blstrs_directory)
pairing_ce_benchmarks_directory = $(benchmark_directory)/$(pairing_ce_directory)
halo2_curves_benchmarks_directory = $(benchmark_directory)/$(halo2_curves_directory)
curve25519_dalek_benchmarks_directory = $(benchmark_directory)/$(curve25519_dalek_directory)
pasta_curves_benchmarks_directory = $(benchmark_directory)/$(pasta_curves_directory)
zkcrypto_benchmarks_directory = $(benchmark_directory)/$(zkcrypto_directory)
gnark_benchmarks_directory = $(benchmark_directory)/$(gnark_directory)
circom_benchmarks_directory = $(benchmark_directory)/$(circom_directory)
snarkjs_benchmarks_directory = $(benchmark_directory)/$(snarkjs_directory)

bellman_ce_directory = bellman_ce
bellman_ce_directory = bellman
gnark_benchmarks_directory = $(benchmark_directory)/$(gnark_directory)
circom_benchmarks_directory = $(benchmark_directory)/$(circom_directory)
snarkjs_benchmarks_directory = $(benchmark_directory)/$(snarkjs_directory)
bellman_ce_benchmarks_directory = $(benchmark_directory)/$(bellman_ce_directory)
bellman_benchmarks_directory = $(benchmark_directory)/$(bellman_directory)
halo2_pse_directory = halo2_pse
halo2_pse_benchmarks_directory = $(benchmark_directory)/$(halo2_pse_directory)


all: init arkworks-arithmetics blstrs-arithmetics benchmark-gnark-arithmetics benchmark-gnark-ec benchmark-gnark-circuits benchmark-snarkjs-arithmetics benchmark-snarkjs-ec benchmark-circom-circuits

test: init bellman-ce-arithmetics

init:
	cargo install cargo-criterion
	mkdir -p $(arkworks_benchmarks_directory)
	mkdir -p $(blstrs_benchmarks_directory)
	mkdir -p $(curve25519_dalek_benchmarks_directory)
	mkdir -p $(pasta_curves_benchmarks_directory)
	mkdir -p $(zkcrypto_benchmarks_directory)
	mkdir -p $(pairing_ce_benchmarks_directory)
	mkdir -p $(halo2_curves_benchmarks_directory)

halo2-curves-arithmetics:
	$(info --------------------------------------------)
	$(info --- HALO2 CURVES ARITHMETICS BENCHMARKS ----)
	$(info --------------------------------------------)
	cd math/halo2_curves; cargo criterion --message-format=json 1> ../../$(halo2_curves_benchmarks_directory)/halo2_curves.json

pairing-ce-arithmetics:
	$(info --------------------------------------------)
	$(info --- Pairing CE   ARITHMETICS BENCHMARKS ----)
	$(info --------------------------------------------)
	cd math/pairing_ce; cargo criterion --message-format=json 1> ../../$(pairing_ce_benchmarks_directory)/pairing_ce.json

blstrs-arithmetics:
	$(info --------------------------------------------)
	$(info --- BLSTRS ARITHMETICS BENCHMARKS ----------)
	$(info --------------------------------------------)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend/blstrs; cargo criterion --message-format=json 1> ../../../../$(blstrs_benchmarks_directory)/blstrs.json

curve25519-dalek-arithmetics:
	$(info --------------------------------------------)
	$(info -- curve25519-dalek ARITHMETICS BENCHMARKS -)
	$(info --------------------------------------------)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend/curve25519-dalek; cargo criterion --message-format=json 1> ../../../../$(curve25519_dalek_benchmarks_directory)/curve25519-dalek.json

pasta-curves-arithmetics:
	$(info --------------------------------------------)
	$(info ------ PASTA ARITHMETICS BENCHMARKS ------)
	$(info --------------------------------------------)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend/pasta_curves; cargo criterion --message-format=json 1> ../../../../$(pasta_curves_benchmarks_directory)/pasta_curves.json

zkcrypto-arithmetics:
	$(info --------------------------------------------)
	$(info ------ ZKCRYPTO ARITHMETICS BENCHMARKS -----)
	$(info --------------------------------------------)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend/zkcrypto; cargo criterion --message-format=json 1> ../../../../$(zkcrypto_benchmarks_directory)/zkcrypto.json

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

benchmark-circomlib:
	$(info --------------------------------------------)
	$(info ---------- CIRCOM TOY BENCHMARKS -----------)
	$(info --------------------------------------------)
	orig_dir=$(shell pwd)
	cd circom/circuits/benchmarks && if [ ! -d "circomlib" ]; then git clone https://github.com/iden3/circomlib.git; fi
	cd $(orig_dir)
	python3 -m _scripts.reader --config _input/config/circom/config_sha.json
	rm -rf circom/circuits/benchmarks/circomlib

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
