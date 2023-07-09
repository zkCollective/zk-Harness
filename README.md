<h1 align="center">zk-Harness</h1>

zk-Harness is a benchmarking framework for *zero-knowledge succinct non-interactive arguments (zkSNARKs)*. 
This repository contains a modular and easily extensible framework for benchmarking zkSNARKs and underlying mathematical primitives.

## Overview

TODO: goal, math (shared with zkalc), circuits, website of zk-bench, website of zkalc.
The benchmark results are hosted at [zk-bench.org](https://www.zk-bench.org).

The backend for mathematical operations has been merged with [zka.lc](https://zka.lc/).

## Structure

* `benchmarks`: Directory containing the results of the benchmarks.
* `src`: Python code to run the benchmarks and parse the results
* `input`: Configuration and input files
  - `input/circuit`: Input files for each circuit, i.e., values to be used as the inputs in circuits for benchmarking.
  - `input/config`: configurations for executing benchmarks for a specific framework using `zkbench`
* `app`: UI code for presenting the results.
* `frameworks`: Directory containing the harness for each framework to benchmark circuits.
* `data`: Auxilary data
  - `data/circuits.json`: Support zkSNARK frameworks
  - `data/math.json`: Support math libraries

### Circuit payloads supported by zk-Harness 

The current framework supports a set of payloads for each library.
We aim to successively extend the following once more circuits are available as `std` in the respective libraries.

|          | Exponentiate       | SHA-256           |
| -------- | ------------------ | ----------------- |
| `Bellman`| :heavy_check_mark: | :heavy_check_mark:|
| `Circom` | :heavy_check_mark: | :heavy_check_mark:|
| `Gnark`  | :heavy_check_mark: | :heavy_check_mark:|
| `Halo2`  | :heavy_check_mark: | :heavy_check_mark:|
| `Starky` | :heavy_check_mark: | :heavy_check_mark:|

### Curves, Fields and Arithmetizations and Backends

You can find the set of elliptic curves and finite fields as implemented in common libraries [here](https://docs.google.com/spreadsheets/d/1tq8lvcg88dE6D-EVJd61hBKhQxpDsZF16UMYpDXjef8/edit#gid=156416826).
We aim to maintain the set of supported functionalities and circuits in each library in the future.

## How to use

To run the benchmarks, you will first need to follow the installation instructions of the respective framework.

To run *all* benchmarks for mathematical operations, run `make math`.

To run *all* benchmarks for end-to-end circuits on standard operations, run `make benchmark-circuits`.

To obtain a Nix environment in which you can successfully run these benchmarks,
first [install Nix](https://nixos.org/download.html)
and [enable flakes](https://nixos.wiki/wiki/Flakes#Enable_flakes),
then run `nix develop`.

When running into issues for nix on M1/M2 Macs, please refer to [this issue](https://github.com/input-output-hk/plutus-pioneer-program/issues/40).  

### Run Benchmarks On Your Own!

zk-Harness is supposed to be easily extensible and modular, which means that you should be able to integrate you own circuits with ease.
Each framework in `framework/<framework_name>` includes a detailed description on how to add a self-developed circuit that goes beyond the standard payloads already integrated.


## Future Work & ZKP Hackathon

We aim to successively expand this list to further include benchmarks for other ZKP frameworks, recursive composition of proofs, and potentially zk-EVMs. 

zk-Harness was developed as a part of the ZKP / Web 3.0 Hackathon at UC Berkeley. You can find the program description detailing future integrations [here](https://drive.google.com/file/d/1Igm47dFXSOFAC_wldfUG4Y9OiITqlbQu/view).

## How to contribute

There are many ways in which you can contribute to the zk-Harness:

- Add benchmarks for circuits in an already integrated framework
- Integrate a new framework into the zk-Harness
- Propose new benchmark categories, such as for recursion and zk-EVMs.

Please read the [Contribution Guidelines](https://github.com/zkCollective/zk-Harness/blob/main/CONTRIBUTING.md) before creating a PR or opening an issue.

zk-Harness is developed as part of the [zk-Hackathon](https://rdi.berkeley.edu/zkp-web3-hackathon/) hosted by the [Berkeley Center for Responsible Decentralized Intelligence](https://rdi.berkeley.edu/).
Further, zk-Harness is part of the [zk-Collective](https://github.com/zkCollective/).
