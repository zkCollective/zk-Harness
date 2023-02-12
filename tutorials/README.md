# Tutorials  for the zk-Harness benchmarking framework

Our benchmarking framework is designed to support all types of different frameworks.
We specify a generic set of interfaces, such that benchmarks can be invoked through a configuration file ``config.json``, which produces a standardized output for a given benchmarking scenario.
You can find a description of the configuration file in the ``config`` sub-folder of this folder, whereas the logging format can be found in the ``logging`` sub-folder of this repository.

## Adding a new framework to the zk-Harness

To integrate a custom, yet not supported framework, one should follow the following steps:

1. First, fork the ``zk-benchmarks`` repository
2. Create a ``./<framework_name>`` folder in the root folder of the repository. You can find the specifciation of the minimal framework folder structure [here]() (==TODO== - Add link & create docs).
3. Create a custom benchmarking script that *(i)* reads from the standardized input of the ``config.json`` as described in the ``config`` folder and outputs *(ii)* the standardized logs as described in the ``logging`` folder.
   1. For example, benchmarking for ``gnark`` is done through a custom CLI, based on [cobra](https://github.com/spf13/cobra)
4. Modify the ``config_reader.sh`` shell script to include the newly created ``bench_<framework_name>.sh``, which is called if the ``project`` field of the respective config contains the ``<framework_name>``.
5. Create a documentation in the ``./<framework_name>/tutorials`` folder to outline how others can include new benchmarks in the framework. Depending on the project you pursue, it should contain documentation on either, or all, of the following:
   1. How to add a new circuit implementation
   2. How to run tests for integrated circuits
   3. ==TODO== - What else should be documented

If you follow the specified interfaces for config and logging, your framework specific benchmarking should seemlessly integrate into the zk-Harness.

Once finished, please create a Pull Request and assign it to one of the maintainers for review and correct implementation.
