use plonky2::fri::reduction_strategies::FriReductionStrategy;
use plonky2::fri::{FriConfig};
use starky::config::StarkConfig;
use std::{env, fs::File};
use std::process;
use std::io::{self, BufReader, Read, Write};
use std::fs;

// Explanation of parameters
// rate_bits - Reed solomon code rate
// cap_height - Height of the Merkle Tree caps
// proof_of_work_bits - The number of leading 0 bits in the grinding mechanism.
// The number of the leading zeros defines a certain amount of work that the prover must perform
// before generating the randomness representing the queries.

// Suggestions from ETHStarK for 128 bits:
// rate_bits = 2 -> 0.25 rate
// proof_of_work_bits = 20
// num_query_rounds = 55
// extension degree = 3

pub fn secure_config() -> StarkConfig {
    StarkConfig {
        security_bits: 128,
        num_challenges: 4,
        fri_config: FriConfig {
            rate_bits: 2,
            cap_height: 8,
            proof_of_work_bits: 20,
            reduction_strategy: FriReductionStrategy::ConstantArityBits(4, 5),
            num_query_rounds: 90,
        },
    }
}

pub fn standard_fast_config() -> StarkConfig {
    StarkConfig {
        security_bits: 100,
        num_challenges: 2,
        fri_config: FriConfig {
            rate_bits: 1,
            cap_height: 4,
            proof_of_work_bits: 16,
            reduction_strategy: FriReductionStrategy::ConstantArityBits(4, 5),
            num_query_rounds: 84,
        },
    }
}

pub fn plonky2_config() -> StarkConfig {
    StarkConfig{
        security_bits: 100,
        num_challenges: 2,
        fri_config: FriConfig {
            rate_bits: 1,
            cap_height: 4,
            proof_of_work_bits: 16,
            reduction_strategy: FriReductionStrategy::ConstantArityBits(4, 5),
            num_query_rounds: 84,
        },
    }
}

pub fn read_file_contents(file_name: String) -> String {
    let mut file = File::open(file_name).expect("Cannot load file");
    let mut file_str = String::new();
    file.read_to_string(&mut file_str).expect("Cannot read file");
    return file_str;
}

pub fn read_file_from_env_var(env_var_name: String) -> String {
    let input_file = env::var(env_var_name.clone()).unwrap_or_else(|_| {
        println!("Please set the {} environment variable to point to the input file", env_var_name);
        process::exit(1);
    });
    return read_file_contents(input_file);
}

pub fn read_env_variable(env_var_name: String) -> String {
    let variable_str = env::var(env_var_name.clone()).unwrap_or_else(|_| {
        println!("Please set the {} environment variable", env_var_name);
        process::exit(1);
    });
    return variable_str;
}

pub fn save_proof(proof_file: String, proof: &[u8]) -> io::Result<()> {
    let mut file = File::create(&proof_file)?;
    file.write_all(proof)?;
    Ok(())
}