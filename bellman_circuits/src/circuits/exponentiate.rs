use bellman::{
    gadgets::num::AllocatedNum,
    Circuit, ConstraintSystem, SynthesisError,
};
use ff::PrimeField;
use serde::{Serialize, Deserialize};
use serde_json;

#[derive(Clone)]
pub struct ExponentiationCircuit<Scalar: PrimeField> {
    pub x: Option<Scalar>,
    pub e: usize, // e is the number of iterations (exponent)
    pub y: Option<Scalar>,
}

impl<Scalar: PrimeField> Circuit<Scalar> for ExponentiationCircuit<Scalar> {
    fn synthesize<CS: ConstraintSystem<Scalar>>(self, cs: &mut CS) -> Result<(), SynthesisError> {
        // Allocate the base value (public)
        let x = AllocatedNum::alloc(cs.namespace(|| "x"), || {
            self.x.ok_or(SynthesisError::AssignmentMissing)
        })?;
        
        let mut exponentiation_result = x.clone(); // Initialize result as x

        // Perform exponentiation by repeated multiplication
        for i in 1..self.e {
            exponentiation_result = exponentiation_result.mul(cs.namespace(|| format!("multiply {}", i)), &x)?;
        }        

        // Allocate the result value (public)
        let y = AllocatedNum::alloc(cs.namespace(|| "y"), || {
            self.y.ok_or(SynthesisError::AssignmentMissing)
        })?;

        // Assert that exponentiation_result equals to y
        cs.enforce(
            || "result = y",
            |lc| lc + exponentiation_result.get_variable(),
            |lc| lc + CS::one(),
            |lc| lc + y.get_variable(),
        );

        Ok(())
    }
}

#[derive(Debug, Deserialize, Serialize)]
pub struct ExponentiateInput {
    X: String,
    Y: String,
    E: String,
}

pub fn get_exponentiate_data (
    input_str: String
) -> (u64, usize, u64){
    let input: ExponentiateInput = serde_json::from_str(&input_str)
        .expect("JSON was not well-formatted");
    let x: u64 = input.X.parse().expect("Failed to parse X as u64");
    let y: u64 = input.Y.parse().expect("Failed to parse Y as u64");
    let e: usize = input.E.parse().expect("Failed to parse E as usize");
    return (x, e, y);
}

#[cfg(test)]
mod exponentiate_tests {
    use super::*;
    use bls12_381::{Scalar};
    use bellman::gadgets::test::TestConstraintSystem;

    #[test]
    fn test_exponentiation_circuit() {
        // Convert x, y from string to Fr and prepare public input
        let x = Scalar::from(2u64);
        let y = Scalar::from(16u64);
        let e = 4;

        // Create an instance of our circuit (with the x, e and y as a witness).
        let c = ExponentiationCircuit {
            x: Some(x),
            e: e,
            y: Some(y),
        };

        // Create a mock constraint system
        let mut cs = TestConstraintSystem::<Scalar>::new();

        // Synthesize the circuit with our mock constraint system
        c.synthesize(&mut cs).unwrap();

        println!("{}", cs.num_constraints());

        // If the constraint system is not satisfied, we can print the unsatisfied constraints to help debugging
        if !cs.is_satisfied() {
            cs.which_is_unsatisfied().map(|constraint| {
                println!("Unsatisfied constraint: {}", constraint);
            });
        }

        // Check if the constraint system is satisfied
        assert!(cs.is_satisfied());
    }

    // FIXME - Write test for circuit correctness. 
    // Currently, proof verification fails due to parsing of public inputs
}
