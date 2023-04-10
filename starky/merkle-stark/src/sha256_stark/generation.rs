#![allow(clippy::needless_range_loop)]

use core::convert::TryInto;
use std::iter::once;
use arrayref::{array_mut_ref, array_ref};
use plonky2::field::{polynomial::PolynomialValues, types::Field};

use super::constants::{HASH_IV, ROUND_CONSTANTS};
use super::layout::*;
use crate::util::trace_rows_to_poly_values;

fn is_power_of_two(n: u64) -> bool {
    n & (n - 1) == 0
}

#[repr(transparent)]
pub struct Sha2Trace<F: Field>(Vec<[F; NUM_COLS]>);

impl<F: Field> Sha2Trace<F> {
    pub fn new(max_rows: usize) -> Sha2Trace<F> {
        assert!(
            is_power_of_two(max_rows as u64),
            "max_rows must be a power of two"
        );
        Sha2Trace(vec![[F::ZERO; NUM_COLS]; max_rows])
    }
}

pub struct Sha2TraceGenerator<F: Field> {
    trace: Sha2Trace<F>,
    hash_idx: usize,
    left_input: [u32; 8],
    right_input: [u32; 8],
    step: usize,
}

impl<F: Field> Sha2TraceGenerator<F> {
    pub fn new(max_rows: usize) -> Sha2TraceGenerator<F> {
        Sha2TraceGenerator {
            trace: Sha2Trace::new(max_rows),
            hash_idx: 1, // hash_idx is 1-indexed
            left_input: [0; 8],
            right_input: [0; 8],
            step: 0,
        }
    }

    fn max_rows(&self) -> usize {
        self.trace.0.len()
    }

    fn curr_row_idx(&self) -> usize {
        (self.hash_idx - 1) * NUM_STEPS_PER_HASH + self.step
    }

    fn get_next_window(&mut self) -> (&mut [[F; NUM_COLS]; 2], usize, usize) {
        let idx = self.curr_row_idx();
        assert!(idx < self.max_rows(), "get_next_window exceeded MAX_ROWS");

        let hash_idx = self.hash_idx;
        let step = self.step;
        self.step += 1;

        (array_mut_ref![self.trace.0, idx, 2], hash_idx, step)
    }

    fn get_next_row(&mut self) -> (&mut [F; NUM_COLS], usize, usize) {
        let idx = self.curr_row_idx();
        assert!(idx < self.max_rows(), "get_next_window exceeded MAX_ROWS");

        let hash_idx = self.hash_idx;
        let step = self.step;
        self.step += 1;

        (&mut self.trace.0[idx], hash_idx, step)
    }

    // returns wi
    fn gen_msg_schedule(next_row: &mut [F; NUM_COLS], w15: u32, w2: u32, w16: u32, w7: u32) -> u32 {
        let mut xor_tmp_0 = rotr(w15, 7) ^ rotr(w15, 18);
        let mut s0 = xor_tmp_0 ^ (w15 >> 3);

        let mut xor_tmp_1 = rotr(w2, 17) ^ rotr(w2, 19);
        let mut s1 = xor_tmp_1 ^ (w2 >> 10);

        let wi = w16.wrapping_add(s0).wrapping_add(w7).wrapping_add(s1);
        let wi_u64 = w16 as u64 + s0 as u64 + w7 as u64 + s1 as u64;
        let quotient = wi_u64 / (1 << 32);

        for bit in 0..29 {
            next_row[xor_tmp_0_bit(bit)] = F::from_canonical_u32(xor_tmp_0 & 1);
            xor_tmp_0 >>= 1;
        }

        for bit in 0..22 {
            next_row[xor_tmp_1_bit(bit)] = F::from_canonical_u32(xor_tmp_1 & 1);
            xor_tmp_1 >>= 1;
        }
        for bit in 0..32 {
            next_row[little_s0_bit(bit)] = F::from_canonical_u32(s0 & 1);
            next_row[little_s1_bit(bit)] = F::from_canonical_u32(s1 & 1);

            s0 >>= 1;
            s1 >>= 1;
        }

        next_row[WI_FIELD] = F::from_canonical_u64(wi_u64);
        next_row[WI_QUOTIENT] = F::from_canonical_u64(quotient);

        wi
    }

    // returns new (abcd, efgh)
    fn gen_round_fn(
        curr_row: &mut [F; NUM_COLS],
        next_row: &mut [F; NUM_COLS],
        wi: u32,
        ki: u32,
        abcd: [u32; 4],
        efgh: [u32; 4],
    ) -> ([u32; 4], [u32; 4]) {
        let mut xor_tmp_2 = rotr(efgh[0], 6) ^ rotr(efgh[0], 11);
        let mut s1 = xor_tmp_2 ^ rotr(efgh[0], 25);
        let mut ch = (efgh[0] & efgh[1]) ^ ((!efgh[0]) & efgh[2]);
        let mut e_and_f = efgh[0] & efgh[1];
        let mut not_e_and_g = (!efgh[0]) & efgh[2];
        let mut xor_tmp_3 = rotr(abcd[0], 2) ^ rotr(abcd[0], 13);
        let mut s0 = xor_tmp_3 ^ rotr(abcd[0], 22);
        let mut xor_tmp_4 = (abcd[0] & abcd[1]) ^ (abcd[0] & abcd[2]);
        let mut maj = xor_tmp_4 ^ (abcd[1] & abcd[2]);
        let mut a_and_b = abcd[0] & abcd[1];
        let mut a_and_c = abcd[0] & abcd[2];
        let mut b_and_c = abcd[1] & abcd[2];

        let temp1_u32 = efgh[3]
            .wrapping_add(s1)
            .wrapping_add(ch)
            .wrapping_add(ki)
            .wrapping_add(wi);
        let temp2_u32 = s0.wrapping_add(maj);
        let temp1_u64 = efgh[3] as u64 + s1 as u64 + ch as u64 + ki as u64 + wi as u64;
        let temp2_u64 = s0 as u64 + maj as u64;

        for bit in 0..32 {
            curr_row[xor_tmp_2_bit(bit)] = F::from_canonical_u32(xor_tmp_2 & 1);
            curr_row[xor_tmp_3_bit(bit)] = F::from_canonical_u32(xor_tmp_3 & 1);
            curr_row[xor_tmp_4_bit(bit)] = F::from_canonical_u32(xor_tmp_4 & 1);

            curr_row[big_s1_bit(bit)] = F::from_canonical_u32(s1 & 1);
            curr_row[big_s0_bit(bit)] = F::from_canonical_u32(s0 & 1);
            curr_row[ch_bit(bit)] = F::from_canonical_u32(ch & 1);
            curr_row[maj_bit(bit)] = F::from_canonical_u32(maj & 1);

            curr_row[e_and_f_bit(bit)] = F::from_canonical_u32(e_and_f & 1);
            curr_row[not_e_and_g_bit(bit)] = F::from_canonical_u32(not_e_and_g & 1);
            curr_row[a_and_b_bit(bit)] = F::from_canonical_u32(a_and_b & 1);
            curr_row[a_and_c_bit(bit)] = F::from_canonical_u32(a_and_c & 1);
            curr_row[b_and_c_bit(bit)] = F::from_canonical_u32(b_and_c & 1);

            xor_tmp_2 >>= 1;
            xor_tmp_3 >>= 1;
            xor_tmp_4 >>= 1;
            s1 >>= 1;
            s0 >>= 1;
            ch >>= 1;
            maj >>= 1;
            e_and_f >>= 1;
            not_e_and_g >>= 1;
            a_and_b >>= 1;
            a_and_c >>= 1;
            b_and_c >>= 1;
        }

        let (mut abcd, mut efgh) = swap(abcd, efgh);

        let a_next_u64 = temp1_u64 + temp2_u64;
        let a_next_quotient = a_next_u64 / (1 << 32);
        let e_next_u64 = efgh[0] as u64 + temp1_u64;
        let e_next_quotient = e_next_u64 / (1 << 32);

        abcd[0] = temp1_u32.wrapping_add(temp2_u32);
        efgh[0] = efgh[0].wrapping_add(temp1_u32);

        let res = (abcd, efgh);

        curr_row[A_NEXT_FIELD] = F::from_canonical_u64(a_next_u64);
        curr_row[A_NEXT_QUOTIENT] = F::from_canonical_u64(a_next_quotient);
        curr_row[E_NEXT_FIELD] = F::from_canonical_u64(e_next_u64);
        curr_row[E_NEXT_QUOTIENT] = F::from_canonical_u64(e_next_quotient);

        for bit in 0..32 {
            next_row[a_bit(bit)] = F::from_canonical_u32(abcd[0] & 1);
            next_row[b_bit(bit)] = F::from_canonical_u32(abcd[1] & 1);
            next_row[c_bit(bit)] = F::from_canonical_u32(abcd[2] & 1);
            next_row[e_bit(bit)] = F::from_canonical_u32(efgh[0] & 1);
            next_row[f_bit(bit)] = F::from_canonical_u32(efgh[1] & 1);
            next_row[g_bit(bit)] = F::from_canonical_u32(efgh[2] & 1);

            abcd[0] >>= 1;
            abcd[1] >>= 1;
            abcd[2] >>= 1;
            efgh[0] >>= 1;
            efgh[1] >>= 1;
            efgh[2] >>= 1;
        }

        next_row[D_COL] = F::from_canonical_u32(abcd[3]);
        next_row[H_COL] = F::from_canonical_u32(efgh[3]);

        res
    }

    // fills in stuff the other fns don't at each row
    fn gen_misc(curr_row: &mut [F; NUM_COLS], step: usize, hash_idx: usize) {
        curr_row[HASH_IDX] = F::from_canonical_u64(hash_idx as u64);

        for i in 0..NUM_STEPS_PER_HASH {
            curr_row[step_bit(i)] = F::ZERO;
        }

        curr_row[step_bit(step)] = F::ONE;

        curr_row[INPUT_FILTER] = F::ZERO;
        curr_row[OUTPUT_FILTER] = F::ZERO;
    }

    fn gen_keep_his_same(curr_row: &mut [F; NUM_COLS], next_row: &mut [F; NUM_COLS]) {
        for i in 0..8 {
            next_row[h_i(i)] = curr_row[h_i(i)]
        }
    }
    // returns wis, abcd, efgh
    fn gen_phase_0(&mut self, his: [u32; 8]) -> ([u32; 16], [u32; 4], [u32; 4]) {
        let left_input = self.left_input;
        let right_input = self.right_input;
        let mut abcd = *array_ref![his, 0, 4];
        let mut efgh = *array_ref![his, 4, 4];

        let mut wis = [0; 16];
        (&mut wis)[..8].copy_from_slice(&left_input);
        (&mut wis)[8..].copy_from_slice(&right_input);
        wis = rotl_wis(wis);

        // left inputs
        for i in 0..16 {
            let ([curr_row, next_row], hash_idx, step) = self.get_next_window();
            Self::gen_misc(curr_row, step, hash_idx);

            if i == 0 {
                let mut abcd = abcd;
                let mut efgh = efgh;
                for bit in 0..32 {
                    curr_row[a_bit(bit)] = F::from_canonical_u32(abcd[0] & 1);
                    curr_row[b_bit(bit)] = F::from_canonical_u32(abcd[1] & 1);
                    curr_row[c_bit(bit)] = F::from_canonical_u32(abcd[2] & 1);
                    curr_row[e_bit(bit)] = F::from_canonical_u32(efgh[0] & 1);
                    curr_row[f_bit(bit)] = F::from_canonical_u32(efgh[1] & 1);
                    curr_row[g_bit(bit)] = F::from_canonical_u32(efgh[2] & 1);

                    abcd[0] >>= 1;
                    abcd[1] >>= 1;
                    abcd[2] >>= 1;
                    efgh[0] >>= 1;
                    efgh[1] >>= 1;
                    efgh[2] >>= 1;
                }
                
                curr_row[D_COL] = F::from_canonical_u32(abcd[3]);
                curr_row[H_COL] = F::from_canonical_u32(efgh[3]);

                // set his to IV
                for j in 0..8 {
                    curr_row[h_i(j)] = F::from_canonical_u32(HASH_IV[j]);
                }

                // load inputs
                for j in 0..8 {
                    curr_row[input_i(j)] = F::from_canonical_u32(left_input[j])
                        + F::from_canonical_u64(hash_idx as u64) * F::from_canonical_u64(1 << 32);

                    curr_row[input_i(j + 8)] = F::from_canonical_u32(right_input[j])
                        + F::from_canonical_u64(hash_idx as u64) * F::from_canonical_u64(1 << 32);
                }

                // load rotated wis
                Self::assign_wis(wis, curr_row);

                // set input filter to 1
                curr_row[INPUT_FILTER] = F::ONE;
            }

            Self::gen_keep_his_same(curr_row, next_row);

            let ki = ROUND_CONSTANTS[i];
            curr_row[KI] = F::from_canonical_u32(ki);
            (abcd, efgh) = Self::gen_round_fn(curr_row, next_row, wis[15], ki, abcd, efgh);

            if i == 15 {
                let w16 = wis[0];
                wis = shift_wis(wis);
                let wi = Self::gen_msg_schedule(next_row, wis[0], wis[13], w16, wis[8]);
                wis[15] = wi;
            } else {
                wis = rotl_wis(wis);
            }

            Self::assign_wis(wis, next_row);
        }

        (wis, abcd, efgh)
    }

    fn assign_wis(wis: [u32; 16], row: &mut [F; NUM_COLS]) {
        for i in 0..16 {
            match i {
                15 => {
                    let mut wi = wis[15];
                    for bit in 0..32 {
                        row[wi_bit(bit)] = F::from_canonical_u32(wi & 1);
                        wi >>= 1;
                    }
                }
                13 => {
                    let mut wi_minus_2 = wis[13];
                    for bit in 0..32 {
                        row[wi_minus_2_bit(bit)] = F::from_canonical_u32(wi_minus_2 & 1);
                        wi_minus_2 >>= 1;
                    }
                },
                0 => {
                    let mut wi_minus_15 = wis[0];
                    for bit in 0..32 {
                        row[wi_minus_15_bit(bit)] = F::from_canonical_u32(wi_minus_15 & 1);
                        wi_minus_15 >>= 1;
                    }
                },
                14 | 1..=12 => {
                    row[wi_field(i)] = F::from_canonical_u32(wis[i])
                }
                _ => unreachable!()
            }
        }
    }

    // returns wis, abcd, efgh, his
    fn gen_phase_1(
        &mut self,
        mut wis: [u32; 16],
        mut abcd: [u32; 4],
        mut efgh: [u32; 4],
        mut his: [u32; 8],
    ) -> ([u32; 16], [u32; 4], [u32; 4], [u32; 8]) {
        for i in 0..48 {
            let ([curr_row, next_row], hash_idx, step) = self.get_next_window();
            Self::gen_misc(curr_row, step, hash_idx);

            let ki = ROUND_CONSTANTS[i + 16];
            curr_row[KI] = F::from_canonical_u32(ki);
            (abcd, efgh) = Self::gen_round_fn(curr_row, next_row, wis[15], ki, abcd, efgh);

            let w16 = wis[0];
            wis = shift_wis(wis);

            if i != 47 {
                Self::gen_keep_his_same(curr_row, next_row);
                let wi = Self::gen_msg_schedule(next_row, wis[0], wis[13], w16, wis[8]);
                wis[15] = wi;
            }

            // update his during last row
            if i == 47 {
                for j in 0..4 {
                    let hj_next_u64 = his[j] as u64 + abcd[j] as u64;
                    let hj_next_quotient = hj_next_u64 / (1 << 32);

                    his[j] = his[j].wrapping_add(abcd[j]);
                    curr_row[h_i_next_field(j)] = F::from_canonical_u64(hj_next_u64);
                    curr_row[h_i_next_quotient(j)] = F::from_canonical_u64(hj_next_quotient);
                    next_row[h_i(j)] = F::from_canonical_u32(his[j]);
                }

                for j in 0..4 {
                    let hj_next_u64 = his[j + 4] as u64 + efgh[j] as u64;
                    let hj_next_quotient = hj_next_u64 / (1 << 32);

                    his[j + 4] = his[j + 4].wrapping_add(efgh[j]);
                    curr_row[h_i_next_field(j + 4)] = F::from_canonical_u64(hj_next_u64);
                    curr_row[h_i_next_quotient(j + 4)] = F::from_canonical_u64(hj_next_quotient);
                    next_row[h_i(j + 4)] = F::from_canonical_u32(his[j + 4]);
                }
            }

            Self::assign_wis(wis, next_row)
        }

        (wis, abcd, efgh, his)
    }

    fn gen_last_step(&mut self, his: [u32; 8]) {
        let (curr_row, hash_idx, step) = self.get_next_row();
        Self::gen_misc(curr_row, step, hash_idx);
        for i in 0..8 {
            curr_row[output_i(i)] = F::from_canonical_u32(his[i])
                + F::from_canonical_u64(hash_idx as u64) * F::from_canonical_u64(1 << 32);
        }
        curr_row[OUTPUT_FILTER] = F::ONE;
    }

    pub fn gen_hash(&mut self, left_input: [u32; 8], right_input: [u32; 8]) -> [u32; 8] {
        self.left_input = left_input;
        self.right_input = right_input;

        let his = HASH_IV;
        let (wis, abcd, efgh) = self.gen_phase_0(his);
        let (_wis, _abcd, _efgh, his) = self.gen_phase_1(wis, abcd, efgh, his);
        self.gen_last_step(his);

        self.hash_idx += 1;
        self.step = 0;
        his
    }

    pub fn into_polynomial_values(self) -> Vec<PolynomialValues<F>> {
        trace_rows_to_poly_values(self.trace.0)
    }
}

pub fn to_u32_array_be<const N: usize>(block: [u8; N * 4]) -> [u32; N] {
    let mut block_u32 = [0; N];
    for (o, chunk) in block_u32.iter_mut().zip(block.chunks_exact(4)) {
        *o = u32::from_be_bytes(chunk.try_into().unwrap());
    }
    block_u32
}

#[inline(always)]
fn shift_wis(mut wis: [u32; 16]) -> [u32; 16] {
    for i in 0..15 {
        wis[i] = wis[i + 1];
    }
    wis[15] = 0;
    wis
}

#[inline(always)]
fn rotl_wis(wis: [u32; 16]) -> [u32; 16] {
    let mut res = wis;
    for i in 0..16 {
        res[i] = wis[(i + 1) % 16];
    }
    res
}

#[inline(always)]
fn swap(abcd: [u32; 4], efgh: [u32; 4]) -> ([u32; 4], [u32; 4]) {
    (
        [0, abcd[0], abcd[1], abcd[2]],
        [abcd[3], efgh[0], efgh[1], efgh[2]],
    )
}

#[inline(always)]
fn rotr(x: u32, n: u32) -> u32 {
    x.rotate_right(n)
}

#[cfg(test)]
mod tests {
    use generic_array::{typenum::U64, GenericArray};
    use plonky2_field::goldilocks_field::GoldilocksField;
    use sha2::compress256;

    use super::*;

    type F = GoldilocksField;

    #[test]
    fn test_hash_of_zero() {
        let block = [0u8; 64];
        let block_arr = GenericArray::<u8, U64>::from(block);
        let mut state = HASH_IV;
        compress256(&mut state, &[block_arr]);

        let left_input = [0u32; 8];
        let right_input = [0u32; 8];
        let mut generator = Sha2TraceGenerator::<F>::new(128);

        let his = generator.gen_hash(left_input, right_input);

        assert_eq!(his, state);
    }

    #[test]
    fn test_hash_of_something() {
        let mut block = [0u8; 64];
        for i in 0..64 {
            block[i] = i as u8;
        }

        let block_arr = GenericArray::<u8, U64>::from(block);
        let mut state = HASH_IV;
        compress256(&mut state, &[block_arr]);

        let block: [u32; 16] = to_u32_array_be(block);
        let left_input = *array_ref![block, 0, 8];
        let right_input = *array_ref![block, 8, 8];
        let mut generator = Sha2TraceGenerator::<F>::new(128);

        let his = generator.gen_hash(left_input, right_input);
        assert_eq!(his, state);
    }
}
