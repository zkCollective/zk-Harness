# ``config.json`` Reference

The ``<project>.json`` file contains information about a particular benchmarking project. The following describes how to run benchmarks given a specific config file. Further, we specify each of the keys in a common config file and their expected values.

## Running benchmarks for a specific config file

Given a config file, you can run ``make config=/path/to/config.json`` to run the evaluation given a specific config file.

## config.json key specification

### ``project``

The name of the project being benchmarked.

### ``project_url``

The URL(s) to the repository for the project.

### ``backend``

The backend algorithm to use for proving the specified circuit(s).

### ``curves``

The curves over which the backend algorithm should be run.

### ``circuits``

The name of the circuit to benchmark. The circuit name used as an input here should be equivalent to the name of the file stored in ``<framework>/circuits/X/<circuit_name>.<extension>``.
Equivalen circuits across frameworks should have the same naming scheme for ease of comparison.

### ``algorithm``

The algorithm to execute.
Valid algorithms to execute in a given framework are currently:

- ``compile``
- ``setup``
- ``prove``
- ``verify``

If a given algorithm is not specified for the configured framework, the execution of ``make config=/path/to/config.json`` will fail.