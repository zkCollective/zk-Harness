# Installation Instructions

To run most of the benchmarks, you will need to have the following software installed on your system:

1. Rust: Rust is a systems programming language that is used by some of the benchmarks. You can install Rust by following the instructions provided on the official website: https://www.rust-lang.org/tools/install.
2. Go: Go is a programming language that is used by some of the benchmarks. You can install Go by following the instructions provided on the official website: https://golang.org/doc/install.
3. Node.js: Node.js is a JavaScript runtime environment that is used by some of the benchmarks. You can install Node.js by following the instructions provided on the official website: https://nodejs.org.
4. Cargo Criterion: Cargo Criterion is a benchmarking library for Rust. To install Cargo Criterion, open a terminal and run the following command:

```bash
cargo install cargo-criterion
```

This will install Cargo Criterion and its dependencies.


Once you have installed the above software, you should be ready to run most of the benchmarks. 

## Additional Dependencies

For running benchmarks for the following projects, you should follow the specific instructions.

### Circom

Install Circom compiler following the instructions here: https://docs.circom.io/getting-started/installation/#installing-dependencies.
You also need to install [jq](https://jqlang.github.io/jq/).

To install `jq` you can use `brew`, `apt`, or `yum`.

* Using SNARKJS

```
npm install -g snarkjs
```

* Using rapidsnark

Note that this would only work in Intel64

```
cd frameworks/circom
git submodule init && git submodule update
cd rapidsnark
# Following the instructions in the README
```

### ffiasm

Note that this would only work in Intel64

* Linux (ubuntu)

```
apt install libgmp-dev nasm
```

* Mac

```
brew install gmp nasm
```
