![Alt text](/logo_harness.png?raw=true "Title")


# zk-Harness - A benchmarking framework for general purpose Zero-Knowledge Proofs

We cordially invite the zk SNARK community to join us in creating a comprehensive benchmarking framework (zk-Harness) for zk SNARKs. This is a crucial step in the important mission to create a reference point for non-experts and experts alike on what zkSNARK scheme best suits their needs, and to also promote further research by identifying performance gaps. We believe that the collective efforts of the community will help to achieve this goal. Whether you are a researcher, developer, or simply passionate about zk SNARKs, we welcome your participation and contribution in this exciting initiative.

It is designed to be modular - new circuit implementations and ZKP-frameworks can be easily added, without extensive developer overhead.
zk-Harness has a standardized set of interfaces for configuring benchmark jobs and formatting log outputs.
Once a new component is included, it's benchmarks will be displayed on [zk-bench.org](https://www.zk-bench.org).


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
