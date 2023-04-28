SHELL = zsh

benchmark_directory = benchmarks

arkworks_directory = arkworks
blstrs_directory = blstrs
curve25519_dalek_directory = curve25519-dalek
pasta_curves_directory = pasta_curves
zkcrypto_directory = zkcrypto
gnark_directory = gnark
circom_directory = circom
snarkjs_directory = snarkjs
bellman_ce_directory = bellman_ce

arkworks_benchmarks_directory = $(benchmark_directory)/$(arkworks_directory)
blstrs_benchmarks_directory = $(benchmark_directory)/$(blstrs_directory)
bellman_ce_benchmarks_directory = $(benchmark_directory)/$(bellman_ce_directory)
curve25519_dalek_benchmarks_directory = $(benchmark_directory)/$(curve25519_dalek_directory)
pasta_curves_benchmarks_directory = $(benchmark_directory)/$(pasta_curves_directory)
zkcrypto_benchmarks_directory = $(benchmark_directory)/$(zkcrypto_directory)
gnark_benchmarks_directory = $(benchmark_directory)/$(gnark_directory)
circom_benchmarks_directory = $(benchmark_directory)/$(circom_directory)
snarkjs_benchmarks_directory = $(benchmark_directory)/$(snarkjs_directory)


all: init arkworks-arithmetics blstrs-arithmetics benchmark-gnark-arithmetics benchmark-gnark-ec benchmark-gnark-circuits benchmark-snarkjs-arithmetics benchmark-snarkjs-ec benchmark-circom-circuits

test: init bellman-ce-arithmetics

init:
	cargo install cargo-criterion
	mkdir -p $(arkworks_benchmarks_directory)
	mkdir -p $(blstrs_benchmarks_directory)
	mkdir -p $(curve25519_dalek_benchmarks_directory)
	mkdir -p $(pasta_curves_benchmarks_directory)
	mkdir -p $(zkcrypto_benchmarks_directory)
	mkdir -p $(bellman_ce_benchmarks_directory)

bellman-ce-arithmetics:
	cd bellman_ce; cargo criterion --message-format=json 1> ../$(bellman_ce_benchmarks_directory)/bellman_ce.json

arkworks-arithmetics: arkworks-curves
	cd arkworks; cargo criterion --message-format=json 1> $(arkworks_benchmarks_directory)/arkworks_arithmetics.json

arkworks-curves:
	rm -rf ./arkworks/curves
	git clone https://github.com/arkworks-rs/curves.git ./arkworks/curves || true
	cd ./arkworks/curves; cargo criterion --features ark-ec/parallel,ark-ff/asm --message-format=json 1> ../../$(arkworks_benchmarks_directory)/arkworks_curves.json

blstrs-arithmetics:
	cd blstrs; cargo criterion --message-format=json 1> ../$(blstrs_benchmarks_directory)/blstrs.json

curve25519-dalek-arithmetics:
	cd curve25519-dalek; cargo criterion --message-format=json 1> ../$(curve25519_dalek_benchmarks_directory)/curve25519-dalek.json

pasta-curves-arithmetics:
	cd pasta_curves; cargo criterion --message-format=json 1> ../$(pasta_curves_benchmarks_directory)/pasta_curves.json

zkcrypto-arithmetics:
	cd zkcrypto; cargo criterion --message-format=json 1> ../$(zkcrypto_benchmarks_directory)/zkcrypto.json

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
