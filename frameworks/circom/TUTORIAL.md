# Tutorial for Adding Circom Circuit in zk-harness

To add a new circuit please follow the next steps.

## Adding Circuit Source Code

You should should add sour source code file(s) in a new directory named `circuits/benchmarks/project_name`.
The `circom` files that includes the main component should be named `circuit.circom`.
If you need to use any external libraries please include them in your project directory.

## Updating Main Config File

Please update the main config file for Circom circuits benchmarking in `../_input/config/circom/config_all_circuits.json`

## Add Input File (Optionally)

If you add a new circuit that is not implemented in another framework please include some input files in `../_input/circuit/circuit_name/`. 
For more information about the input files please check `../documentation/`.

## Verify Circuit

Please make sure that the circuit is correct by providing some additional 
inputs and document how you implemented the circuit in the PR and in the source
code. Finally, include the command to run the new circuit using the 
`scripts/run_circuit.sh` script in the PR.
