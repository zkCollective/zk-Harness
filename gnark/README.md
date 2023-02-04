# gnark Library

The development version documentation of gnark can be found [gnark](https://docs.gnark.consensys.net/en/latest/Concepts/schemes_curves/)

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
  - By default, this this tests on all curves and proving schemes supported by gnark by leveraging the ``` github.com/ConsenSys/gnark/test ``` package and the ``` AssertNew ``` method. See detailed descriptiont on what exactly is tested in ``` toy/cubic/cubic_test.go ```
