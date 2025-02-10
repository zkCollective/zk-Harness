# Tutorial for Adding gnark Circuit in zk-harness

To add a new circuit please follow the next steps.

## Adding a new circuit

1. Add the circuit source code: You should should add sour source code file(s) in a new directory named `circuits/<category>/<circuit_name>.go`.
If you need to use any external libraries please include them in your project directory.

2. Adding a test compliant with the gnark testing suite: Each new circuit should be tested with a test compliant with the gnark testing suite.
By default, the gnark testing suite runs the circuit over all curves and checks whether tests work for all backends. It is important that a newly added circuit passes these tests.

3. Edit circuits.go: In ``circuits.go`` add a new case handling your newly added circuit in the ``Witness`` function. Similarly, you should add a case in the ``Circuit`` function and the ``Init()`` function in the format ``BenchCircuits["<circuit_cmd>"] = &defaultCircuit{}``.

4. Updating the main config file: Please update the main config file for gnark circuits benchmarking in `../_input/config/gnark/config_all_circuits.json`. This config can be run to benchmark the whole gnark integration over all fields, curves and circuits.

## Add Input File (Optionally)

If you add a new circuit that is not implemented in another framework please include some input files in `../_input/circuit/circuit_name/`.
For more information about the input files please check `../documentation/`.
