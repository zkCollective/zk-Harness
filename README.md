![Alt text](/logo_harness.png?raw=true "Title")


# zk-Harness - A benchmarking framework for general purpose Zero-Knowledge Proofs


## What is it?


zk-Harness is a benchmarking framework for general purpose zero-knowledge proofs.
It is designed to be modular - new circuit implementations and ZKP-frameworks can be easily added, without extensive developer overhead.
zk-Harness has a standardized set of interfaces for configuring benchmark jobs and formatting log outputs.
Once a new component is included, it's benchmarks will be displayed on [zk-harness.org](zk-harness.org).


**Note: zk-Harness is a WIP. Its architecture may change over time.**


## Main Features


zk-Harness currently includes the following:


- Benchmarks for field arithmetic
- Benchmarks for Elliptic curve group operations
- Benchmarks for circuit implementations
- In the following frameworks:
 - gnark
 - circom


A detailed list of included sub-components can be found in the respective subdirectory of the ZKP-framework.


## How to contribute


There are many ways in which you can contribute to the zk-Harness:


- Add benchmarks for circuits in an already integrated framework
- Integrate a new framework into the zk-Harness
- Propose new benchmark categories, such as for recursion and zk-EVMs.


zk-Harness is developed as part of the [zk-Hackathon](https://rdi.berkeley.edu/zkp-web3-hackathon/) hosted by the [Berkeley Center for Responsible Decentralized Intelligence](https://rdi.berkeley.edu/).
Further, zk-Harness is part of the [zk-Collective](https://github.com/zkCollective/).