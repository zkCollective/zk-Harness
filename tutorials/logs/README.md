# Logs Description

This document provides an overview of the log format for benchmarks in the ZKP community. To learn how to generate logs for a specific language or library, please refer to the respective page.

There are 4 different benchmark categories:

* Field Arithmetic
* Elliptic Curve Group Operations
* Circuits
* Recursion

The logs are saved as CSV files and are processed to compare the performance of different ZKP languages and libraries.

## Logs Directory

All benchmark logs should be saved in a directory within the `benchmark` directory.
For example, the logs for the cubic circuit benchmark can be found at: `benchmarks/gnark/circuits/cubic.csv`.
The directory structure should follow this format:

```
benchmarks/[language/library name]/[benchmark category]/benchmark name
```

## Logs Content

In the following sections, we describe the columns in the CSV file for each benchmark category.

### Field Arithmetic

The following information is recorded for each field arithmetic benchmark:

* framework: the name of the framework (e.g., gnark)
* category: the category of the benchmark (i.e., arithmetic)
* curve: the curve of which field we use
* field: the benchmarked field (base or scalar)
* operation: the operation performed (add, sub, mul, inv, exp)
* input: file path of the input used 
* ram: memory consumed in bytes
* time: elapsed time in nanoseconds
* nbPhysicalCores: number of physical cores used
* nbLogicalCores: number of logical cores used
* machine: the machine used for benchmarking

### Elliptic Curve Group Operations

The following information is recorded for each elliptic curve group operation benchmark:

* framework: the name of the framework (e.g., gnark)
* category: the category of the benchmark (i.e., ec)
* curve: the benchmarked curve
* operation: the operation performed -- MSM, FFT/NTT, Pairing
* input: file path of the input used 
* ram: memory consumed in bytes
* time: elapsed time in milliseconds
* nbPhysicalCores: number of physical cores used
* nbLogicalCores: number of logical cores used
* machine: the machine used for benchmarking

### Circuits

The following information is recorded for each circuit benchmark:

* framework: the name of the framework (e.g., gnark)
* category: the category of the benchmark (i.e., circuit)
* backend: the backend used (e.g., groth16)
* curve: the curve used (e.g., bn256)
* circuit: the circuit being run
* input: file path of the input used 
* operation: the step being measured -- compile, witness, setup, prove, verify 
* nbConstraints: the number of constraints in the circuit
* nbSecret: number of secret inputs
* nbPublic: number of public inputs
* ram: memory consumed in bytes
* time: elapsed time in milliseconds
* proofSize: the size of the proof in bytes -- empty value when Operation != proving
* nbPhysicalCores: number of physical cores used
* nbLogicalCores: number of logical cores used
* machine: the machine used for benchmarking

### Recursion

The contents of the recursion logs are yet to be determined and will be added at a later date.
