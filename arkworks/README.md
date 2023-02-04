# arkworks library

## Plain Setup

### Installation

Arkworks operates as a set of libraries in the Rust ecosystem.
To compile Rust programs, you need to install cargo and the rust compiler which can be done through ``` rustup ```.

``` curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh ```

and successively run ``` source "$HOME/.cargo/env" ```.

### Compilation & Proof - Toy Examples

This folder provides arkworks implementationsof the following toy examples/circuits:

- Merkle Tree Membership Proof

To run commands / proofs for the toy examples you can leverage the following commands:

- ``` make test-toy-<name> ``` - Run tests for a toy example, e.g. ``` make test-toy-merkle ```
