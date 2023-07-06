# Tutorials  for the zk-Harness benchmarking framework


Our benchmarking framework is designed to support all types of different frameworks.
We specify a generic set of interfaces, such that benchmarks can be invoked through a configuration file ``config.json``, which produces a standardized output for a given benchmarking scenario.
You can find a description of the configuration file in the ``config`` sub-folder of this folder, whereas the logging format can be found in the ``logging`` sub-folder of this repository.


![Alt text](./HarnessSpecification.jpg?raw=true "Title")

## Overview

On a high level, zk-Harness takes as input a configuration file. The “Config Reader” reads the standardized config and invokes the ZKP framework as specified in the configuration file. You can find a description of the configuration file in the tutorials/config sub-folder of the GitHub repository. Each integrated ZKP framework exposes a set of functions that take as an input the standardized configuration parameters to execute the corresponding benchmarks. The output of benchmarking a given ZKP framework is a log file in csv format with standardized metrics. The log file is read by the “Log Analyzer”, which compiles the logs into pandas dataframes that are used by the front-end and displayed on the public website. You can find the standardized logging format in the tutorials/logs sub-folder.


## Adding a new framework to the zk-Harness


To integrate a framework, one should follow the following steps:


1. First, fork the ``zk-benchmarks`` repository
2. Create a ``./<framework_name>`` folder in the root folder of the repository.
3. Create a custom benchmarking script that *(i)* reads from the standardized input of the ``config.json`` as described in the ``config`` folder and outputs *(ii)* the standardized logs as described in the ``logs`` folder.
    * For example, benchmarking for ``gnark`` is done through a custom CLI, based on [cobra](https://github.com/spf13/cobra)
    * Your script should be able to take a variety of arguments as specified in the config.json, such that benchmarks can be easily executed and extended. E.g., a common command in the gnark integration would be ``./gnark groth16 --circuit=sha256 --input=_input/circuit/sha256/input_3.json --curve=bn254``
4. Modify the ``_scripts/reader/X`` scripts to include your newly created script as described in step 3, which is called if the ``project`` field of the respective config contains the ``<framework_name>`` of your newly added ZKP framework.
  * The ``_scripts/reader/X`` processing python scripts are invoked by ``__main__.py`` based on the ``category`` (currently arithmetics, ec, circuit) field as specified in the config.
5. Create a documentation in the ``./<framework_name>/tutorials`` folder to outline how others can include new circuits / benchmarks for another ``category`` in the framework. Depending on the project you pursue, it should contain documentation on either, or all, of the following:
    * How to add a new circuit implementation
    * How to run tests for integrated circuits
    * How to benchmark a new ``category`` in the ZKP-framework
6. Add config files for running the benchmarks in `_input/config/` and add a make rule for the new framework in the Makefile.

If you follow the specified interfaces for config and logs, your framework specific benchmarking should seamlessly integrate into the zk-Harness frontend.

Once finished, please create a Pull Request and assign it to one of the maintainers for review and correct implementation.
