ROOT_DIR=$(pwd)
benchmark_directory = benchmarks
MATH = math
MACHINE := $(shell cat machine 2> /dev/null || echo DEFAULT)
ZKALC = https://github.com/asn-d6/zkalc.git
INPUTS = input
FRAMEWORK = src

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

all: math circuits

init:
	cargo install cargo-criterion

math-init:
	@if [ ! -d "math" ]; then mkdir -p math; fi

############################# ARITHMETICS ######################################
# TODO include ffiasm
math: math-arkworks math-arkworks-curves math-blstrs math-curve25519-dalek math-pasta-curves math-halo2-curves math-zkcrypto math-pairing-ce math-ffjavascript math-gnark

math-arkworks: init math-init
	$(info --------------------------------------------)
	$(info --------- ARKWORKS MATH BENCHMARKS ---------)
	$(info --------------------------------------------)
	mkdir -p $(arkworks_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/arkworks; cargo criterion --message-format=json 1> ../../../../$(arkworks_benchmarks_directory)/arkworks.json || true

math-arkworks-curves: init math-init
	$(info --------------------------------------------)
	$(info ------ ARKWORKS CURVES MATH BENCHMARKS -----)
	$(info --------------------------------------------)
	mkdir -p $(arkworks_curves_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend && git clone https://github.com/arkworks-rs/curves.git || true
	cd math/zkalc/backend/curves; git fetch; git checkout releases; cargo criterion --features ark-ec/parallel,ark-ff/asm --message-format=json 1> ../../../../$(arkworks_curves_benchmarks_directory)/ark-curves.json || true

math-blstrs: init math-init
	$(info --------------------------------------------)
	$(info -------- BLSTRS MATH BENCHMARKS ------------)
	$(info --------------------------------------------)
	mkdir -p $(blstrs_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/blstrs; cargo criterion --message-format=json 1> ../../../../$(blstrs_benchmarks_directory)/blstrs.json || true

math-curve25519-dalek: init math-init
	$(info --------------------------------------------)
	$(info ---- curve25519-dalek MATH BENCHMARKS ------)
	$(info --------------------------------------------)
	mkdir -p $(curve25519_dalek_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/curve25519-dalek; cargo criterion --message-format=json 1> ../../../../$(curve25519_dalek_benchmarks_directory)/curve25519-dalek.json || true

math-pasta-curves: init math-init
	$(info --------------------------------------------)
	$(info ----------- PASTA MATH BENCHMARKS ----------)
	$(info --------------------------------------------)
	mkdir -p $(pasta_curves_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/pasta_curves; cargo criterion --message-format=json 1> ../../../../$(pasta_curves_benchmarks_directory)/pasta_curves.json || true

math-halo2-curves: init math-init
	$(info --------------------------------------------)
	$(info ----------- HALO2 MATH BENCHMARKS ----------)
	$(info --------------------------------------------)
	mkdir -p $(halo2_curves_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/halo2_curves; cargo criterion --message-format=json 1> ../../../../$(halo2_curves_benchmarks_directory)/halo2_curves.json || true

math-zkcrypto: init math-init
	$(info --------------------------------------------)
	$(info -------- ZKCRYPTO MATH BENCHMARKS ----------)
	$(info --------------------------------------------)
	mkdir -p $(zkcrypto_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/zkcrypto; cargo criterion --message-format=json 1> ../../../../$(zkcrypto_benchmarks_directory)/zkcrypto.json || true

math-pairing-ce: init math-init
	$(info --------------------------------------------)
	$(info --------- Pairing CE MATHBENCHMARKS --------)
	$(info --------------------------------------------)
	mkdir -p $(pairing_ce_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/pairing_ce; cargo criterion --message-format=json 1> ../../../../$(pairing_ce_benchmarks_directory)/pairing_ce.json || true

math-ffjavascript: math-init
	$(info --------------------------------------------)
	$(info ------- ffjavascript MATHBENCHMARKS --------)
	$(info --------------------------------------------)
	mkdir -p $(ffjavascript_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/ffjavascript; node bench.js > ../../../../$(ffjavascript_benchmarks_directory)/ffjavascript.json || true

math-gnark: math-init
	$(info --------------------------------------------)
	$(info ----------- gnark MATHBENCHMARKS -----------)
	$(info --------------------------------------------)
	mkdir -p $(gnark_crypto_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend; if [ ! -d "gnark-crypto" ]; then git clone -b zkalc https://github.com/ConsenSys/gnark-crypto.git; fi
	cd math/zkalc/backend/gnark-crypto && \
	bash ./zkalc.sh bls12-381 | tee ../../../../$(gnark_crypto_benchmarks_directory)/gnark-bls12-381.txt && \
	bash ./zkalc.sh bls12-377 | tee ../../../../$(gnark_crypto_benchmarks_directory)/gnark-bls12-377.txt && \
	bash ./zkalc.sh bn254     | tee ../../../../$(gnark_crypto_benchmarks_directory)/gnark-bn254.txt && \
	bash ./zkalc.sh secp256k1 | tee ../../../../$(gnark_crypto_benchmarks_directory)/gnark-secp256k1.txt

math-ffiasm: math-init
	$(info --------------------------------------------)
	$(info ---------- ffiasm MATHBENCHMARKS -----------)
	$(info --------------------------------------------)
	mkdir -p $(ffiasm_benchmarks_directory)
	cd math && if [ ! -d "zkalc" ]; then git clone $(ZKALC); fi
	cd math/zkalc/backend && make init
	cd math/zkalc/backend/ffiasm; node scripts/bench.js > ../../../../$(ffjavascript_benchmarks_directory)/ffjavascript.json || true

################################################################################

############################## CIRCUITS ########################################

circuits-test: benchmark-bellman-test-circuit benchmark-halo2-pse-test-circuit benchmark-circom-test-circuit gnark-init benchmark-gnark-test-circuit benchmark-starky-test-circuit

circuits: benchmark-bellman-circuits benchmark-halo2-pse-circuits benchmark-circom-circuits benchmark-gnark-circuits benchmark-starky-circuits

log-init:
	mkdir -p .logs

gnark-init: 
	cd frameworks/gnark && go build

circom-init:
	cd frameworks/circom/circuits/benchmarks && if [ ! -d "circomlib" ]; then git clone https://github.com/iden3/circomlib.git; fi

benchmark-bellman-test-circuit: init log-init
	$(info --------------------------------------------)
	$(info ----- BELLMAN TEST CIRCUIT BENCHMARKS  -----)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/bellman/config_test.json --machine $(MACHINE) 2>&1 | tee -a .logs/bellman.log

benchmark-bellman-circuits: init log-init
	$(info --------------------------------------------)
	$(info ------- BELLMAN CIRCUIT BENCHMARKS ---------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/bellman/config_circuits.json --machine $(MACHINE) 2>&1 | tee -a .logs/bellman.log

benchmark-halo2-pse-test-circuit: init log-init
	$(info --------------------------------------------)
	$(info ----- HALO2-PSE TEST CIRCUIT BENCHMARKS -----)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/halo2_pse/config_test.json --machine $(MACHINE) 2>&1 | tee -a .logs/halo2_pse.log

benchmark-halo2-pse-circuits: init log-init
	$(info --------------------------------------------)
	$(info ----- HALO2-PSE ARITHMETICS BENCHMARKS ------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/halo2_pse/config_circuits.json --machine $(MACHINE) 2>&1 | tee -a .logs/halo2_pse.log

benchmark-circom-test-circuit: init circom-init log-init
	$(info --------------------------------------------)
	$(info ----- CIRCOM TEST CIRCUIT BENCHMARKS -------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/circom/config_test.json --machine $(MACHINE) 2>&1 | tee -a .logs/circom.log

benchmark-exponentiate-circom: init log-init
	$(info --------------------------------------------)
	$(info ----- CIRCOM EXPONENTIATE BENCHMARKS -------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/circom/config_exponentiate.json --machine $(MACHINE) 2>&1 | tee -a .logs/circom.log

benchmark-sha-circom: init circom-init log-init
	$(info --------------------------------------------)
	$(info -------- CIRCOM SHA256 BENCHMARKS ----------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/circom/config_sha.json --machine $(MACHINE) 2>&1 | tee -a .logs/circom.log

benchmark-circom-circuits: init circom-init log-init
	$(info --------------------------------------------)
	$(info -------- CIRCOM CIRCUIT BENCHMARKS ---------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/circom/config_circuits.json --machine $(MACHINE) 2>&1 | tee -a .logs/circom.log

benchmark-gnark-test-circuit: gnark-init log-init
	$(info --------------------------------------------)
	$(info ------ GNARK TEST CIRCUIT BENCHMARKS -------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/gnark/config_test.json --machine $(MACHINE) 2>&1 | tee -a .logs/gnark.log

benchmark-gnark-circuits: gnark-init log-init
	$(info --------------------------------------------)
	$(info -------- GNARK CIRCUITS BENCHMARKS ---------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/gnark/config_circuits.json --machine $(MACHINE) 2>&1 | tee -a .logs/gnark.log

benchmark-starky-test-circuit: init log-init
	$(info --------------------------------------------)
	$(info -------- STARKY TEST CIRCUIT BENCHMARKS --------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/starky/config_test.json --machine $(MACHINE) 2>&1 | tee -a .logs/starky.log

benchmark-starky-circuits: init log-init
	$(info --------------------------------------------)
	$(info -------- STARKY CIRCUIT BENCHMARKS --------)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/starky/config_circuits.json --machine $(MACHINE) 2>&1 | tee -a .logs/starky.log

################################################################################

############################## RECURSION #######################################

benchmark-gnark-recursion: gnark-init
	$(info --------------------------------------------)
	$(info ----------- GNARK RECURSION BENCHMARKS -----)
	$(info --------------------------------------------)
	python3 -m $(FRAMEWORK).reader --config $(INPUTS)/config/gnark/config_recursion.json --machine $(MACHINE)

################################################################################

clean:
	rm -rf $(benchmark_directory)/*
