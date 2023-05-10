extern crate rand;
extern crate bellman_ce;
extern crate criterion;

use std::rc::Rc;
use rand::{Rng};
use bellman_ce::{Circuit, ConstraintSystem, SynthesisError};
use bellman_ce::pairing::bn256::{Bn256};
use bellman_ce::pairing::ff::{Field, ScalarEngine};
use bellman_ce::groth16::{Parameters, generate_random_parameters};
use bellman_ce::pairing::{Engine};

const MIMC_ROUNDS: usize = 322;

// New function to generate the MiMC round constants
pub fn generate_constants<R: Rng>(rng: &mut R) -> Rc<Vec<<Bn256 as ScalarEngine>::Fr>> {
    Rc::new((0..MIMC_ROUNDS).map(|_| rng.gen()).collect())
}

// New function to generate parameters for the given circuit
pub fn generate_circuit_parameters<C: Circuit<Bn256>, R: Rng>(
    circuit: C,
    rng: &mut R,
) -> Result<Parameters<Bn256>, SynthesisError> {
    generate_random_parameters(circuit, rng)
}

// The rest of the MiMCDemo struct and impl remains the same
/// This is our demo circuit for proving knowledge of the
/// preimage of a MiMC hash invocation.
#[derive(Clone)]
pub struct MiMCDemo<E: Engine> {
    pub xl: Option<E::Fr>,
    pub xr: Option<E::Fr>,
    pub constants: Rc<Vec<E::Fr>>,
}

/// Our demo circuit implements this `Circuit` trait which
/// is used during paramgen and proving in order to
/// synthesize the constraint system.
impl<'a, E: Engine> Circuit<E> for MiMCDemo<E> {
    fn synthesize<CS: ConstraintSystem<E>>(
        self,
        cs: &mut CS
    ) -> Result<(), SynthesisError>
    {
        assert_eq!(self.constants.len(), MIMC_ROUNDS);

        // Allocate the first component of the preimage.
        let mut xl_value = self.xl;
        let mut xl = cs.alloc(|| "preimage xl", || {
            xl_value.ok_or(SynthesisError::AssignmentMissing)
        })?;

        // Allocate the second component of the preimage.
        let mut xr_value = self.xr;
        let mut xr = cs.alloc(|| "preimage xr", || {
            xr_value.ok_or(SynthesisError::AssignmentMissing)
        })?;

        for i in 0..MIMC_ROUNDS {
            // xL, xR := xR + (xL + Ci)^3, xL
            let cs = &mut cs.namespace(|| format!("round {}", i));

            // tmp = (xL + Ci)^2
            let tmp_value = xl_value.map(|mut e| {
                e.add_assign(&self.constants[i]);
                e.square();
                e
            });
            let tmp = cs.alloc(|| "tmp", || {
                tmp_value.ok_or(SynthesisError::AssignmentMissing)
            })?;

            cs.enforce(
                || "tmp = (xL + Ci)^2",
                |lc| lc + xl + (self.constants[i], CS::one()),
                |lc| lc + xl + (self.constants[i], CS::one()),
                |lc| lc + tmp
            );

            // new_xL = xR + (xL + Ci)^3
            // new_xL = xR + tmp * (xL + Ci)
            // new_xL - xR = tmp * (xL + Ci)
            let new_xl_value = xl_value.map(|mut e| {
                e.add_assign(&self.constants[i]);
                e.mul_assign(&tmp_value.unwrap());
                e.add_assign(&xr_value.unwrap());
                e
            });

            let new_xl = if i == (MIMC_ROUNDS-1) {
                // This is the last round, xL is our image and so
                // we allocate a public input.
                cs.alloc_input(|| "image", || {
                    new_xl_value.ok_or(SynthesisError::AssignmentMissing)
                })?
            } else {
                cs.alloc(|| "new_xl", || {
                    new_xl_value.ok_or(SynthesisError::AssignmentMissing)
                })?
            };

            cs.enforce(
                || "new_xL = xR + (xL + Ci)^3",
                |lc| lc + tmp,
                |lc| lc + xl + (self.constants[i], CS::one()),
                |lc| lc + new_xl - xr
            );

            // xR = xL
            xr = xl;
            xr_value = xl_value;

            // xL = new_xL
            xl = new_xl;
            xl_value = new_xl_value;
        }

        Ok(())
    }
}
