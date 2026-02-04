# Tutorial for Adding OPOCH Circuits in zk-Harness

OPOCH is a specialized proof system for SHA-256 hash chains. Unlike general-purpose ZK frameworks, adding new circuit types requires modifying the core OPOCH library.

## Current Limitations

OPOCH currently only supports the `sha256_chain` circuit. This is by design - OPOCH uses semantic hashing optimized specifically for SHA-256 computation chains.

## Adding New Input Sizes

To benchmark additional input sizes (N = number of SHA-256 iterations):

1. Create a new input file in `input/circuit/opoch/sha256_chain/`:

```json
{"n": 2048}
```

2. The benchmark will automatically pick up all JSON files in the input directory.

## Adding New Circuit Types

Adding new circuit types to OPOCH requires modifications to the core library at https://github.com/chetannothingness/opoch-hash.

If you extend OPOCH with a new circuit type:

1. Add the circuit implementation in the OPOCH library
2. Update `src/bin/zkharness_bench.rs` to handle the new circuit:
   - Add argument parsing for the new circuit name
   - Add a benchmark function similar to `benchmark_sha256_chain`
3. Create input files in `input/circuit/opoch/<circuit_name>/`
4. Update the config files in `input/config/opoch/` to include the new circuit

## Add Input File (Optionally)

If you add a new circuit that is not implemented in another framework, please include some input files in `input/circuit/<circuit_name>/`.

For more information about the input files please check the zk-Harness documentation.
