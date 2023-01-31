# ZK Compilers

This repository provides an overview of existing zk compilers and DSLs.
In the followingm, we initially separate between *Succinct Non-Interactive Arguments of Knowledge* (SNARKs), *Succinct Transparent Arguments of Knowledge* (STARKs) and *Dedicated Verifier Proofs*, which may require interaction with a dedicated verifier and therefore lose non-interactivity.

## Overview of ZK SNARK Compilers

| Name                                                                  | Language       | Arithmetization       | Commitment Scheme         |
| :---                                                                  |    :----:     |          ---:         |          ---:             |
| [libsnark](https://github.com/scipr-lab/libsnark)                     | C++         | Here's this           |                           |
| [Bellman](https://github.com/zkcrypto/bellman)                        | Text          | Here's this           |                           |
| [jsnark](https://github.com/akosba/jsnark)                            | Text          | Here's this           |                           |
| [snarky](https://github.com/o1-labs/snarky)                           | Text          | Here's this           |                           |
| [gnark](https://github.com/ConsenSys/gnark)                           | Text          | Here's this           |                           |
| [PLONK Dusk Network](https://github.com/dusk-network/plonk)           | Text          | Here's this           |                           |
| [circom](https://github.com/iden3/circom)                             | Text          | Here's this           |                           |
| [arkworks](https://github.com/arkworks-rs)                            | Text          | Here's this           |                           |
| [jellyfish](https://github.com/EspressoSystems/jellyfish)             | Text          | Here's this           |                           |
| [halo2](https://github.com/zcash/halo2)                               | Text          | Here's this           |                           |
| [adjoint-io bulletproofs](https://github.com/sdiehl/bulletproofs)     | Text          | Here's this           |                           |
| [DIZK](https://github.com/scipr-lab/dizk)                             | Text          | Here's this           |                           |
| [Spartan](https://github.com/microsoft/Spartan)                       | Text          | Here's this           |                           |
| [Anoma Vamp-IR](https://github.com/anoma/vamp-ir)                     | Text          | Here's this           |                           |
| [VIRGO](https://github.com/sunblaze-ucb/Virgo)                        | Text          | Here's this           |                           |
| [Hyperplonk](https://github.com/EspressoSystems/hyperplonk)           | Text          | Here's this           |                           |

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

- The Halo2 Book [here]()
- 

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
- Arkworks Polynomial Commitmente implementations [here](https://github.com/arkworks-rs/poly-commit/tree/master/src) (KZG, IPA, Marlin KZG, Sonic KZG)

### Additional Info on Polynomial Commitments and Multilinear Extenstions

- Interactive Proofs & Arguments, Low-Degree & Multilinear Extensions - Justin Thaler [here](https://people.cs.georgetown.edu/jthaler/IPsandextensions.pdf)

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

<!-- TODO -->

## Cryptographic Primitives

### Optimized for SNARKs

- MIMC hash function [here](https://byt3bit.github.io/primesym/mimc/)
- Poseidon hash function []() 
  - Scroll implementation can be found [here](https://github.com/scroll-tech/poseidon-circuit))
  - Arkworks implementation can be found [here](https://github.com/arkworks-rs/crypto-primitives/tree/main/src/sponge)

### Traditional Schemes

- Merkle Patricia Tree []() 
  - Implemenation by Scroll can be found [here](https://github.com/scroll-tech/mpt-circuit)
- SHA-256 []()
  - Arkworks implementation can be found [here]()

## Existing Performance Comparisons

- Comparison pairings in gnark/arkworks [here](https://eprint.iacr.org/2022/1162.pdf)
- Performance comparison of polynomial commitments []()
- SHA-256 comparison for halo2, plonky2, circom [here](https://github.com/Sladuca/sha256-prover-comparison)
- Master Thesis on a theoretical / practical Overview of SNARKs [here](https://is.muni.cz/th/ovl3c/SNARKs_STARKs_introduction.pdf)
