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

* Framework: the name of the framework (e.g., gnark)
* Category: the category of the benchmark (i.e., arithmetic)
* Field: the benchmarked field (e.g., Native)
* p: the order of the field
* Operation: the operation performed
* Input: file path of the input used 
* Ram: memory consumed in bytes
* Time: elapsed time in milliseconds
* nbCores: number of cores used
* Machine: the machine used for benchmarking

### Elliptic Curve Group Operations

The following information is recorded for each elliptic curve group operation benchmark:

* Framework: the name of the framework (e.g., gnark)
* Category: the category of the benchmark (i.e., ec)
* Curve: the benchmarked curve
* Operation: the operation performed -- MSM, FFT/NTT, Pairing
* Input: file path of the input used 
* Ram: memory consumed in bytes
* Time: elapsed time in milliseconds
* nbCores: number of cores used
* Machine: the machine used for benchmarking

### Circuits

The following information is recorded for each circuit benchmark:

* Framework: the name of the framework (e.g., gnark)
* Category: the category of the benchmark (i.e., circuit)
* Backend: the backend used (e.g., groth16)
* Curve: the curve used (e.g., bn256)
* Circuit: the circuit being run
* Input: file path of the input used 
* Operation: the step being measured -- compile, setup, proving, verifying 
* nbConstraints: the number of constraints in the circuit
* nbSecret: number of secret inputs
* nbPublic: number of public inputs
* Ram: memory consumed in bytes
* Time: elapsed time in milliseconds
* ProofSize: the size of the proof -- empty value when Operation != proving
* nbCores: number of cores used
* Machine: the machine used for benchmarking

### Recursion

The contents of the recursion logs are yet to be determined and will be added at a later date.
