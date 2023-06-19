# Benchmarking the gnark Library

The development version documentation of gnark can be found [here](https://docs.gnark.consensys.net/en/latest/Concepts/schemes_curves/).
Note: this is an integration of the [gnark-bench](https://github.com/ConsenSys/gnark-bench), extended and adapted to zk-Harness.

## Plain Setup

### Installation

Installation and setup descriptions can be found [here](https://docs.gnark.consensys.net/en/latest/HowTo/get_started/)
gnark is written in Golang and hence requires the system wide installation of [go](https://go.dev/doc/install) to compile circuits and run proofs.

To write gnark code, the gnark module needs needs to be installed by running ``` go get github.com/consensys/gnark@v0.7.0 ```.

### Compilation & Proof

In general, to create a proof in gnark the following steps need to be completed - 1) Writing the circuit, 2) Compiling the circuit to an intermediary representation, 3) Debugging and Testing the circuit, 4) Create Proofs, 5) Verify Proofs.

To run the provided code you can leverage the ``` Makefile ``` with the following commands:

- ``` make test-toy```
  - Runs all toy examples provided in the ``` /toy ``` directory for gnark
  - By default, this tests on all curves and proving schemes supported by gnark by leveraging the ``` github.com/ConsenSys/gnark/test ``` package and the ``` AssertNew ``` method. See detailed descriptiont on what exactly is tested in ``` toy/cubic/cubic_test.go ```.

- ``` make test-toy```
  - Test all prf implementations

Running a benchmark can be facilitated through the exposed Cobra CLI:

``./gnark groth16 --circuit=sha256 --input=_input/circuit/sha256/input_3.json --curve=bn254``

would run an evaluation of SHA-256 circuit with the specified input over curve BN254.

### Memory Benchmarks 

Memory is benchmarked by buildiong independent binaries, which include serialization and deserialization of the constraint system, witness and proof.

Run ```./build_memory.sh ``` to build the relevant binaries.

Run e.g. ```/usr/bin/time -h -l ./groth16_memory_compile --circuit=sha256 --curve=bn254 --input=_input/circuit/sha256/input_10.json --count=10``` to benchmark realisitc memory consumption of circuit compilation.

Note that the overhead of the cobra CLI code is marginal, e.g. for an input of size $$2^{10}$$, the in-code memory consumption is roughly 790 MB for circuit compilation without serialization.
On the other hand, when including the code for serialization to a ```.dat``` file, max RSS memory consumption is roughly 2.2 GB, which is a significant overhead.

You can also run ```make benchmark-hash``` in the root directory to generate memory and time benchmarks for the same inputs on the SHA-256 circuit.

FIXME - Deleting the previously created  ```gnark_groth16_sha256.csv``` is mandatory at this point, otherwise memory benchmarks are overwritten and falsified. Hence, binary benchmarks are added in a new column 'ramReal'.

## Adding new circuits

See `TUTORIAL.md`

