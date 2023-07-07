ROOT_DIR=$(pwd)
benchmark_directory = benchmarks
MATH = math
MACHINE := $(shell cat machine 2> /dev/null || echo DEFAULT)

# Math variables
arkworks_directory = arkworks
arkworks_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(arkworks_directory)
arkworks_curves_directory = arkworks_curves
arkworks_curves_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(arkworks_curves_directory)
blstrs_directory = blstrs
blstrs_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(blstrs_directory)
curve25519_dalek_directory = curve25519-dalek
curve25519_dalek_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(curve25519_dalek_directory)
pasta_curves_directory = pasta_curves
pasta_curves_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(pasta_curves_directory)
halo2_curves_directory = halo2_curves
halo2_curves_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(halo2_curves_directory)
zkcrypto_directory = zkcrypto
zkcrypto_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(zkcrypto_directory)
pairing_ce_directory = pairing_ce
pairing_ce_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(pairing_ce_directory)
gnark_crypto_directory = gnark_crypto
gnark_crypto_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(gnark_crypto_directory)
ffjavascript_directory = ffjavascript
ffjavascript_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(ffjavascript_directory)
ffiasm_directory = ffiasm
ffiasm_benchmarks_directory = $(benchmark_directory)/$(MATH)/$(MACHINE)/$(ffiasm_directory)

# Circuits variables

all: init math circuits

init:
	cargo install cargo-criterion
	@if [ ! -d "math" ]; then mkdir -p math; fi

############################# ARITHMETICS ######################################
# TODO include ffiasm
math: math-arkworks math-arkworks-curves math-blstrs math-curve25519-dalek math-pasta-curves math-halo2-curves math-zkcrypto math-pairing-ce math-ffjavascript math-gnark

math-arkworks:
	$(info --------------------------------------------)
	$(info --------- ARKWORKS MATH BENCHMARKS ---------)
	$(info --------------------------------------------)
	mkdir -p $(arkworks_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/arkworks; cargo criterion --message-format=json 1> ../../../../$(arkworks_benchmarks_directory)/arkworks.json || true

math-arkworks-curves:
	$(info --------------------------------------------)
	$(info ------ ARKWORKS CURVES MATH BENCHMARKS -----)
	$(info --------------------------------------------)
	mkdir -p $(arkworks_curves_benchmarks_directory)
	cd math/zkalc/backend && make init
	cd math/zkalc/backend && git clone https://github.com/arkworks-rs/curves.git || true
	cd math/zkalc/backend/curves; git fetch; git checkout releases; cargo criterion --features ark-ec/parallel,ark-ff/asm --message-format=json 1> ../../../../$(arkworks_curves_benchmarks_directory)/ark-curves.json || true

math-blstrs:
	$(info --------------------------------------------)
	$(info -------- BLSTRS MATH BENCHMARKS ------------)
	$(info --------------------------------------------)
	mkdir -p $(blstrs_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/blstrs; cargo criterion --message-format=json 1> ../../../../$(blstrs_benchmarks_directory)/blstrs.json || true

math-curve25519-dalek:
	$(info --------------------------------------------)
	$(info ---- curve25519-dalek MATH BENCHMARKS ------)
	$(info --------------------------------------------)
	mkdir -p $(curve25519_dalek_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/curve25519-dalek; cargo criterion --message-format=json 1> ../../../../$(curve25519_dalek_benchmarks_directory)/curve25519-dalek.json || true

math-pasta-curves:
	$(info --------------------------------------------)
	$(info ----------- PASTA MATH BENCHMARKS ----------)
	$(info --------------------------------------------)
	mkdir -p $(pasta_curves_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/pasta_curves; cargo criterion --message-format=json 1> ../../../../$(pasta_curves_benchmarks_directory)/pasta_curves.json || true

math-halo2-curves:
	$(info --------------------------------------------)
	$(info ----------- HALO2 MATH BENCHMARKS ----------)
	$(info --------------------------------------------)
	mkdir -p $(halo2_curves_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/halo2_curves; cargo criterion --message-format=json 1> ../../../../$(halo2_curves_benchmarks_directory)/halo2_curves.json || true

math-zkcrypto:
	$(info --------------------------------------------)
	$(info -------- ZKCRYPTO MATH BENCHMARKS ----------)
	$(info --------------------------------------------)
	mkdir -p $(zkcrypto_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/zkcrypto; cargo criterion --message-format=json 1> ../../../../$(zkcrypto_benchmarks_directory)/zkcrypto.json || true

math-pairing-ce:
	$(info --------------------------------------------)
	$(info --------- Pairing CE MATHBENCHMARKS --------)
	$(info --------------------------------------------)
	mkdir -p $(pairing_ce_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/pairing_ce; cargo criterion --message-format=json 1> ../../../../$(pairing_ce_benchmarks_directory)/pairing_ce.json || true

math-ffjavascript:
	$(info --------------------------------------------)
	$(info ------- ffjavascript MATHBENCHMARKS --------)
	$(info --------------------------------------------)
	mkdir -p $(ffjavascript_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/ffjavascript; node bench.js > ../../../../$(ffjavascript_benchmarks_directory)/ffjavascript.json || true

math-gnark:
	$(info --------------------------------------------)
	$(info ----------- gnark MATHBENCHMARKS -----------)
	$(info --------------------------------------------)
	mkdir -p $(gnark_crypto_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend; if [ ! -d "gnark-crypto" ]; then git clone -b zkalc https://github.com/ConsenSys/gnark-crypto.git; fi
	cd math/zkalc/backend/gnark-crypto && \
	bash ./zkalc.sh bls12-381 | tee ../../../../$(gnark_crypto_benchmarks_directory)/gnark-bls12-381.txt && \
	bash ./zkalc.sh bls12-377 | tee ../../../../$(gnark_crypto_benchmarks_directory)/gnark-bls12-377.txt && \
	bash ./zkalc.sh bn254     | tee ../../../../$(gnark_crypto_benchmarks_directory)/gnark-bn254.txt && \
	bash ./zkalc.sh secp256k1 | tee ../../../../$(gnark_crypto_benchmarks_directory)/gnark-secp256k1.txt

math-ffiasm:
	$(info --------------------------------------------)
	$(info ---------- ffiasm MATHBENCHMARKS -----------)
	$(info --------------------------------------------)
	mkdir -p $(ffiasm_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone https://github.com/asn-d6/zkalc.git; fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/ffiasm; node scripts/bench.js > ../../../../$(ffjavascript_benchmarks_directory)/ffjavascript.json || true

################################################################################

############################## CIRCUITS ########################################

circuits-test: benchmark-bellman-test-circuit benchmark-halo2-pse-test-circuit benchmark-circom-test-circuit benchmark-gnark-test-circuit

circuits: benchmark-bellman-test-circuit

benchmark-bellman-test-circuit:
	$(info --------------------------------------------)
	$(info ----- BELLMAN TEST CIRCUIT BENCHMARKS  -----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/bellman/config_test.json --machine $(MACHINE)

benchmark-bellman-circuits:
	$(info --------------------------------------------)
	$(info ------- BELLMAN CIRCUIT BENCHMARKS ---------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/bellman/config_circuits.json --machine $(MACHINE)

benchmark-halo2-pse-test-circuit:
	$(info --------------------------------------------)
	$(info ----- HALO-PSE TEST CIRCUIT BENCHMARKS -----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/halo2_pse/config_test.json --machine $(MACHINE)

benchmark-halo2-pse-circuits:
	$(info --------------------------------------------)
	$(info ----- HALO-PSE ARITHMETICS BENCHMARKS ------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/halo2_pse/config_circuits.json --machine $(MACHINE)

benchmark-circom-test-circuit:
	$(info --------------------------------------------)
	$(info ----- CIRCOM TEST CIRCUIT BENCHMARKS -------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/circom/config_test.json --machine $(MACHINE)

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

benchmark-gnark-test-circuit:
	$(info --------------------------------------------)
	$(info ------ GNARK TEST CIRCUIT BENCHMARKS -------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_test.json --machine $(MACHINE)

benchmark-gnark-circuits: 
	$(info --------------------------------------------)
	$(info -------- GNARK CIRCUITS BENCHMARKS ---------)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_circuits.json --machine $(MACHINE)

################################################################################

############################## RECURSION #######################################

benchmark-gnark-recursion:
	$(info --------------------------------------------)
	$(info ----------- GNARK RECURSION BENCHMARKS -----)
	$(info --------------------------------------------)
	python3 -m _scripts.reader --config _input/config/gnark/config_recursion.json --machine $(MACHINE)

################################################################################

clean:
	rm -rf $(benchmark_directory)/*
