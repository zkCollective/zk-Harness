<h1 align="center">zk-Harness</h1>

There is a large and increasing number of libraries that enable verifiable computation with *zero-knowledge succinct non-interactive arguments (zkSNARKs)*. 
Whereas the performance of zkSNARKs is well-understood in theory, various factors make it difficult to compare different proof systems without favoring some approaches over others.
Evaluating the practical performance of a specific library can be difficult due to various factors, such as the underlying elliptic curve, the proof system at hand, the arithmetization supported, or due to application-specific factors such as the desired security level.

**zk-Harness** is a benchmarking framework that aims to address these barriers by providing a unified
benchmark for standardized evaluation of existing libraries for zkSNARKS. 
It is designed to be easily extendable - libraries can be easily added and compared for standardized functionalities, whereas circuit developers can simply add their already developed circuit to evaluate its performance. 
zk-Harness provides benchmarks along the whole zkSNARK stack.
The backend for mathematical operations (field & curve) is inherited from [zka.lc](https://zka.lc/), a project from Michele Orrù and George Kadianakis.
The benchmark results are hosted at [zk-bench.org](https://www.zk-bench.org).

## Clone

```
git clone --recurse-submodules -j8 git@github.com:zkCollective/zk-Harness.git
# or 
git clone git@github.com:zkCollective/zk-Harness.git
git submodule update --init --recursive
```

## Update benchmarks

To pull the latest results run:

```
git submodule update --remote benchmarks
```

## Structure

* `benchmarks`: Directory containing the results of the benchmarks.
* `src`: Python code to run the benchmarks and parse the results
* `input`: Configuration and input files
  - `input/circuit`: Input files for each circuit, i.e., values to be used as the inputs in circuits for benchmarking.
  - `input/config`: configurations for executing benchmarks for a specific framework using `zkbench`
* `app`: UI code for presenting the results.
* `frameworks`: Directory containing the harness for each framework to benchmark circuits.
* `data`: Auxilary data
  - `data/circuits.json`: Supported zkSNARK frameworks
  - `data/math.json`: Supported math libraries
* `scripts`: Other auxiliary scripts

### Circuit payloads supported by zk-Harness 

The current framework supports a set of payloads for each library.

|          | Exponentiate        | SHA-256             |
| -------- | ------------------- | ------------------- |
| [Bellman](https://github.com/zkcrypto/bellman) | :heavy_check_mark: (custom) | :heavy_check_mark: ([implementation](https://github.com/zkcrypto/bellman/blob/main/src/gadgets/sha256.rs)) |
| [Circom](https://github.com/iden3/circom) | :heavy_check_mark: (custom) | :heavy_check_mark: ([implementation](https://github.com/iden3/circomlib/tree/master/circuits/sha256)) |
| [Gnark](https://github.com/Consensys/gnark) | :heavy_check_mark: (custom) | :heavy_check_mark: (custom) |
| [Halo2-PSE](https://github.com/privacy-scaling-explorations/halo2/) | :heavy_check_mark: (custom) | :heavy_check_mark: ([implementation](https://github.com/privacy-scaling-explorations/halo2/blob/main/halo2_gadgets/benches/sha256.rs)) |
| [Starky](https://github.com/mir-protocol/plonky2) | :heavy_check_mark: (custom) | :x: |

### Curves and Fields 

|           | Language | Curves/Fields | Frameworks |
| --------- | -------- | ------------- | ---------- |
| [blstrs](https://github.com/filecoin-project/blstrs) | Rust | BLS12-381 |  |
| [gnark-crypto](https://github.com/Consensys/gnark-crypto) | Go | BN254, BLS12-377, BLS12-378, BLS12-381, BLS12-387, BLS24-315, BLS24-317, BW6-761, BW6-756, BW6-633, secp256k1, stark-curve, goldilocks | [gnark](https://github.com/Consensys/gnark) |
| [arkworks-curves](https://github.com/arkworks-rs/curves) | Rust | BN254, BLS12-377, BLS12-381, MNT4-298, MNT4-753, MNT6-298, MNT6-753, Grumpkin, BW6-761, CP6-782, secp256k1, secp256r1, secp384r1, secq256k1 | [arkworks](https://github.com/arkworks-rs/snark) |
| [curve25519-dalek](https://github.com/dalek-cryptography/curve25519-dalek) | Rust | Curve25519 |  |
| [ffjavascript](https://github.com/iden3/ffjavascript) | JavaScript/WASM | BN128, BLS12-381 | [snarkjs](https://github.com/iden3/snarkjs) |
| [ffiasm](https://github.com/iden3/ffiasm) | C++ | BN128, BLS12-381 | [rapidsnark](https://github.com/iden3/rapidsnark) |
| [halo2curves](https://github.com/privacy-scaling-explorations/halo2curves) | Rust | BN256, Pallas, Vesta | [halo2-PSE](https://github.com/privacy-scaling-explorations/halo2) |
| [pairing_ce](https://github.com/matter-labs/pairing) | Rust | BN256, BLS12-381 | [bellman-ce](https://github.com/matter-labs/bellman) |
| [pairing](https://github.com/zkcrypto/pairing) | Rust | jubjub, BLS12-381 | [bellman](https://github.com/zkcrypto/bellman) |
| [pasta_curves](https://github.com/zcash/pasta_curves) | Rust | Pallas, Vesta | [halo2](https://github.com/zcash/halo2) |

## How to use

To run the benchmarks, you will first need to follow the [installation instructions](INSTALL.md).

To run *all* benchmarks for mathematical operations, run `make math`.

Run *test* benchmarks for end-to-end circuits, run `make circuits-test`.

To run *all* benchmarks for end-to-end circuits, run `make circuits`.

To keep logs for the runs you can use the `tee` command, e.g., `make circuits 2>&1 | tee logs`

### Add benchmarks for new circuits

zk-Harness is easily extensible and modular, which means that you should be able to integrate you own circuits with ease.
Each framework in `framework/<framework_name>` includes a detailed description on how to add a self-developed circuit that goes beyond the standard payloads already integrated.

### zk-Harness architecture and add support for new frameworks

See [ARCHITECTURE.md](ARCHITECTURE.md)

## How to contribute / TODOs

We aim to successively expand this list to further include benchmarks for other ZKP frameworks, more circuits, and recursive composition of proofs.

There are many ways in which you can contribute to the zk-Harness:

- [] Add benchmarks for circuits in an already integrated framework
- [] Integrate a new framework into the zk-Harness
- [] Integrate new math benchmarks in [zkalc](https://github.com/asn-d6/zkalc/)
- [] Run existing benchmarks to additional machines
- [] Work on any open GitHub issue
- [] Propose new visualizations for the results
- [] Propose new benchmark categories, such as recursion benchmarking

Please read the [Contribution Guidelines](https://github.com/zkCollective/zk-Harness/blob/main/CONTRIBUTING.md) before creating a PR or opening an issue.

zk-Harness is developed as part of the [zk-Hackathon](https://rdi.berkeley.edu/zkp-web3-hackathon/) hosted by the [Berkeley Center for Responsible Decentralized Intelligence](https://rdi.berkeley.edu/).
Further, zk-Harness is part of the [zk-Collective](https://github.com/zkCollective/).
