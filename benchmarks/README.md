# Benchmarks

This folder contains the raw benchmark results.

## Structure

```
benchmarks/<circuits or math>/<machine>/<framework>/
```

In the `math` directory, we save the raw results for low-level arithmetic (e.g., field operations and curve operations),
whereas the `circuits` directory contains the results for the circuits (e.g., exponentiate circuit)

## Supported Machines

* __m5.2xlarge__: AWS machine, 8 vCPU, 32 GB RAM, Intel Xeon Platinum 8000 @3.1 GHz 
* __MAC_M1_PRO__: MAC laptop, 10 CPU, 64 GB RAM, Apple M1 Max @3.22 GHz

__NOTE__: The experiments in the laptops are not performed in a controlled environment.
