use plonky2::fri::reduction_strategies::FriReductionStrategy;
use plonky2::fri::{FriConfig};
use starky::config::StarkConfig;

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
            cap_height: 4,
            proof_of_work_bits: 20,
            reduction_strategy: FriReductionStrategy::ConstantArityBits(4, 5),
            num_query_rounds: 90,
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
