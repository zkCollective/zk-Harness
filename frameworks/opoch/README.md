# Benchmarking the OPOCH Library

OPOCH is a STARK-based proof system optimized for SHA-256 hash chain verification.

**Key Result**: Sub-millisecond verification (~0.006ms) and constant 321-byte proofs regardless of computation size.

## Overview

| Property | Value |
|----------|-------|
| Proof System | STARK (FRI-based) |
| Field | Goldilocks (p = 2^64 - 2^32 + 1) |
| Security | 128-bit (68 FRI queries) |
| Trusted Setup | None (transparent) |
| Post-Quantum | Yes |

Currently the following circuits have been implemented:

* SHA-256 Chain (`sha256_chain`) - proves `y = SHA-256^N(x)`

## Installation

No external installation required. The benchmark binary builds automatically using a Cargo git dependency.

```bash
# Optional: build manually
cd frameworks/opoch
cargo build --release --bin zkharness_bench
```

## Run the benchmarks

Running a benchmark can be facilitated through the following commands:

```bash
# Run quick test (N=64)
make benchmark-opoch-test-circuit

# Run all benchmarks (N=64, 256, 512, 1024)
make benchmark-opoch-circuits
```

Or run directly:

```bash
cd frameworks/opoch
cargo run --release --bin zkharness_bench -- --circuit sha256_chain --input 64 --count 1
```

## Benchmark Results

| Input (N) | Prove Time | Verify Time | Proof Size | Constraints |
|-----------|------------|-------------|------------|-------------|
| 64        | ~6,400 ms  | ~0.006 ms   | 321 bytes  | 4,096       |
| 256       | ~27,000 ms | ~0.006 ms   | 321 bytes  | 16,384      |
| 512       | ~54,000 ms | ~0.006 ms   | 321 bytes  | 32,768      |
| 1024      | ~107,000 ms| ~0.006 ms   | 321 bytes  | 65,536      |

**Key Properties:**
- Verification time: O(1) - constant ~0.006ms regardless of N
- Proof size: O(1) - constant 321 bytes regardless of N
- Prover time: O(N) - linear in chain length

## Adding new circuits

See `TUTORIAL.md`

## Repository

- **Source**: https://github.com/chetannothingness/opoch-hash
- **License**: MIT
