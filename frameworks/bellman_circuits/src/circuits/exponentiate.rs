use bellman::{
    gadgets::num::AllocatedNum,
    gadgets::multipack,
    Circuit, ConstraintSystem, SynthesisError,
};
use ff::PrimeField;
use serde::{Serialize, Deserialize};
use bls12_381::{Scalar};
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
) -> (Scalar, usize, Scalar){
    let input: ExponentiateInput = serde_json::from_str(&input_str)
        .expect("JSON was not well-formatted");
    let x_64: u64 = input.X.parse().expect("Failed to parse X as u64");
    let e: usize = input.E.parse().expect("Failed to parse E as usize");

    let mut x = Scalar::from(x_64);

    // Compute y as a Scalar value in the field
    for _ in 1..e {
        x = x.mul(&x);
    }
    let y = x;

    return (x, e, y);
}

#[cfg(test)]
mod exponentiate_tests {
    use super::*;
    use bls12_381::{Scalar};
    use bellman::gadgets::test::TestConstraintSystem;
    use bellman::groth16;
    use rand::rngs::OsRng;
    use bls12_381::Bls12;

    #[test]
    fn test_cs_satisfaction() {
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

    #[test]
    fn test_circuit_correctness() {
        // Convert x, y from string to Fr and prepare public input
        let x = Scalar::from(1u64);
        let y = Scalar::from(1u64);
        let e = 4;

        // Define the circuit
        let circuit = ExponentiationCircuit {
            x: Some(x),
            e: e,
            y: Some(y),
        };

        // Generate Parameters
        let params = groth16::generate_random_parameters::<Bls12, _, _>(circuit.clone(), &mut OsRng).unwrap();

        // Create a mock constraint system
        let mut cs = TestConstraintSystem::<Scalar>::new();
        // Synthesize the circuit with our mock constraint system
        circuit.clone().synthesize(&mut cs).unwrap();
        println!("Number of constraints: {}", cs.num_constraints());

        let rng = &mut OsRng;
        let pvk = groth16::prepare_verifying_key(&params.vk);
        let proof = groth16::create_random_proof(circuit.clone(), &params, rng).unwrap(); 

        // Check the proof!
        assert!(groth16::verify_proof(&pvk, &proof, &[]).is_ok());
    }
}
