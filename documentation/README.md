# Documentation for the zk-Harness benchmarking framework

Our benchmarking framework is designed to support all types of different frameworks.
We specify a generic set of interfaces, such that benchmarks can be invoked through a configuration file `config.json`, which produces a standardized output for a given benchmarking scenario.
You can find a description of the configuration file in the `config` sub-folder of this folder, whereas the logging format can be found in the `logging` sub-folder of this repository.

## Adding a new framework to the zk-Harness

```diff
- TODO - Fix this once fully implemented - So far only first draft
```

To integrate a new framework, one should follow the following steps:

1. First, fork the `zk-benchmarks` repository
2. Create a `./<framework_name>` folder in the root folder of the repository. This directory should not have a particular structure but it should:
    1. Have a `README` file describing how to install all required dependencies and run an example circuit
    2. Provide script(s) for running arithmetic, elliptic curve, and circuit benchmarks (see step 3 below)
3. Create a custom benchmarking script that *(i)* reads from the standardized input of the ``config.json`` as described in the ``config`` folder and outputs *(ii)* the standardized logs as described in the ``logging`` folder.
   1. For example, benchmarking for `gnark` is done through a custom CLI, based on [cobra](https://github.com/spf13/cobra)
   2. In contrast, benchmarking for Circom/snarkjs is performed using a bash script (for circuits) and a JavaScript script (for arithmetics and elliptict curves). 
4. Modify scripts `process_arithmetic.py`, `process_ec.py`, and `process_circuit.py` in `_scripts/reader` directory. Specifically, you should add a function named `build_command_project_name` that given a `Payload` object (read from the config) produces commands to be executed (using the script described in step 3) that will run the benchmarks.
5. Create three configuration files in `_input/config/project_name` for running the default benchmarks for arithmetic, elliptic curve, and circuit.
6. Add rules for running the benchmarks in the main Makefile.
7. Create a documentation in the `<framework_name>/TUTORIAL.md` directory to outline how others can include new benchmarks in the framework. Depending on the project you pursue, it should contain documentation on either, or all, of the following:
   1. How to add a new circuit implementation
   2. How to run tests for integrated circuits

If you follow the specified interfaces for config and logging, your framework specific benchmarking should seamlessly integrate into the zk-Harness.

Once finished, please create a Pull Request and assign it to one of the maintainers for review.
