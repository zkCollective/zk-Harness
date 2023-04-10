use arrayref::array_mut_ref;
use plonky2::field::{
    polynomial::PolynomialValues,
    types::{Field, PrimeField64},
};

use super::layout::*;
use crate::util::{is_power_of_two, compress, trace_rows_to_poly_values};

struct TreeTrace<F: Field>(Vec<[F; NUM_COLS]>);

impl<F: Field> TreeTrace<F> {
    fn new(max_rows: usize) -> TreeTrace<F> {
        assert!(
            is_power_of_two(max_rows as u64),
            "max_rows must be a power of two"
        );

        TreeTrace(vec![[F::ZERO; NUM_COLS]; max_rows])
    }
}

pub struct TreeTraceGenerator<F: Field> {
    trace: TreeTrace<F>,
    leaves: [[u32; 8]; TREE_WIDTH],
    idx: usize,
    level: usize,
    level_idx: usize,
}

impl<F: Field + PrimeField64> TreeTraceGenerator<F> {
    pub fn new(max_rows: usize, leaves: [[u32; 8]; TREE_WIDTH]) -> TreeTraceGenerator<F> {
        TreeTraceGenerator {
            trace: TreeTrace::new(max_rows),
            leaves,
            idx: 0,
            level: 0,
            level_idx: 0,
        }
    }

    fn max_rows(&self) -> usize {
        self.trace.0.len()
    }

    fn get_next_window(&mut self) -> (&mut [[F; NUM_COLS]; 2], usize, usize, usize) {
        let idx = self.idx;
        assert!(idx < self.max_rows(), "get_next_window exceeded MAX_ROWS");

        let level = self.level;
        let level_idx = self.level_idx;

        self.idx += 1;
        self.level_idx += 1;
        if self.level_idx == level_width(level + 1) {
            self.level_idx = 0;
            self.level += 1;
        }

        (array_mut_ref![self.trace.0, idx, 2], idx, level_idx, level)
    }

    fn gen_misc(
        curr_row: &mut [F; NUM_COLS],
        idx: usize,
        level_idx: usize,
        level: usize,
        max_rows: usize,
    ) {
        let level_width = level_width(level);
        let half_level_width = level_width / 2;

        // set pc and half level width
        curr_row[PC] = F::from_canonical_u64(level_idx as u64);
        curr_row[HALF_LEVEL_WIDTH] = F::from_canonical_u64(half_level_width as u64);

        // set hash idx and lookup filters for all rows but the last
        curr_row[HASH_IDX] = F::from_canonical_u64(idx as u64 + 1);
        if idx < max_rows - 1 {
            // hash idx is 1-indexed for domain separation
            curr_row[INPUT_FILTER] = F::ONE;
            curr_row[OUTPUT_FILTER] = F::ONE;
        } else {
            curr_row[INPUT_FILTER] = F::ZERO;
            curr_row[OUTPUT_FILTER] = F::ZERO;
        }

        // set level done flag
        if level_idx == half_level_width - 1 {
            curr_row[LEVEL_DONE_FLAG] = F::ONE;
        } else {
            curr_row[LEVEL_DONE_FLAG] = F::ZERO;
        }

        // set level flags
        match level {
            0 => {
                curr_row[level_flag(0)] = F::ONE;
                curr_row[level_flag(1)] = F::ZERO;
                curr_row[level_flag(2)] = F::ZERO;
                curr_row[level_flag(3)] = F::ZERO;
            }
            1 => {
                curr_row[level_flag(0)] = F::ZERO;
                curr_row[level_flag(1)] = F::ONE;
                curr_row[level_flag(2)] = F::ZERO;
                curr_row[level_flag(3)] = F::ZERO;
            }
            2 => {
                curr_row[level_flag(0)] = F::ZERO;
                curr_row[level_flag(1)] = F::ZERO;
                curr_row[level_flag(2)] = F::ONE;
                curr_row[level_flag(3)] = F::ZERO;
            }
            3 => {
                curr_row[level_flag(0)] = F::ZERO;
                curr_row[level_flag(1)] = F::ZERO;
                curr_row[level_flag(2)] = F::ZERO;
                curr_row[level_flag(3)] = F::ONE;
            }
            4 => {
                curr_row[level_flag(0)] = F::ZERO;
                curr_row[level_flag(1)] = F::ZERO;
                curr_row[level_flag(2)] = F::ZERO;
                curr_row[level_flag(3)] = F::ZERO;
            }
            _ => unreachable!(),
        }
    }

    pub fn gen(&mut self) -> ([u32; 8],[F; NUM_PUBLIC_INPUTS]) {
        let max_rows = self.max_rows();

        // load leaves into first row of val cols
        let first_row = &mut self.trace.0[0];
        for i in 0..TREE_WIDTH {
            for word in 0..WORDS_PER_HASH {
                first_row[val_i_word(i, word)] = F::from_canonical_u32((&self.leaves)[i][word]);
            }
        }

        for level in 0..(TREE_DEPTH - 1) {
            for _ in 0..level_width(level + 1) {
                let ([curr_row, next_row], idx, level_idx, level) = self.get_next_window();
                Self::gen_misc(curr_row, idx, level_idx, level, max_rows);

                let curr_level_width = level_width(level);
                let half_level_width = curr_level_width / 2;
                let level_done = level_idx == half_level_width - 1;

                // load in current hash's inputs
                for word in 0..WORDS_PER_HASH {
                    curr_row[hash_input_0_word(word)] = curr_row[val_i_word(0, word)]
                        + curr_row[HASH_IDX] * F::from_canonical_u64(1 << 32);
                    curr_row[hash_input_1_word(word)] = curr_row[val_i_word(1, word)]
                        + curr_row[HASH_IDX] * F::from_canonical_u64(1 << 32);
                }

                // compute output hash to be looked up
                let mut left = [0u32; 8];
                let mut right = [0u32; 8];
                for word in 0..WORDS_PER_HASH {
                    let left_with_idx = curr_row[hash_input_0_word(word)];
                    let right_with_idx = curr_row[hash_input_1_word(word)];
                    let left_field =
                        left_with_idx - curr_row[HASH_IDX] * F::from_canonical_u64(1 << 32);
                    let right_field =
                        right_with_idx - curr_row[HASH_IDX] * F::from_canonical_u64(1 << 32);

                    left[word] = left_field
                        .to_canonical_u64()
                        .try_into()
                        .expect("expected hash word to fit into u32");
                    right[word] = right_field
                        .to_canonical_u64()
                        .try_into()
                        .expect("expected hash word to fit into u32");
                }

                let output_hash = compress(left, right);
                for word in 0..WORDS_PER_HASH {
                    curr_row[hash_output_word(word)] = F::from_canonical_u32(output_hash[word])
                        + curr_row[HASH_IDX] * F::from_canonical_u64(1 << 32);
                }

                // shift vals
                // if we're not at the end of a level, shift var cols left by two and copy hash output to rightmost place
                // at the end of a level, we need to "unspread" the vals (see `mod.rs` for more explanation),
                // get them all to the left, and append the output hash
                if level_done && level < TREE_DEPTH - 1 {
                    let next_level_width = level_width(level + 1);
                    for i in 0..(next_level_width - 1) {
                        let shift_amount = get_level_end_shift(i, level);
                        for word in 0..WORDS_PER_HASH {
                            next_row[val_i_word(i, word)] =
                                curr_row[val_i_word(i + shift_amount, word)];
                        }
                    }

                    // append the hash from OUTPUT_COL
                    for word in 0..WORDS_PER_HASH {
                        next_row[val_i_word(next_level_width - 1, word)] = curr_row
                            [hash_output_word(word)]
                            - curr_row[HASH_IDX] * F::from_canonical_u64(1 << 32);
                    }

                    // zero the rest of the next row's var cols
                    for i in next_level_width..TREE_WIDTH {
                        for word in 0..WORDS_PER_HASH {
                            next_row[val_i_word(i, word)] = F::ZERO;
                        }
                    }
                } else {
                    for i in 0..(TREE_WIDTH - 2) {
                        for word in 0..WORDS_PER_HASH {
                            next_row[val_i_word(i, word)] = curr_row[val_i_word(i + 2, word)];
                        }
                    }

                    // load hash output into rightmost var cols of next row
                    for word in 0..WORDS_PER_HASH {
                        next_row[val_i_word(TREE_WIDTH - 1, word)] = curr_row
                            [hash_output_word(word)]
                            - curr_row[HASH_IDX] * F::from_canonical_u64(1 << 32);
                    }
                }

                // let mut curr_values = Vec::new();
                // for i in 0..TREE_WIDTH {
                //     for word in 0..WORDS_PER_HASH {
                //         curr_values.push(curr_row[val_i_word(i, word)]);
                //     }
                // }
                // println!("curr_values: {:?}", curr_values);
            }
        }

        let last_row = &self.trace.0[self.idx];
        let mut root = [0; WORDS_PER_HASH];
        for word in 0..WORDS_PER_HASH {
            root[word] = last_row[val_i_word(0, word)]
                .to_canonical_u64()
                .try_into()
                .expect("expected hash word to fit in u32");
        }

        (root, self.get_pis(root))
    }

    fn get_pis(&self, root: [u32; WORDS_PER_HASH]) -> [F; NUM_PUBLIC_INPUTS] {
        let mut pis = [F::ZERO; NUM_PUBLIC_INPUTS];
        for i in 0..TREE_WIDTH {
            for word in 0..WORDS_PER_HASH {
                pis[pi_leaf_i_word(i, word)] = F::from_canonical_u32((&self.leaves)[i][word]);
            }
        }

        for word in 0..WORDS_PER_HASH {
            pis[pi_root_word(word)] = F::from_canonical_u32(root[word]);
        }

        pis
    }

    pub fn into_polynomial_values(self) -> Vec<PolynomialValues<F>> {
        trace_rows_to_poly_values(self.trace.0)
    }
}

#[cfg(test)]
mod tests {
    use plonky2_field::goldilocks_field::GoldilocksField;

    use super::*;

    type F = GoldilocksField;

    fn merkle_root(leaves: &[[u32; 8]]) -> [u32; 8] {
        assert!(leaves.len() > 0);
        assert!(is_power_of_two(leaves.len() as u64));

        if leaves.len() == 2 {
            return compress(leaves[0], leaves[1]);
        } else {
            let half_len = leaves.len() / 2;
            let left = merkle_root(&leaves[0..half_len]);
            let right = merkle_root(&leaves[half_len..]);
            compress(left, right)
        }
    }

    #[test]
    fn test_build_tree() {
        let mut leaves = [[0; 8]; TREE_WIDTH];
        for i in 0..TREE_WIDTH {
            for word in 0..WORDS_PER_HASH {
                leaves[i][word] = (i * WORDS_PER_HASH + word) as u32;
            }
        }

        // compute tree without generator
        let correct_root = merkle_root(&leaves[..]);

        // compute trace with generator
        let mut generator = TreeTraceGenerator::<F>::new(16, leaves);
        let (_root, pis) = generator.gen();
        let mut root = [0u32; 8];
        for i in 0..WORDS_PER_HASH {
            root[i] = pis[pi_root_word(i)].to_canonical_u64().try_into().expect("expected pi_root_word to fit in u32");
        }

        assert_eq!(root, correct_root);
    }
}
