# rapidsnark

[rapidsnark](https://github.com/iden3/rapidsnark) is a fast zkSNARK prover implemented in C++.
It typically used along with Circom.

## Installation

You can install rapidsnark with `npm install -g rapidsnark` using [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm).

## Arithmetics and elliptic curves benchmarks

If you want to run the benchmarks for Arithmetics and EC then you need to 
execute the following command inside `rapidsnark` directory to download the
required library.

```
npm install tmp-promise ffiasm big-integer
```

### Example 

* Arithmetics benchmarking

```
./scripts/arithmetics.js bn128 scalar add 100 ../_input/arithmetic/add/input_1.json res.csv
```

* EC benchmarking

```
./scripts/curves.js bn128 g1 multi-scalar-multiplication 10 ../_input/ec/multi-scalar-multiplication/input_1.json res.csv
```
