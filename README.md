# ZK Compilers

This is a software to benchmark various implementations of polynomial commitment schemes, curve implementations, circuit implementations that leverage different proving schemes.

Google drive link: <https://drive.google.com/drive/u/0/folders/1zkiGrN1xA4FIfAk4N3P8v-\_QIeA8P4Jb>

<!--ts-->
- [ZK Compilers](#zk-compilers)
  - [Benchmarks](#benchmarks)
    - [Toy Examples](#toy-examples)
  - [Overview of ZK SNARK Compilers](#overview-of-zk-snark-compilers)
  - [Overview of ZK STARK Compilers](#overview-of-zk-stark-compilers)
  - [Overview of Compilers for Dedicated Verifier Proofs](#overview-of-compilers-for-dedicated-verifier-proofs)
  - [General Purpose Frameworks](#general-purpose-frameworks)
  - [Implementations leveraging different ZK Compilers](#implementations-leveraging-different-zk-compilers)
    - [PLONK](#plonk)
    - [Halo2](#halo2)
    - [Circom](#circom)
    - [Bulletproofs](#bulletproofs)
    - [Sigma Protocols](#sigma-protocols)
  - [Polynomial Commitments](#polynomial-commitments)
    - [Additional Info on Polynomial Commitments and Multilinear Extenstions](#additional-info-on-polynomial-commitments-and-multilinear-extenstions)
  - [Curves and Pairings](#curves-and-pairings)
  - [Hardware Accelerations](#hardware-accelerations)
    - [GPU-based Hardware acceleration](#gpu-based-hardware-acceleration)
    - [FPGA-based Hardware acceleration](#fpga-based-hardware-acceleration)
  - [Trusted Setup](#trusted-setup)
  - [Toy Examples](#toy-examples-1)
  - [Cryptographic Primitives](#cryptographic-primitives)
    - [Optimized for SNARKs](#optimized-for-snarks)
    - [Traditional Schemes](#traditional-schemes)
  - [Existing Performance Comparisons](#existing-performance-comparisons)
  - [What are questions that we should answer in the paper?](#what-are-questions-that-we-should-answer-in-the-paper)
  - [Questions for us to resolve in order to determine the focus](#questions-for-us-to-resolve-in-order-to-determine-the-focus)
    - [Questions to focus on](#questions-to-focus-on)
    - [Micro-Benchmark Metrics](#micro-benchmark-metrics)

<!-- Created by https://github.com/ekalinin/github-markdown-toc -->
<!-- Added by: serious, at: Mon Feb  6 13:44:17 CET 2023 -->

<!--te-->

## Benchmarks

TODO:

- [ ] Include Toy Example multiplication Golang
- [ ] Output benchmarks in a csv file

Benchmarks can be run through the Makefile. Currently, the following benchmarks are supported:

### Toy Examples

The benchmarks for Toy examples can be run by executing ``` make benchmarks-toy ```
The following Toy example benchmarks will be executed:

- Exponentiation
  - Gnark

## Overview of ZK SNARK Compilers

In the following, we initially separate between *Succinct Non-Interactive Arguments of Knowledge* (SNARKs), *Succinct Transparent Arguments of Knowledge* (STARKs) and *Dedicated Verifier Proofs*, which may require interaction with a dedicated verifier and therefore lose non-interactivity.

| Name                                                                  | Lang.     | Arith.  | IOP     | Front/Back    |
| :---                                                                  | :----:    | :----:  | :---:    |   ---:   |
| [libsnark](https://github.com/scipr-lab/libsnark)                     | C++       | R1CS    | Groth16 | Back          |
| [Bellman](https://github.com/zkcrypto/bellman)                        | Rust      | R1CS    | Groth16 | Back          |
| [jsnark](https://github.com/akosba/jsnark)                            | JS        | libsnark    |  libsnark  | Front          |
| [snarky](https://github.com/o1-labs/snarky)                           | OCAML     | libsnark    |  libsnark  | Front          |
| [gnark](https://github.com/ConsenSys/gnark)                           | Go     | Plonk    |  Plonk  | Front + Back          |
| [PLONK Dusk Network](https://github.com/dusk-network/plonk)           | Rust     | Plonk    |  Plonk  | Front + Back          |
| [circom](https://github.com/iden3/circom)                             | Rust     | R1CS    |  Groth16  | Front + Back          |
| [arkworks](https://github.com/arkworks-rs)                            | Rust     | R1CS    |  Groth16/Spartan/Marlin  | Front + Back          |
| [jellyfish](https://github.com/EspressoSystems/jellyfish)             | Rust     | Plonk    |  Plonk  | Front + Back          |
| [halo2](https://github.com/zcash/halo2)                               | Rust     | UltraPlonk    |  Halo2  | Front + Back          |
| [adjoint-io bulletproofs](https://github.com/sdiehl/bulletproofs)     | Rust     | - (DLOG)    |  Bulletproof  | Front + Back          |
| [DIZK](https://github.com/scipr-lab/dizk)                             | Java     | R1CS    |  Groth16  | Front + Back          |
| [Spartan](https://github.com/microsoft/Spartan)                       | Rust     | R1CS    |  Spartan  | Front + Back          |
| [Anoma Vamp-IR](https://github.com/anoma/vamp-ir)                     | Rust     | ?    |  X  | Front          |
| [VIRGO](https://github.com/sunblaze-ucb/Virgo)                        | C++     | ?    |  ?  | ?          |
| [Hyperplonk](https://github.com/EspressoSystems/hyperplonk)           | Rust      | Plonk    |  ?  | ?          |

## Overview of ZK STARK Compilers

| Name                                                          | ZK Type       | Arithmetization       | Commitment Scheme         |
| :---                                                          |    :----:     |          ---:         |          ---:             |
| [libstark](https://github.com/elibensasson/libSTARK)          | Title         | Here's this           |                           |
| [OpenZKP](https://github.com/0xProject/OpenZKP)               | Text          | Here's this           |                           |
| [genSTARK](https://github.com/GuildOfWeavers/genSTARK)        | Text          | Here's this           |                           |
| [Hodor](https://github.com/matter-labs/hodor)                 | Text          | Here's this           |                           |
| [Winterfell](https://github.com/facebook/winterfell)          | Text          | Here's this           |                           |
| [ethSTARK](https://github.com/starkware-libs/ethSTARK)        | Text          | Here's this           |                           |

## Overview of Compilers for Dedicated Verifier Proofs

| Name                                                      | ZK Type       | Arithmetization       | Commitment Scheme         |
| :---                                                      |    :----:     |          ---:         |          ---:             |
| [emp-zk](https://github.com/emp-toolkit/emp-zk)           | Text          | Here's this           |                           |
| [FETA](https://github.com/KULeuven-COSIC/Feta)            | Text          | Here's this           |                           |

## General Purpose Frameworks

We should have a discussion section describing frontend approach implementing
circuits that essentially execute step-by-step some simple CPU, also called a 
virtual machine (VM). Such approaches are: StarkWare's (Cairo), zkEVMS,
Polygon's VM, RISC-V (?), and zkLLVM.

## Implementations leveraging different ZK Compilers

### PLONK

Plonk is commonly used in:

- Plonk Aztex 2.0 Monorepo [here](https://github.com/neidis/plonk-plookup)
- Barretenberg (Aztec), an optimized elliptic curve library for the bn128 curve, and PLONK SNARK prover [here](https://github.com/AztecProtocol/barretenberg)
- Sec-Bit implementation of Plonk & Plookup with arkworks libraries [here](https://github.com/sec-bit/plonk) 
- zk Garage PLONK compatible with the arkworks library [here](https://github.com/ZK-Garage/plonk)
- Plonky - Recursive arguments based on PLONK and Halo [archived here](https://github.com/mir-protocol/plonky)
  - Documentation for Plonky can be found [here](https://mirprotocol.org/blog/Fast-recursive-arguments-based-on-Plonk-and-Halo)
- Plonky2 - Recursive arguments based on PLONK and FRI polynomial commitments [here](https://github.com/mir-protocol/plonky2)
  - Writeup documenting details about Plonky2 can be found [here](https://github.com/mir-protocol/plonky2/blob/main/plonky2/plonky2.pdf)

Recent research work on PLONK:

- HyperPlonk - Plonk with multilinear polynomial commitments [here](https://github.com/darkrenaissance/darkfi)

### Halo2

- combines efficient accumulation scheme with PLONKish arithmetization and needs no trusted setup
- based on IPA commitment scheme
- flourishing developer ecosystem
- prover time: O(N*log N)
- verifier time: O(1)>Groth16
- proof size: O(log N)
- assumption: discrete log

The arithemtization for Halo2 comes from PLONK, more precisely UltraPlonk, which is an extended version of Plonk that supports both lookup arguments and custom gates.

Detailed Informations about Halo2:

- The Halo2 Book [here](https://zcash.github.io/halo2/)

Halo2 is commonly used in:

- [zk-EVM circuits Privacy Scaling Explorations](https://github.com/privacy-scaling-explorations/zkevm-circuits)
- Privacy Scaling Explorations - halo2 with KZG instead of IPA [here](https://github.com/privacy-scaling-explorations/halo2wrong)
- halo2 with FRI instead of IPA [here](https://github.com/Orbis-Tertius/halo2)
- Scroll zk-EVM recursive aggregation [here](https://github.com/scroll-tech/halo2-snark-aggregator)
  - This library leverages recursive aggregation for constructing a zkEVM by using Halo2 with KZG polynomial commitments
  - At the time of writing, Halo2 does not yet support recursive aggregation
- Orbis zkEVM on Cardano [here](https://github.com/Orbis-Tertius)
- DarkFi Layer1

### Circom

- [circomlib](https://github.com/iden3/circomlib/tree/master/circuits) provides loads or circuit templates for use in applications (see detailed list in the sheets)
- Sismo Hydra S1 ZKPS [here](https://github.com/sismo-core/hydra-s1-zkps)
 
### Bulletproofs

### Sigma Protocols

Sigma protocols are commonly used in:

- Signal for group chats [see details here](https://signal.org/blog/signal-private-group-system/)
- Proving the equality of discrete logarithms & Proofs of Knowledge of discrete log (see the paper by Dan Boneh [here](https://eprint.iacr.org/2018/1188.pdf)

## Polynomial Commitments

Example of an univariate polynomial:    p = x^4 - 4*x^2 - 7*x + 9

Example of a multivariate polynomial:   p = 2*xz^4 - 4*xz^2 - 7*xy

Example of a multilinear polynomial:    p = 2*xz - 4*xz - 7*xy

- Hyrax ([Paper](https://eprint.iacr.org/2017/1132.pdf) / [Code](https://github.com/TAMUCrypto/hyrax-bls12-381))
- KZG ([Paper]())
  - Prima One - Simple KZG implementation [here](https://github.com/proxima-one/kzg)
  - Espresso Systems - Implementation of Multilinear KZG commitments [here](https://github.com/EspressoSystems/hyperplonk)
- Arkworks Polynomial Commitment implementations [here](https://github.com/arkworks-rs/poly-commit/tree/master/src) (KZG, IPA, Marlin KZG, Sonic KZG)
- FRI (Risc0, Plonky2, Winterfell) [Paper](Fast Reed-Solomon Interactive Oracle Proofs of Proximity)
  - Risc0 [implementation](https://github.com/risc0/risc0/blob/main/risc0/zkp/src/prove/fri.rs) / [Description](https://www.risczero.com/docs/reference-docs/about-fri) of FRI

### Additional Info on Polynomial Commitments and Multilinear Extenstions

- Interactive Proofs & Arguments, Low-Degree & Multilinear Extensions - Justin Thaler [here](https://people.cs.georgetown.edu/jthaler/IPsandextensions.pdf)

## Curves and Pairings

<!-- https://github.com/supranational/blst BLS12-381 -->
<!-- 2-adicity of curves see https://www.cryptologie.net/article/559/whats-two-adicity/ -->

## Hardware Accelerations

### GPU-based Hardware acceleration

- Mina GPU Groth16 accelerated prover [here](https://github.com/MinaProtocol/gpu-groth16-prover-3x)

### FPGA-based Hardware acceleration

<!-- TODO -->

## Trusted Setup

Pairing-based SNARKs require a trusted setup (universal/per-circuit) to achieve high efficiency.
The parameters generated during a trusted setup procedure have to remain secret to ensure the security of the overall system.
Commonly, trusted setup procedures are instantiated through MPC in a distributed manner.
Whereas e.g. Groth16 demands for a trusted setup per circuit, e.g. PLONK requires a universal trusted setup ceremony to instantiate the 

- Aleo Setup procedure [here](https://github.com/AleoHQ/aleo-setup)

## Toy Examples

Potential Toy Examples that we can evaluate in the paper:

- Multiplication
- Cubic polynomial
- Dot product
- Inner product
- Sorting
- Modulo
- Exponentiation

## Cryptographic Primitives

### Optimized for SNARKs

- MIMC hash function [here](https://byt3bit.github.io/primesym/mimc/)
- Poseidon hash function []() 
  - Scroll implementation can be found [here](https://github.com/scroll-tech/poseidon-circuit))
  - Arkworks implementation can be found [here](https://github.com/arkworks-rs/crypto-primitives/tree/main/src/sponge)

### Traditional Schemes

- Merkle Patricia Tree [docs](https://ethereum.org/en/developers/docs/data-structures-and-encoding/patricia-merkle-trie/) 
  - Implemenation by Scroll can be found [here](https://github.com/scroll-tech/mpt-circuit)
- SHA-256 [RFC](https://www.rfc-editor.org/rfc/rfc6234)
  - Arkworks implementation can be found [here](https://github.com/arkworks-rs/crypto-primitives/tree/main/src/crh)
  - Risc0 implementation can be found [here](https://github.com/risc0/risc0-rust-examples/tree/main/sha)

## Existing Performance Comparisons

- Comparison pairings in gnark/arkworks [here](https://eprint.iacr.org/2022/1162.pdf)
- Polynomial Commitment benchmark can be found [here](https://2π.com/23/pc-bench/index.html)
- SHA-256 comparison for halo2, plonky2, circom [here](https://github.com/Sladuca/sha256-prover-comparison)
- Master Thesis on a theoretical / practical Overview of SNARKs [here](https://is.muni.cz/th/ovl3c/SNARKs_STARKs_introduction.pdf)
- ZKP Test Repository [here](https://github.com/spalladino/zkp-tests)

## What are questions that we should answer in the paper?

- Which compilers are mostly used?
- Which compiler should a developer choose given a specific project in mind?
  - E.g. Recursion
- What is the state-of-the-art in current ZKP implementations?
  - Leave hands of off this
- As a designer of a compiler, what are limitations of existing approaches?
  - Leave as discussion
- As policy makers - what are the assurances that ZKPs can provide - and where can they fail?
  - Leave as discussion
  - Search in Repos where they fixed any bugs

- What is the performance of different Curve implementations for SNARKs that apply bilinear pairings?
  - Can do that - very important primitive
- What is the performance of different polynomial commitment schemes? (For SNARKs that rely on polynomial commitment schemes)
- How do optimizations, such as lookups and custom gates, influence the performance of the baseline SNARK?
- What are current bottlenecks to SNARKs?
  - Multi Scalar Multiplication
  - Fast Fourier Transformation
  - Trusted Setup Ceremony

## Questions for us to resolve in order to determine the focus

### Questions to focus on

1. Performance Toy Examples
   1. Create a Table of availabel toy examples
   2. Decide which ones to compare
      1. circom
      2. gnark
2. What are the current bottlenecks in SNARKs?
3. Which are the most commonly use libraries / compilers / HDLs?
  
### Micro-Benchmark Metrics

- What are impor
