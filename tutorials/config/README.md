# ``config.json`` Reference

The ``config.json`` file contains information about what benchmarks to execute when invoking the zk-Harness. The following describes how to run benchmarks given a specific config file. Further, we specify each of the keys in a common config file and their expected values.

## Running benchmarks for a specific config file

Given a config file, you can run ``python3 -m _scripts.reader --config path/to/config.json`` to run the benchmarks defined by a given config file.

## config.json key specification

### ``project``

The name of the project being benchmarked.

### ``project_url``

The URL(s) to the repository for the project.

### ``category``

The category for which the zk-Harness should benchmark a ZKP-framework. Supports any of the following values:

- ``arithmetic``
  - Benchmarks finite field arithmetic
- ``curve_operation``
  - Benchmarks group operations
- ``circuit``
  - Benchmarks common circuit implementations of cryptographic primitives

### ``payload``

``payload`` specifies the exact algorithms to benchmark. The format of the payload depends on the mode the zk-Harness operates in.

#### ``payload`` specification for ``mode``==``arithmetic``

##### ``operation``

##### ``fields``

``fields`` specifies the fields that should be benchmarked. The current values are supported for the different ZKP-frameworks:

- ``gnark`` - Benchmarks are executed over the subgroup of prime order ``r`` of the based field F_p.
  - ``bn254``, ``bls12_381``, ``bls12_377``, ``bls24_315``, ``bw6_633``, ``bw6_761``

##### ``operations``

``operations`` specifies the operations that should be benchmarkes. The current values should be supported by a newly added ZKP-framework:

- ``Add`` (Addition)
- ``Sub`` (Subtraction)
- ``Mul`` (Multiplication)
- ``Div`` (Division)
- ``Exp`` (Modular Exponentiation)


#### ``payload`` specification for ``mode``==``curve_operation``

```diff
- TODO - Specify payload for arithmetic mode
```

#### ``payload`` specification for ``mode``==``circuit``

##### ``backend``

The backend algorithm(s) to use for proving the specified circuit(s).

##### ``curves``

The curve(s) for which the ZKP-framework should be benchmarked.

##### ``circuits``

The name of the circuit to benchmark. The circuit name used as an input here should be equivalent to the name of the file stored in ``<framework>/circuits/X/<circuit_name>.<extension>``.
Equivalent circuits across frameworks should have the same naming scheme for ease of comparison.

If a new circuit is added, which does not yet exist in any framework, one should create a new input specification in the ``/_input/<circuit_name>/input_<circuit_name>.json``.

##### ``algorithm``

The algorithm to execute.
Valid algorithms to execute in a given framework are currently:

- ``compile``
- ``setup``
- ``prove``
- ``verify``

If a given algorithm is not specified for the configured framework, the execution of ``make config=/path/to/config.json`` will fail.
