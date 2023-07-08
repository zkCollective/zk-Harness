use std::marker::PhantomData;

use plonky2::field::extension::{Extendable, FieldExtension};
use plonky2::field::packed::PackedField;
use plonky2::field::polynomial::PolynomialValues;
use plonky2::hash::hash_types::RichField;
use plonky2::plonk::circuit_builder::CircuitBuilder;
use plonky2::field::types::Field;

use plonky2::plonk::config::{GenericConfig};
use starky::constraint_consumer::{ConstraintConsumer, RecursiveConstraintConsumer};
use starky::stark::Stark;
use starky::util::trace_rows_to_poly_values;
use starky::vars::{StarkEvaluationTargets, StarkEvaluationVars};
use serde::{Serialize, Deserialize};

/// Toy STARK system used for testing, exponentiation circuis
// y == x**e, x = 2 e = 3 y = 8
//     x    |    e      |    y    
//
//     2    |    0      |    1
//     2    |    1      |    2
//     2    |    2      |    4

#[derive(Copy, Clone)]
pub struct ExponentiateStark<F: RichField + Extendable<D>, const D: usize> {
    num_rows: usize,
    _phantom: PhantomData<F>,
}

impl<F: RichField + Extendable<D>, const D: usize> ExponentiateStark<F, D> {
    // The first public input is `x`.
    const PI_INDEX_X: usize = 0;
    // The second public input is `e`.
    const PI_INDEX_E: usize = 1;
    // The third public input is the second element of the last row, which should be equal to the
    // `num_rows`-th exponentiation number.
    const PI_INDEX_RES: usize = 2;

    pub fn new(num_rows: usize) -> Self {
        Self {
            num_rows,
            _phantom: PhantomData,
        }
    }

    pub fn generate_trace(&self, x: F, e: F, y: F) -> Vec<PolynomialValues<F>> {
        let mut trace_rows = (0..self.num_rows)
            .scan([x, e, y, F::ZERO], |acc, _| {
                let tmp = *acc;
                acc[0] = tmp[0];
                acc[1] = tmp[1] + F::ONE;
                acc[2] = tmp[2] * tmp[0];
                acc[3] = tmp[3] + F::ONE; // This is a dummy column 
                Some(tmp)
            })
            .collect::<Vec<_>>();
        trace_rows[self.num_rows - 1][3] = F::ZERO; // So that column 2 and 3 are permutation of one another.
        trace_rows_to_poly_values(trace_rows)
    }
}

impl<F: RichField + Extendable<D>, const D: usize> Stark<F, D> for ExponentiateStark<F, D> {
    // FIXME - Set to 4 as COLUMNS = 3 leads to an error in line 68
    // Once Fixed, delete dummy column in trace generation
    const COLUMNS: usize = 4; 
    const PUBLIC_INPUTS: usize = 3;

    fn eval_packed_generic<FE, P, const D2: usize>(
        &self,
        vars: StarkEvaluationVars<FE, P, { Self::COLUMNS }, { Self::PUBLIC_INPUTS }>,
        yield_constr: &mut ConstraintConsumer<P>,
    ) where
        FE: FieldExtension<D2, BaseField = F>,
        P: PackedField<Scalar = FE>,
    {
        // Check public inputs.
        yield_constr
            .constraint_first_row(vars.local_values[0] - vars.public_inputs[Self::PI_INDEX_X]);
        yield_constr
            .constraint_first_row(vars.local_values[1] - vars.public_inputs[Self::PI_INDEX_E]);
        yield_constr
            .constraint_last_row(vars.local_values[2] - vars.public_inputs[Self::PI_INDEX_RES]);

        // x' = x
        yield_constr.constraint_transition(vars.next_values[0] - vars.local_values[0]);
        // e' = e - 1
        yield_constr.constraint_transition(vars.next_values[1] - (vars.local_values[1] + FE::ONE));
        // y' = y * x
        yield_constr.constraint_transition(vars.next_values[2] - (vars.local_values[2] * vars.local_values[0]));
    }


    fn eval_ext_circuit(
        &self,
        _builder: &mut CircuitBuilder<F, D>,
        _vars: StarkEvaluationTargets<D, { Self::COLUMNS }, { Self::PUBLIC_INPUTS }>,
        _yield_constr: &mut RecursiveConstraintConsumer<F, D>,
    ) {
        // TODO
        // Currently not filled as this is only relevant for recursion
    }

    fn constraint_degree(&self) -> usize {
        2
    }
}

#[derive(Debug, Deserialize, Serialize)]
pub struct ExponentiateInput {
    X: String,
    Y: String,
    E: String,
}

pub fn get_exponentiate_data<C: GenericConfig<D>, const D: usize>(
    input_str: String
) -> (C::F, i32, C::F)
{
    let input: ExponentiateInput = serde_json::from_str(&input_str)
        .expect("JSON was not well-formatted");

    let x_64: u64 = input.X.parse().expect("Failed to parse X as u64");
    let e: i32 = input.E.parse().expect("Failed to parse E as usize");

    let mut x = C::F::from_canonical_u64(x_64);

    // Compute y as a value in the field
    for _ in 1..e {
        x = x.scalar_mul(x);
    }
    let y = x;

    return (x, e, y);
}

pub fn exponentiate<F: Field>(n: usize, x: F) -> F {
    (0..n).fold((F::ONE, x), |acc, _| (acc.1, acc.1 * x)).1
}

#[cfg(test)]
mod tests {
    use anyhow::Result;
    use plonky2::field::types::Field;
    use plonky2::plonk::config::{
        GenericConfig, PoseidonGoldilocksConfig,
    };
    use plonky2::util::timing::TimingTree;

    use starky::config::StarkConfig;
    use crate::circuits::exponentiate::ExponentiateStark;
    use crate::circuits::exponentiate::{
        exponentiate,
        get_exponentiate_data
    };
    use starky::prover::prove;
    use starky::verifier::verify_stark_proof;

    use starky_utils;

    #[test]
    fn test_exponentiate_stark() {
        const D: usize = 2;
        type C = PoseidonGoldilocksConfig;
        type F = <C as GenericConfig<D>>::F;
        type S = ExponentiateStark<F, D>;

        let config = StarkConfig::standard_fast_config();
        let num_rows = 1 << 10;
        let public_inputs = [F::ONE, F::from_canonical_usize(num_rows), exponentiate(num_rows - 1, F::ONE)];
        let stark = S::new(num_rows);
        let trace = stark.generate_trace(public_inputs[0], public_inputs[1], public_inputs[2]);
        let proof = prove::<F, C, S, D>(
            stark,
            &config,
            trace,
            public_inputs,
            &mut TimingTree::default(),
        ).expect("Proof generation failed");

        verify_stark_proof(stark, proof, &config).expect("Proof verification failed");
    }

    #[test]
    fn test_get_data() -> Result<(), Box<dyn std::error::Error>> { 
        let filename = "../../_input/circuit/exponentiate_2/input_16.json";
        let input_str = starky_utils::read_file_contents(filename.to_string());
        let (_x, _e, _y) = get_exponentiate_data::<PoseidonGoldilocksConfig, 2>(input_str);
        Ok(())
    }
    
}