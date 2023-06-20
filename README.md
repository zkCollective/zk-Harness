<h1 align="center">zk-Harness</h1>

zk-Harness is a benchmarking framework for *zero knowledge succinct non-interactive arguments (zkSNARKs)*. 
This repository contains a modular and easily extensible framework for benchmarking zkSNARKs and underlying mathematical primitives.

## 📊 Benchmark Results

The benchmark results are currently hosted at [zk-bench.org](https://www.zk-bench.org).

The backend for mathematical operations has been merged with [zka.lc](https://zka.lc/).

### Currently Supported Standard Payloads 

The current framework supports a set of payloads for each library.
We aim to successively extend the following once more circuits are available as `std` in the respective libraries.

|          | Exponentiate       | SHA-256           |
| -------- | ------------------ | ----------------- |
| `Bellman`| :heavy_check_mark: | :heavy_check_mark:|
| `Circom` | :heavy_check_mark: | :heavy_check_mark:|
| `Gnark`  | :heavy_check_mark: | :heavy_check_mark:|
| `Halo2`  | :heavy_check_mark: | :heavy_check_mark:|
| `Starky` | :heavy_check_mark: | :heavy_check_mark:|

### Backends, Curves, Fields and Arithmetizations

<!-- TODO -->

### Run Benchmarks On Your Own!

zk-Harness is supposed to be easily extensible and modular, which means that you should be able to integrate you own circuits with ease.
Each framework in `framework/<framework_name>` includes a detailed description on how to add a self-developed circuit that goes beyond the standard payloads already integrated.

To validate and re-run any of the standard payloads, please refer to the the `Makefile`, which utilizes standard configuration files.

To obtain a Nix environment in which you can successfully run these benchmarks, first [install Nix](https://nixos.org/download.html) and [enable flakes](https://nixos.wiki/wiki/Flakes#Enable_flakes), then run `nix develop`.

When running into issues for nix on M1/M2 Macs, please refer to [this issue](https://github.com/input-output-hk/plutus-pioneer-program/issues/40).  


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
