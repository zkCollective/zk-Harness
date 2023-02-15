# zk-Harness - Benchmarking Hackathon

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

## Program Description

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

Develop a circuit implementation that does not yet exist and benchmark it against other implementations.

#### Task Description / Steps Involved

Given a framework, choose a cryptographic primitive that is not yet implemented, implement it and benchmark your implementation as compared to existing implementations in other frameworks. You can find a list of already implemented primitives here. The description on how to contribute a novel circuit implementation can be found in the respective folder of the framework.

This task comprises the following steps:

1. Choose a primitive to implement in a given framework
2. Read the tutorial on “How to contribute a new circuit” in the chosen framework
3. If the library is not yet supported, follow steps X-Y of #4
4. Create Tests for your circuit

#### Designated Tasks

Comparative Evaluation / implementation of one of the following primitives:

- SHA-256
- Blake 2
- AES-128
- ECDSA (see here)

in the following frameworks:

- circom
- gnark
- halo2

Prize: X$

### Supporting New Libraries

#### Goal

Integrate a framework into the zk-Harness for benchmarking.

#### Task Description / Steps Involved

There are plenty of implementations of SNARKs that we intend to include in the benchmarking framework. Hence, we encourage hackathon participants to integrate novel libraries that are not yet supported to support a holistic comparison of heterogeneous SNARK implementations.

This task comprises the following steps:

1. Choose a framework to implement (e.g. artworks / plonky2 / halo2)2. Support the data loading of configuration files
3. Configure the framework behavior based on the configuration file
4. Generate logs for a specified logging format (see logging formats here)
5. Integrate the logs in the frontend implementation
6. Create a pull request to integrate your framework in the harness and display the evaluation results on the public website.

You can find the detailed description on how to add a new library here. You can find the standardized, cross framework, log format which is consumed by the log analyzer here and the description of the generic config files here. To fully integrate a framework, you’ll need to adapt the config reader to invoke your benchmarking script and adapt the log analyzer to display your benchmarks on the webpage.

#### Designated Tasks

Integrate one or more of the following libraries into the zk-Harness:

- plonky2
- jellyfish
- arkworks

### Benchmarking Recursion

#### Goal

Benchmark implementations of recursive proofs.

#### Task Description / Steps Involved

Commonly, proof recursion can be achieved through the following approaches:

1. Encoding arithmetic statements of a field into equations over a different field (described above as non-native field arithmetic, e.g. here)
2. 2-chains and cycles of elliptic curves - use matching base fields to implement one curve inside the other

In this task, you should benchmark common implementations of recursion in popular frameworks.

#### Designated Tasks

Comparative Evaluation of recursion as implemented in the following frameworks:

- plonky2
- halo2
- Nova
- gnark (BW6 on BLS12-381 - see here)

Prize: X$

## Hackthon Awards and Prize

In the hackathon, there are two types of projects: self-selected projects and designated tasks. Teams can come up with their own creative project ideas and independently submit contributions beyond the ones specified in the sections “Recommended Starting Points” in each of the challenges above. Each category has its own prize, which is detailed below.
In the case of designated tasks, specific tasks critical to the ZK bridge have been identified. Detailed instructions and awards are provided for these tasks.

### Self-selected project

1. Grand Prize: The team that develops the best overall project, as evaluated by the panel of judges, will receive a grand prize of $10,000.

2. Runner-Up Prize: The team that develops the second-best project, as evaluated by the panel of judges, will receive a runner-up prize of $5,000.

3. Most Creative Idea: The team that develops the most creative idea, as evaluated by the panel of judges, will receive a prize of $2,000.

### Designated Tasks

Each designated task is assigned a prize value, based on the estimated difficulty and work effort.

## Review Mechanism

We will carefully review the correctness of benchmarks to integrate in the benchmarking framework. On completing a novel benchmark that is not yet integrated in the zk-Harness, we recommend that you create a pull request that can be independently reviewed.
Your implementation will be evaluated based on the following criteria:

1. Correctness of the implementation.
2. Efficiency of the implementation.
3. Quality of the documentation.
