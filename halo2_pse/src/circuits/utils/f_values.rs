use halo2_proofs::arithmetic::Field;

pub fn get_one<F: Field>() -> F {
    // Depending on the version of the library one of the following is the correct one
    // F::ONE;
    F::ONE
}

pub fn get_zero<F: Field>() -> F {
    // Depending on the version of the library one of the following is the correct one
    // F::ZERO;
    F::ZERO
}
