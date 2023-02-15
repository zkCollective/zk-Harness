 ZK Compilers

We cordially invite the zk SNARK community to join us in creating a comprehensive benchmarking framework for zk SNARKs. As part of our efforts to further advance the technology and promote its widespread adoption, we have organized a Hackathon to bring together experts and enthusiasts from the community to collaborate and contribute to the establishment of a standardized benchmarking framework. This is a crucial step in our mission to create a reference point for non-experts and experts alike on what zkSNARK scheme best suits their needs, and to also promote further research by identifying performance gaps. We believe that the collective efforts of the community will help to achieve this goal. Whether you are a researcher, developer, or simply passionate about zk SNARKs, we welcome your participation and contribution in this exciting initiative.

Google drive link: <https://drive.google.com/drive/u/0/folders/1zkiGrN1xA4FIfAk4N3P8v-\_QIeA8P4Jb>

![Alt text](/HarnessSpecification.png?raw=true "Title")

## Introduction

There is a large and ever-increasing number of SNARK implementations. Although the theoretical complexity of the underlying proof systems is well understood, the concrete costs rely on a number of factors that depend on the underlying computation model and its compatibility with the specific application. The difference in implementation is evident in multiple layers of the SNARK stack, summarized as follows:

### Field Arithmetic

- Native Field Arithmetic

  Proof systems encode a polynomial equation over a finite field F_r, where r is generally a prime number. In order to perform arithmetic operations in finite fields, it is often necessary to work with large integers, which can have hundreds or thousands of digits. SNARKs rely on the use of finite fields for their construction, and native field arithmetic is used to perform arithmetic operations on the elements of the field (for a reference implementation - see here). Benchmarking native field arithmetic involves benchmarking the following operations:

      - Addition
      - Subtraction
      - Multiplication
      - Division
      - (Modular) Exponentiation
      - Inverse Exponentiation

- Non-Native Field Arithmetic

  Encoding programs that are not designed to operate over F_r is generally expensive. This problem is generally solved by simulating the operations in the target field in the base field F_r (for a reference on how this is done in practice, see here). Benchmarking non-native field arithmetic involves benchmarking the same operations as for native field arithmetic.

### Operations over Elliptic Curves

- Scalar Multiplication

  Scalar multiplication involves multiplying a single point on an elliptic curve by a scalar value, which is typically an integer. The result is another point on the curve. Scalar multiplication is a fundamental operation in elliptic curve cryptography, and is used in many cryptographic protocols such as key exchange, digital signatures, and encryption.

- Multi-Scalar Multiplication (MSM)

  Large MSMs over points on an elliptic curve are required for instance in the Setup and Prove algorithms of common proving systems. For example, Groth16 requires (3n G_1 + m G_2) MSMs in the per-circuit setup phase, where n is the number of multiplication gates and m is the number of wires. Therefore, an efficient implementation of MSM is an important feature in general purpose SNARK framework. If you are unfamiliar with the problem of Multi-Scalar Multiplication, you can find a good introduction here.

- Pairings

  Pairing operations over pairing-friendly elliptic curves are essential in the verification algorithm of pairing-based SNARKs, such as PLONK and Groth16. For example, the verification algorithm of PLONK with KZG polynomial commitments requires computing 2 pairings.

### Circuits

- Circuits for SNARK-optimized primitives

- Circuits for CPU-optimized primitives

  Even though it would be beneficial to only rely on SNARK optimized primitives, practical applications often don’t allow for the usage of e.g. Poseidon hash functions or SNARK friendly signature schemes. For example, verifying ECDSA signatures in SNARKs is crucial when building e.g. zkOracles, however an implementation requires for non-native field arithmetic and  in-circuit (see #2 and here), and therefore yields many constraints. Similarly, for building applications such as TLS Notary, one has to prove SHA-256 hash functions and AES-128 encryption which yields many constraints.

We integrated gnark to exemplify how to integrate libraries into the benchmarking harness. You can find a description on how to run benchmarks for gnark here.

## Review Mechanism

We will carefully review the correctness of benchmarks to integrate in the benchmarking framework. On completing a novel benchmark that is not yet integrated in the zk-Harness, we recommend that you create a pull request that can be independently reviewed.
Your implementation will be evaluated based on the following criteria:

1. Correctness of the implementation.
2. Efficiency of the implementation.
3. Quality of the documentation.

## Program Desctiption

### Benchmarking Mathematical Operations

#### Goal

Benchmark framework specific implementations of native / non-native field arithmetic and elliptic curve group operations

#### Task Description / Steps Involved

The purpose of this task is to benchmark the performance of implementation of field arithmetic and elliptic curve group operations, in order to assess their efficiency and identify areas for improvement.

Benchmarking native & non-native field arithmetic involves benchmarking the following operations:

- Addition
- Subtraction
- Multiplication
- Division
- (Modular) Exponentiation
- Inverse Exponentiation

Benchmarking elliptic curve group operations involves benchmarking the following operations:

- Scalar Multiplication
- Multi-Scalar Multiplication (MSM)
- Pairing

#### Designated Tasks

Comparative Evaluation of native field arithmetic in the following frameworks:

- circom
- gnark
- halo2

Prize: X$

Comparative Evaluation of non-native field arithmetic in the following frameworks:

- circom
- gnark
- halo2

Prize: X$

Comparative Evaluation of the above mentioned operations over all curves as supported by the following frameworks:

- circom
- gnark
- halo2

Prize: X$

### Benchmarking Circuit Implementations

#### Goal

#### Task Description / Steps Involved

#### Designated Tasks

### Supporting New Libraries

#### Goal

#### Task Description / Steps Involved

#### Designated Tasks

### Benchmarking Recursion

#### Goal

#### Task Description / Steps Involved

#### Designated Tasks