![Alt text](/logo_harness.png?raw=true "Title")

# zk-Harness - A benchmarking framework for general purpose Zero-Knowledge Proofs

We cordially invite the zk SNARK community to join us in creating a comprehensive benchmarking framework (zk-Harness) for zk SNARKs. This is a crucial step in the important mission to create a reference point for non-experts and experts alike on what zkSNARK scheme best suits their needs, and to also promote further research by identifying performance gaps. We believe that the collective efforts of the community will help to achieve this goal. Whether you are a researcher, developer, or simply passionate about zk SNARKs, we welcome your participation and contribution in this exciting initiative.

It is designed to be modular - new circuit implementations and ZKP-frameworks can be easily added, without extensive developer overhead.
zk-Harness has a standardized set of interfaces for configuring benchmark jobs and formatting log outputs.
Once a new component is included, it's benchmarks will be displayed on [zk-bench.org](https://www.zk-bench.org).

**Note: zk-Harness is a WIP. Its architecture may change over time.**

## Main Features

There is a large and ever-increasing number of SNARK implementations. Although the theoretical complexity of the underlying proof systems is well understood, the concrete costs rely on a number of factors such as the efficiency of the field and curve implementations, the underlying proof techniques, and the computation model and its compatibility with the specific application. To elicit the concrete performance differences in different proof systems, it is important to separately benchmark the following:

### Field and Curve computations

All popular SNARKs operate over prime fields, which are basically integers modulo p, i.e,. F_p. While some SNARKs are associated with a single field F_p, there are many SNARKs that rely on elliptic curve groups for security. For such SNARKs, the scalar field of the elliptic curve is F_p, and the base field is a different field F_q. Thus, the aim is to benchmark the field F_p, along with the field F_q and the elliptic curve group (if applicable). Benchmarking F_p and F_q involves benchmarking the following operations:

- Addition
- Subtraction
- Multiplication
- (Modular) Exponentiation
- Inverse Exponentiation

An elliptic curve is defined over a prime field of specific order (F_q). The elliptic curve group (E(F_q)) consists of the subgroup of points in the field that are on the curve, including a special point at infinity. While some SNARKs operate over elliptic curves without requiring pairings, others require pairings and therefore demand for pairing-friendly elliptic curves. The pairing operation takes an element from G_1 and an element from G_2 and computes an element in G_T. The elements of G_T are typically denoted by e(P, Q), where P is an element of G_1 and Q is an element of G_2. For efficiency, it is required that not only is the finite field arithmetic fast, but also the arithmetic in groups G_1 and G_2 as well as pairings are efficient. Therefore, we intend to benchmark the following operations over pairing-friendly elliptic curves:

- Scalar Multiplication
  - in G for single elliptic curves
  - in G_1 and G_2 for pairing-friendly elliptic curves
- Multi-Scalar Multiplication (MSM)
  - in G for single elliptic curves
  - in G_1 and G_2 for pairing-friendly elliptic curves
- Parings
  - for pairing-friendly elliptic curves

### Circuits

Many end-to-end applications require proving a specific cryptographic primitive,
which requires the specification of said cryptographic primitive in a specific ZKP
framework.

- *Circuits for native field operations* -  These operations, namely, addition and multiplication in F_p, are supported by each SNARK library, and they are the most efficient to prove with a SNARK because arithmetic modulo F_p is the native computation model of a SNARK. This provides a good understanding of the efficiency of the core
SNARK implementation.
- *Circuits for non-native field operations* - All computations we want to prove do not belong to arithmetic modulo p. For instance, Z_{2^64} or uint64/int64 is a popular data type in traditional programming languages. Or, we might want to prove arithmetic on a different field, say Z_q. This usually happens when we want to verify elliptic-curve based cryptographic primitives. An example of this is supporting verification of ECDSA signatures. The native field of elliptic curve underlying the chosen SNARK typically differs from the base field of the secp256k1 curve
- *Circuits for SNARK-optimized primitives* - One of the challenges in the practically using SNARKs is their inefficiency with regard to traditional hash algorithms, like SHA-2, and traditional signature algorithms, such as ECDSA. They are fast when executed on a CPU, but prohibitively slow when used in a SNARK. As a result, the community has proposed several hash functions and signature algorithms that are SNARK-friendly, such as the following:
  - Poseidon Hash
  - Pedersen Hash
  - MIMC Hash
  - Ed25519 (EdDSA signature)
- *Circuits for CPU-optimized primitives* - Even though it would be beneficial to only rely on SNARK optimized primitives, practical applications often don’t allow for the usage of e.g. Poseidon hash functions or SNARK friendly signature schemes. For example, verifying ECDSA signatures in SNARKs is crucial when building e.g. zkBridge, however an implementation requires for non-native field arithmetic, and therefore yields many constraints. Similarly, for building applications such as TLS Notary, one has to prove SHA-256 hash functions and AES-128 encryption which yields many constraints. Hence, we intend to benchmark the performance of the following cryptographic primitives and their
circuit implementations in different ZKP-frameworks:
  - SHA-256
  - Blake2
  - ECDSA

### Current Features

On a high level, zk-Harness takes as input a configuration file. The “Config Reader” reads the standardized config and invokes the ZKP framework as specified in the configuration file. You can find a description of the configuration file in the tutorials/config sub-folder of the GitHub repository. Each integrated ZKP framework exposes a set of functions that take as an input the standardized configuration parameters to execute the corresponding benchmarks. The output of benchmarking a given ZKP framework is a log file in csv format with standardized metrics. The log file is read by the “Log Analyzer”, which compiles the logs into pandas dataframes that are used by the front-end and displayed on the public website. You can find the standardized logging format in the tutorials/logs sub-folder.

Currently, zk-Harness includes the following components as a starting point:

- Benchmarks for field arithmetic
- Benchmarks for Elliptic curve group operations
- Benchmarks for circuit implementations
- In the following frameworks:
  - gnark
  - circom

We aim to successively expand this list to further include benchmarks for other ZKP frameworks, recursion and zk-EVMs. As a part of the ZKP/Web3 Hackathon hosted by UC Berkeley RDI, we aim to further develop the frameworks integrated into zk-Harness. You can find the program description detailing future integrations [here](https://drive.google.com/file/d/1Igm47dFXSOFAC_wldfUG4Y9OiITqlbQu/view). A detailed list of currently included sub-components and the framework architecture can be found in the [GitHub](https://github.com/zkCollective/zk-Harness) repository.

## How to use

Run any of the various targets in the `Makefile`.

To obtain a Nix environment in which you can successfully run these benchmarks,
first [install Nix](https://nixos.org/download.html)
and [enable flakes](https://nixos.wiki/wiki/Flakes#Enable_flakes),
then run `nix develop`.


## How to contribute

There are many ways in which you can contribute to the zk-Harness:

- Add benchmarks for circuits in an already integrated framework
- Integrate a new framework into the zk-Harness
- Propose new benchmark categories, such as for recursion and zk-EVMs.

zk-Harness is developed as part of the [zk-Hackathon](https://rdi.berkeley.edu/zkp-web3-hackathon/) hosted by the [Berkeley Center for Responsible Decentralized Intelligence](https://rdi.berkeley.edu/).
Further, zk-Harness is part of the [zk-Collective](https://github.com/zkCollective/).

### How to add new frameworks to the Nix environment

The `./flake.nix` file contains the Nix code which constructs the Nix environment. To integrate a new
framework into the Nix environment, this code must be changed so that it adds to the environment all
of the dependencies required to run the new framework-specific benchmark commands in the Makefile.
In general, the following steps will be involved:

 * Nix the framework, if it has not already been Nixed.
 * Add the Nixed framework to the flake `inputs`.
 * Add the framework output to the `devShells.default.packages` in the flake `outputs`.
 * Add any other dependencies needed to run the benchmarks to the `devShells.default.packages` in the flake `outputs`.
 * Add any necessary setup code to the `devShells.default.shellHook`.

A few things to note:

 * Nix packages are supposed to have reproducible builds, and there are restrictions within the Nix build system designed to rule out common causes of non-reproducibility. Therefore, a build process designed to work outside Nix may not work within Nix without modification.
    * In particular, the FHS (Filesystem Hierarchy Standard) does not apply in Nix; for instance, `/usr/bin/bash` cannot be expected to point to anything in the Nix environment. Instead of `/usr/bin/foo`, use `/usr/bin/env foo`.
 * Nix is a Turing complete programming language, and there cannot be a set recipe for Nixing software, since every software package is different.
 * Although Nix builds are in theory reproducible, the way this benchmark suite is set up does not run the benchmarks within Nix. In theory, running the benchmarks within the Nix environment should produce similar results on similar machines. However, since the benchmarks are not being run in a hermetic Nix environment, the benchmarks may work on your machine and not work at all on another machine of the same architecture, due to differences in software configuration. To avoid this, you should make sure that the commands in your Makefile only invoke dependencies that are in the Nix environment.

Here are some hopefully useful pointers for Nixing pre-existing code bases:

 * [cargo2nix](https://github.com/cargo2nix/cargo2nix) for Nixing Rust projects;
 * [Nixing Go projects](https://nixos.wiki/wiki/Go);
 * The Nix environment for this project contains Python and Node.js, and you can incorporate
   Node and Python dependencies by adding them to the `./requirements.txt` or `./package.json`
   file (as applicable).
 * You can learn more about Nix by reading the [Nix pills](https://nixos.org/guides/nix-pills/)
   and referring to the [Nix reference manual](https://nixos.org/manual/nix/stable/).
