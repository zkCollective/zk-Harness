#![allow(clippy::needless_range_loop)]

use std::marker::PhantomData;

use plonky2::{
    field::{
        extension::{Extendable, FieldExtension},
        packed::PackedField,
    },
    hash::hash_types::RichField,
    plonk::circuit_builder::CircuitBuilder,
};

use crate::{
    constraint_consumer::{ConstraintConsumer, RecursiveConstraintConsumer},
    stark::Stark,
    vars::{StarkEvaluationTargets, StarkEvaluationVars},
};

pub mod generation;
pub mod layout;

use layout::*;

#[derive(Copy, Clone)]
pub struct MerkleTree5STARK<F: RichField + Extendable<D>, const D: usize> {
    _phantom: PhantomData<F>,
}

impl<F: RichField + Extendable<D>, const D: usize> MerkleTree5STARK<F, D> {
    pub fn new() -> Self {
        Self::default()
    }
}

impl<F: RichField + Extendable<D>, const D: usize> Default for MerkleTree5STARK<F, D> {
    fn default() -> Self {
        Self {
            _phantom: PhantomData,
        }
    }
}

// Padding unnecessary since we're doing always 2^n - 1 hashes for some n, and we have 1 row per hash + 1 for output

impl<F: RichField + Extendable<D>, const D: usize> Stark<F, D> for MerkleTree5STARK<F, D> {
    const COLUMNS: usize = NUM_COLS;
    const PUBLIC_INPUTS: usize = NUM_PUBLIC_INPUTS;

    fn eval_packed_generic<FE, P, const D2: usize>(
        &self,
        vars: StarkEvaluationVars<FE, P, { Self::COLUMNS }, { Self::PUBLIC_INPUTS }>,
        yield_constr: &mut ConstraintConsumer<P>,
    ) where
        FE: FieldExtension<D2, BaseField = F>,
        P: PackedField<Scalar = FE>,
    {
        let curr_row = vars.local_values;
        let next_row = vars.next_values;
        let pis = vars.public_inputs;

        // load leaves in at first row
        for i in 0..TREE_WIDTH {
            for word in 0..WORDS_PER_HASH {
                // degree 1
                yield_constr.constraint_first_row(
                    curr_row[val_i_word(i, word)] - pis[pi_leaf_i_word(i, word)],
                );
            }
        }

        // set PC to 0 first row, hash idx to 1, curr half level width to TREE_WIDTH / 2
        // degree 1
        yield_constr.constraint_first_row(curr_row[PC]);
        yield_constr.constraint_first_row(
            curr_row[HALF_LEVEL_WIDTH] - FE::from_canonical_u64(TREE_WIDTH as u64 / 2),
        );
        yield_constr.constraint_first_row(P::ONES - curr_row[HASH_IDX]);

        let level_selectors = [
            curr_row[level_flag(0)],
            curr_row[level_flag(1)],
            curr_row[level_flag(2)],
            curr_row[level_flag(3)],
            P::ONES
                - curr_row[level_flag(0)]
                - curr_row[level_flag(1)]
                - curr_row[level_flag(2)]
                - curr_row[level_flag(3)],
        ];

        let next_level_selectors = [
            next_row[level_flag(0)],
            next_row[level_flag(1)],
            next_row[level_flag(2)],
            next_row[level_flag(3)],
            P::ONES
                - next_row[level_flag(0)]
                - next_row[level_flag(1)]
                - next_row[level_flag(2)]
                - next_row[level_flag(3)],
        ];

        // set level flags for first row
        // 1000 => 0th row
        // 0100 -> 1st row
        // 0010 -> 2nd row
        // 0001 -> 3rd row
        // 0000 -> 4th row (root)
        yield_constr.constraint_first_row(P::ONES - level_selectors[0]);
        yield_constr.constraint_first_row(level_selectors[1]);
        yield_constr.constraint_first_row(level_selectors[2]);
        yield_constr.constraint_first_row(level_selectors[3]);

        let is_flag_transition_i = |i| {
            level_selectors[i] * next_level_selectors[i + 1]
        };
        let flag_transition_i = |i| {
            level_selectors[i] - next_level_selectors[i + 1]
        };
        let is_transition: P = (0..TREE_DEPTH-1).map(|level| is_flag_transition_i(level)).sum();
        yield_constr.constraint_transition(curr_row[LEVEL_DONE_FLAG] - is_transition);

        // abvance level flags at correct indices
        // one_if_end_of_level is 0 during the single row for root, 
        let one_if_end_of_level = curr_row[PC] - curr_row[HALF_LEVEL_WIDTH];
        for level in 0..TREE_DEPTH-1 {
            let transition = flag_transition_i(level);
            // degree 2
            yield_constr.constraint_transition(one_if_end_of_level * curr_row[LEVEL_DONE_FLAG] * transition);
        }

        // each row, increment PC unless we're done with the current level
        // degree 2
        yield_constr.constraint_transition(
            (P::ONES - curr_row[LEVEL_DONE_FLAG]) * (next_row[PC] - (curr_row[PC] + P::ONES)),
        );

        // divide half level width by two when we reset PC unless it's the last time
        // degree 3
        yield_constr.constraint_transition(
           curr_row[LEVEL_DONE_FLAG] * (P::ONES - next_level_selectors[4])
                * (next_row[HALF_LEVEL_WIDTH] * FE::TWO - curr_row[HALF_LEVEL_WIDTH]),
        );

        // load leftmost two hashes in the val cols into input cols for lookup except for last row
        // note: constraint_transition doesn't apply to last row
        for word in 0..WORDS_PER_HASH {
            // degree 1
            yield_constr.constraint_transition(
                curr_row[hash_input_0_word(word)]
                    - (curr_row[val_i_word(0, word)]
                        + curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32)),
            );
            yield_constr.constraint_transition(
                curr_row[hash_input_1_word(word)]
                    - (curr_row[val_i_word(1, word)]
                        + curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32)),
            );
        }

        // load the output hash into the rightmost val col of the next row unless unless we're moving to the next level.
        for word in 0..WORDS_PER_HASH {
            yield_constr.constraint_transition(
                (P::ONES - curr_row[LEVEL_DONE_FLAG])
                    * (next_row[val_i_word(15, word)]
                        - (curr_row[hash_output_word(word)]
                        - curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32))),
            );
        }

        // set i/o filters to 1 except for last row (constraint_transition does the filtering for us)
        // degree 1
        yield_constr.constraint_transition(P::ONES - curr_row[INPUT_FILTER]);
        yield_constr.constraint_transition(P::ONES - curr_row[OUTPUT_FILTER]);

        // shift hash inputs left by two unless we're at the end of a level (i.e. PC == HALF_LEVEL_WIDTH - 1)
        for i in 0..(TREE_WIDTH - 2) {
            for word in 0..WORDS_PER_HASH {
                // degree 2
                yield_constr.constraint_transition(
                    (P::ONES - curr_row[LEVEL_DONE_FLAG])
                        * (next_row[val_i_word(i, word)] - curr_row[val_i_word(i + 2, word)]),
                );
            }
        }
        // zero the second to last hash in var cols of next row (the hash output went into the last one above)
        for word in 0..WORDS_PER_HASH {
            yield_constr.constraint_transition(
                (P::ONES - curr_row[LEVEL_DONE_FLAG]) * next_row[val_i_word(14, word)],
            );
        }

        // if we're at the end of a level, coalesce the level's outputs to the left. At the last row of the 0th level, the vals will look like this:
        // h_14_0 | h_15_0 | 0 | h_0_1 | 0 | h_1_1 | 0 | h_2_1 | 0 | h_3_1 | 0 | h_4_1 | 0 | h_5_1 | 0 | h_6_1 |
        // and h_7_1 will be OUTPUT_COL, where h_i_j deonetes the ith hash for the th level
        // what we want to do is, in one fell swoop, shift left by two *and* "unspread" the the h_i_1s so the next row looks like this:
        // h_0_1 | h_1_1 | h_2_1 | h_3_1 | h_4_1 | h_5_1 | h_6_1 | h_7_1 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |
        // here, we want to shift h_0_1 left by 3, h_1_1 by 4, h_2_1 by 5. This generalizes to shifting h_i_1 by 3 + i.
        // however, at the end of the 1st level, it'll look like this:
        // h_6_1 | h_7_1 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | h_0_2 | 0 | h_1_2 | 0 | h_2_2 |
        // here, we want to shift h_0_2 by 11, h_1_2 by 12, h_2_2 by 13. This generalizes to shifting h_i_2 by 11 + i.
        // what we actually want is to shift h_i_j left by TREE_WIDTH - LEVEL_WIDTH + 3 + i
        // e.g. the 0th level has a width of 16, formula gives 16 - 16 + 3 + i = 3 + i
        // e.g. 1st level has a width of 8, formula gives 16 - 8 + 3 + i = 11 + i
        // `get_level_end_shift` in `layout.rs` computes this formula

        for level in 0..4 {
            let next_level_width = level_width(level + 1);
            let sel = is_flag_transition_i(level);
            for i in 0..(next_level_width - 1) {
                let shift_amount = get_level_end_shift(i, level);
                for word in 0..WORDS_PER_HASH {
                    // degree 3
                    yield_constr.constraint_transition(
                        sel * (next_row[val_i_word(i, word)]
                            - curr_row[val_i_word(i + shift_amount, word)]),
                    )
                }
            }

            // append the hash from OUTPUT_COL
            for word in 0..WORDS_PER_HASH {
                // degree 3
                yield_constr.constraint_transition(
                    sel * (next_row[val_i_word(next_level_width - 1, word)]
                        - (curr_row[hash_output_word(word)]
                        - curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32))),
                );
            }

            // zero the rest of the next row's var cols
            for i in next_level_width..TREE_WIDTH {
                for word in 0..WORDS_PER_HASH {
                    // degree 3
                    yield_constr.constraint_transition(sel * next_row[val_i_word(i, word)]);
                }
            }
        }

        // in the last row, check that the root hash given in PIs is the same as the one sitting in the leftmost val cols
        for word in 0..WORDS_PER_HASH {
            yield_constr
                .constraint_last_row(curr_row[val_i_word(0, word)] - pis[pi_root_word(word)]);
        }

        // ensure binary flags are binary
        yield_constr.constraint(curr_row[LEVEL_DONE_FLAG] * (P::ONES - curr_row[LEVEL_DONE_FLAG]));
        yield_constr.constraint(curr_row[INPUT_FILTER] * (P::ONES - curr_row[INPUT_FILTER]));
        yield_constr.constraint(curr_row[OUTPUT_FILTER] * (P::ONES - curr_row[OUTPUT_FILTER]));
        for level in 0..4 {
            yield_constr
                .constraint(curr_row[level_flag(level)] * (P::ONES - curr_row[level_flag(level)]));
        }

        // ensure level selectors sum to 1 => at most one level flag is active since we checked the flags are all binary
        yield_constr.constraint(P::ONES - level_selectors.into_iter().sum::<P>());
    }

    fn eval_ext_circuit(
        &self,
        _builder: &mut CircuitBuilder<F, D>,
        _vars: StarkEvaluationTargets<D, { Self::COLUMNS }, { Self::PUBLIC_INPUTS }>,
        _yield_constr: &mut RecursiveConstraintConsumer<F, D>,
    ) {
        todo!()
    }

    fn constraint_degree(&self) -> usize {
        3
    }
}

#[cfg(test)]
mod tests {
    use anyhow::Result;
    use plonky2::plonk::config::{GenericConfig, PoseidonGoldilocksConfig};
    use plonky2::util::timing::TimingTree;

    use super::*;
    use crate::config::StarkConfig;
    use crate::prover::prove;
    use crate::tree_stark::generation::TreeTraceGenerator;
    use crate::stark_testing::test_stark_low_degree;
    use crate::verifier::verify_stark_proof;

    #[test]
    fn test_stark_degree() -> Result<()> {
        const D: usize = 2;
        type C = PoseidonGoldilocksConfig;
        type F = <C as GenericConfig<D>>::F;
        type S = MerkleTree5STARK<F, D>;

        let stark = S::new();
        test_stark_low_degree(stark)
    }

    // #[test]
    // fn test_stark_circuit() -> Result<()> {
    //     const D: usize = 2;
    //     type C = PoseidonGoldilocksConfig;
    //     type F = <C as GenericConfig<D>>::F;
    //     type S = Sha2CompressionStark<F, D>;

    //     let stark = S::new();

    //     test_stark_circuit_constraints::<F, C, S, D>(stark)
    // }

    #[test]
    fn test_tree_stark() -> Result<()> {
        const D: usize = 2;
        type C = PoseidonGoldilocksConfig;
        type F = <C as GenericConfig<D>>::F;
        type S = MerkleTree5STARK<F, D>;

        let mut leaves = [[0; 8]; TREE_WIDTH];
        for i in 0..TREE_WIDTH {
            for word in 0..WORDS_PER_HASH {
                leaves[i][word] = (i * WORDS_PER_HASH + word) as u32;
            }
        }

        let mut generator = TreeTraceGenerator::<F>::new(16, leaves);
        let (_root, pis) = generator.gen();
        let trace = generator.into_polynomial_values();

        let config = StarkConfig::standard_fast_config();
        let stark = S::new();
        let mut timing = TimingTree::default();
        let proof = prove::<F, C, S, D>(stark, &config, trace, pis, &mut timing)?;

        verify_stark_proof(stark, proof, &config)?;

        Ok(())
    }
}
