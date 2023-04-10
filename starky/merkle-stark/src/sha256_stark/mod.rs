use std::iter::once;
use std::marker::PhantomData;

use itertools::Itertools;
use plonky2::field::{polynomial::PolynomialValues, types::Field};
use plonky2::iop::ext_target::ExtensionTarget;
use plonky2::{
    field::{
        extension::{Extendable, FieldExtension},
        packed::PackedField,
    },
    hash::hash_types::RichField,
    plonk::circuit_builder::CircuitBuilder,
};
use plonky2_util::log2_ceil;

use crate::{
    constraint_consumer::{ConstraintConsumer, RecursiveConstraintConsumer},
    stark::Stark,
    vars::{StarkEvaluationTargets, StarkEvaluationVars},
};

pub mod constants;
pub mod generation;
pub mod layout;

use generation::{to_u32_array_be, Sha2TraceGenerator};
use layout::*;

use self::constants::{HASH_IV, ROUND_CONSTANTS};

/// compute field_representation of a sequence of 32 bits interpreted big-endian u32 of a specific element of an trace array
macro_rules! bit_decomp_32_at_idx {
    ($row:expr, $idx:expr, $col_fn:ident, $f:ty, $p:ty) => {
        ((0..32).fold(<$p>::ZEROS, |acc, i| {
            acc + $row[$col_fn($idx, i)] * <$f>::from_canonical_u64(1 << i)
        }))
    };
}

/// compute field_representation of a sequence of 32 bits interpreted big-endian u32
macro_rules! bit_decomp_32 {
    ($row:expr, $col_fn:ident, $f:ty, $p:ty) => {
        ((0..32).fold(<$p>::ZEROS, |acc, i| {
            acc + $row[$col_fn(i)] * <$f>::from_canonical_u64(1 << i)
        }))
    };
}

// compute field_representation of a sequence of 32 bits interpreted big-endian u32 of a specific element of an trace array, in a circuit
macro_rules! bit_decomp_32_at_idx_circuit {
    ($builder:expr, $row:expr, $idx:expr, $col_fn:ident, $f:ty) => {{
        let addends = ((0..32).map(|i| {
            $builder.mul_const_extension(<$f>::from_canonical_u64(1 << i), $row[$col_fn($idx, i)])
        }))
        .collect::<Vec<_>>();

        $builder.add_many_extension(addends)
    }};
}

/// compute field_representation of a sequence of 32 bits interpreted big-endian u32, in a circuit
macro_rules! bit_decomp_32_circuit {
    ($builder:expr, $row:expr, $col_fn:ident, $f:ty) => {{
        let addends = ((0..32).map(|i| {
            $builder.mul_const_extension(<$f>::from_canonical_u64(1 << i), $row[$col_fn(i)])
        }))
        .collect::<Vec<_>>();

        $builder.add_many_extension(addends)
    }};
}

/// Computes the arithmetic generalization of `xor(x, y)`, i.e. `x + y - 2 x y`.
pub(crate) fn xor_gen<P: PackedField>(x: P, y: P) -> P {
    x + y - x * y.doubles()
}

/// Computes the arithmetic generalization of `xor(x, y)`, i.e. `x + y - 2 x y`, in a circuit
pub(crate) fn xor_gen_circuit<F: RichField + Extendable<D>, const D: usize>(
    builder: &mut CircuitBuilder<F, D>,
    x: ExtensionTarget<D>,
    y: ExtensionTarget<D>,
) -> ExtensionTarget<D> {
    let lhs = builder.add_extension(x, y);
    let rhs = builder.mul_extension_with_const(F::TWO, x, y);
    builder.sub_extension(lhs, rhs)
}

#[derive(Copy, Clone)]
pub struct Sha2CompressionStark<F: RichField + Extendable<D>, const D: usize> {
    _phantom: PhantomData<F>,
}

impl<F: RichField + Extendable<D>, const D: usize> Sha2CompressionStark<F, D> {
    pub fn new() -> Self {
        Self::default()
    }
}

impl<F: RichField + Extendable<D>, const D: usize> Default for Sha2CompressionStark<F, D> {
    fn default() -> Self {
        Self {
            _phantom: PhantomData,
        }
    }
}

impl<F: RichField + Extendable<D>, const D: usize> Stark<F, D> for Sha2CompressionStark<F, D> {
    const COLUMNS: usize = NUM_COLS;
    const PUBLIC_INPUTS: usize = 0;

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

        // set hash idx to 1 at the start. hash_idx should be 1-indexed.
        yield_constr.constraint_first_row(P::ONES - curr_row[HASH_IDX]);

        // set his to initial values at start of hash
        let is_hash_start = curr_row[step_bit(0)];
        for i in 0..8 {
            // degree 3
            yield_constr.constraint(
                is_hash_start * (curr_row[h_i(i)] - FE::from_canonical_u32(HASH_IV[i])),
            );
        }

        // ensure his stay the same outside last two rows of hash
        let his_should_change = next_row[step_bit(64)] + curr_row[step_bit(64)];
        for i in 0..8 {
            // degree 2
            yield_constr.constraint_transition(
                (P::ONES - his_should_change) * (next_row[h_i(i)] - curr_row[h_i(i)]),
            );
        }

        let next_is_phase_0: P = (0..16).map(|i| next_row[step_bit(i)]).sum();
        let next_is_phase_1: P = (16..64).map(|i| next_row[step_bit(i)]).sum();
        let next_is_not_padding = next_is_phase_0 + next_is_phase_1 + next_row[step_bit(64)];

        // increment hash idx if we're at the last step and next isn't padding
        let is_last_step = curr_row[step_bit(64)];
        let transition_to_next_hash = is_last_step * next_is_not_padding;
        // degree 3
        yield_constr.constraint_transition(
            transition_to_next_hash * (P::ONES - (next_row[HASH_IDX] - curr_row[HASH_IDX])),
        );

        // otherwise ensure hash idx stays the same unless next row is padding
        // degree 3
        yield_constr.constraint_transition(
            (P::ONES - is_last_step)
                * next_is_not_padding
                * (next_row[HASH_IDX] - curr_row[HASH_IDX]),
        );

        // load input into wis rotated left by one at start
        // wi
        let decomp = bit_decomp_32!(curr_row, wi_bit, FE, P)
            + curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32);
        yield_constr.constraint(is_hash_start * (decomp - curr_row[input_i(0)]));

        // wi minus 2
        let decomp = bit_decomp_32!(curr_row, wi_minus_2_bit, FE, P)
            + curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32);
        yield_constr.constraint(is_hash_start * (decomp - curr_row[input_i(14)]));

        // wi minus 15
        let decomp = bit_decomp_32!(curr_row, wi_minus_15_bit, FE, P)
            + curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32);
        yield_constr.constraint(is_hash_start * (decomp - curr_row[input_i(1)]));

        // the rest
        for i in (1..=12).chain(once(14)) {
            let decomp =
                curr_row[wi_field(i)] + curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32);

            yield_constr.constraint(is_hash_start * (decomp - curr_row[input_i(i + 1)]));
        }

        // set input filter to 1 iff we're at the start, zero otherwise
        yield_constr.constraint(is_hash_start * (P::ONES - curr_row[INPUT_FILTER]));
        yield_constr.constraint((P::ONES - is_hash_start) * curr_row[INPUT_FILTER]);

        // rotate wis when next step is phase 0 and we're not starting a new hash
        let rotate_wis = next_is_phase_0 * (P::ONES - next_row[step_bit(0)]);
        // shift wis left when next step is phase 1
        let shift_wis = next_is_phase_1;

        // wi next when rotating
        let c = bit_decomp_32!(next_row, wi_bit, FE, P)
            - bit_decomp_32!(curr_row, wi_minus_15_bit, FE, P);
        yield_constr.constraint_transition(rotate_wis * c);

        // wi_minus_2 next
        let decomp = bit_decomp_32!(next_row, wi_minus_2_bit, FE, P);
        yield_constr
            .constraint_transition((rotate_wis + shift_wis) * (decomp - curr_row[wi_field(14)]));

        // wi_minus_15 next
        let decomp = bit_decomp_32!(next_row, wi_minus_15_bit, FE, P);
        yield_constr
            .constraint_transition((rotate_wis + shift_wis) * (decomp - curr_row[wi_field(1)]));

        // wi_minus_1 next
        let decomp = bit_decomp_32!(curr_row, wi_bit, FE, P);
        yield_constr
            .constraint_transition((rotate_wis + shift_wis) * (next_row[wi_field(14)] - decomp));

        // wi_i_minus_3 next
        let decomp = bit_decomp_32!(curr_row, wi_minus_2_bit, FE, P);
        yield_constr
            .constraint_transition((rotate_wis + shift_wis) * (next_row[wi_field(12)] - decomp));

        for i in 1..12 {
            yield_constr.constraint_transition(
                (rotate_wis + shift_wis) * (next_row[wi_field(i)] - curr_row[wi_field(i + 1)]),
            );
        }

        // round fn in phase 0 or 1
        let is_phase_0_or_1: P = (0..64).map(|i| curr_row[step_bit(i)]).sum();

        // S1 := (e >>> 6) xor (e >>> 11) xor (e >>> 25)
        for bit in 0..32 {
            let computed_bit = xor_gen(
                curr_row[e_bit((bit + 6) % 32)],
                curr_row[e_bit((bit + 11) % 32)],
            );
            // degree 3
            yield_constr
                .constraint(is_phase_0_or_1 * (curr_row[xor_tmp_2_bit(bit)] - computed_bit));

            let computed_bit = xor_gen(
                curr_row[xor_tmp_2_bit(bit)],
                curr_row[e_bit((bit + 25) % 32)],
            );
            // degree 3
            yield_constr.constraint(is_phase_0_or_1 * (curr_row[big_s1_bit(bit)] - computed_bit));
        }

        // ch := (e and f) xor ((not e) and g)
        for bit in 0..32 {
            let computed_bit = curr_row[e_bit(bit)] * curr_row[f_bit(bit)];
            // degree 3
            yield_constr.constraint(is_phase_0_or_1 * (curr_row[e_and_f_bit(bit)] - computed_bit));

            let computed_bit = (P::ONES - curr_row[e_bit(bit)]) * curr_row[g_bit(bit)];
            // degree 3
            yield_constr
                .constraint(is_phase_0_or_1 * (curr_row[not_e_and_g_bit(bit)] - computed_bit));

            let computed_bit = xor_gen(curr_row[e_and_f_bit(bit)], curr_row[not_e_and_g_bit(bit)]);
            // degree 3
            yield_constr.constraint(is_phase_0_or_1 * (curr_row[ch_bit(bit)] - computed_bit));
        }

        // S0 := (a >>> 2) xor (a >>> 13) xor (a >>> 22)
        for bit in 0..32 {
            let computed_bit = xor_gen(
                curr_row[a_bit((bit + 2) % 32)],
                curr_row[a_bit((bit + 13) % 32)],
            );
            // degree 3
            yield_constr
                .constraint(is_phase_0_or_1 * (curr_row[xor_tmp_3_bit(bit)] - computed_bit));

            let computed_bit = xor_gen(
                curr_row[xor_tmp_3_bit(bit)],
                curr_row[a_bit((bit + 22) % 32)],
            );
            // degree 3
            yield_constr.constraint(is_phase_0_or_1 * (curr_row[big_s0_bit(bit)] - computed_bit));
        }

        // maj := (a and b) xor (a and c) xor (b and c)
        for bit in 0..32 {
            // degree 3
            yield_constr.constraint(
                is_phase_0_or_1
                    * (curr_row[a_and_b_bit(bit)] - curr_row[a_bit(bit)] * curr_row[b_bit(bit)]),
            );

            // degree 3
            yield_constr.constraint(
                is_phase_0_or_1
                    * (curr_row[a_and_c_bit(bit)] - curr_row[a_bit(bit)] * curr_row[c_bit(bit)]),
            );

            // degree 3
            yield_constr.constraint(
                is_phase_0_or_1
                    * (curr_row[b_and_c_bit(bit)] - curr_row[b_bit(bit)] * curr_row[c_bit(bit)]),
            );

            let computed_bit = xor_gen(curr_row[a_and_b_bit(bit)], curr_row[a_and_c_bit(bit)]);
            // degree 3
            yield_constr
                .constraint(is_phase_0_or_1 * (curr_row[xor_tmp_4_bit(bit)] - computed_bit));

            let computed_bit = xor_gen(curr_row[xor_tmp_4_bit(bit)], curr_row[b_and_c_bit(bit)]);
            // degree 3
            yield_constr.constraint(is_phase_0_or_1 * (curr_row[maj_bit(bit)] - computed_bit));
        }

        // set round constant
        for step in 0..64 {
            // degree 2
            yield_constr.constraint(
                curr_row[step_bit(step)]
                    * (curr_row[KI] - FE::from_canonical_u32(ROUND_CONSTANTS[step])),
            );
        }

        // temp1 := h + S1 + ch + k[i] + w[i]
        // e := d + temp1
        let h_field = curr_row[H_COL];
        let big_s1_field = bit_decomp_32!(curr_row, big_s1_bit, FE, P);
        let ch_field = bit_decomp_32!(curr_row, ch_bit, FE, P);
        let wi_u32 = bit_decomp_32!(curr_row, wi_bit, FE, P);
        let temp1_minus_ki = h_field + big_s1_field + ch_field + wi_u32;

        let d_field = curr_row[D_COL];
        let e_u32_next = bit_decomp_32!(next_row, e_bit, FE, P);

        // degree 2
        yield_constr.constraint(
            is_phase_0_or_1 * (curr_row[E_NEXT_FIELD] - (d_field + temp1_minus_ki + curr_row[KI])),
        );
        // degree 3
        yield_constr.constraint_transition(
            is_phase_0_or_1
                * (curr_row[E_NEXT_FIELD]
                    - (e_u32_next + curr_row[E_NEXT_QUOTIENT] * FE::from_canonical_u64(1 << 32))),
        );

        // temp2 := S0 + maj
        // a := temp1 + temp2
        let s0_field = bit_decomp_32!(curr_row, big_s0_bit, FE, P);
        let maj_field = bit_decomp_32!(curr_row, maj_bit, FE, P);
        let temp2 = s0_field + maj_field;
        let a_u32_next = bit_decomp_32!(next_row, a_bit, FE, P);

        // degree 2
        yield_constr.constraint(
            is_phase_0_or_1 * (curr_row[A_NEXT_FIELD] - (temp2 + temp1_minus_ki + curr_row[KI])),
        );
        // degree 3
        yield_constr.constraint_transition(
            is_phase_0_or_1
                * (curr_row[A_NEXT_FIELD]
                    - (a_u32_next + curr_row[A_NEXT_QUOTIENT] * FE::from_canonical_u64(1 << 32))),
        );

        // update local vars when not in last step
        // h := g
        // g := f
        // f := e
        // d := c
        // c := b
        // b := a
        let decomp = bit_decomp_32!(curr_row, g_bit, FE, P);
        yield_constr.constraint_transition(is_phase_0_or_1 * (next_row[H_COL] - decomp));

        let decomp = bit_decomp_32!(curr_row, c_bit, FE, P);
        yield_constr.constraint_transition(is_phase_0_or_1 * (next_row[D_COL] - decomp));
        for bit in 0..32 {
            // degree 2
            yield_constr.constraint_transition(
                is_phase_0_or_1 * (next_row[g_bit(bit)] - curr_row[f_bit(bit)]),
            );
            yield_constr.constraint_transition(
                is_phase_0_or_1 * (next_row[f_bit(bit)] - curr_row[e_bit(bit)]),
            );
            yield_constr.constraint_transition(
                is_phase_0_or_1 * (next_row[c_bit(bit)] - curr_row[b_bit(bit)]),
            );
            yield_constr.constraint_transition(
                is_phase_0_or_1 * (next_row[b_bit(bit)] - curr_row[a_bit(bit)]),
            );
        }

        // update his in last step of phase 1
        let update_his = curr_row[step_bit(63)];
        let vars = [
            bit_decomp_32!(next_row, a_bit, FE, P),
            bit_decomp_32!(next_row, b_bit, FE, P),
            bit_decomp_32!(next_row, c_bit, FE, P),
            next_row[D_COL],
            bit_decomp_32!(next_row, e_bit, FE, P),
            bit_decomp_32!(next_row, f_bit, FE, P),
            bit_decomp_32!(next_row, g_bit, FE, P),
            next_row[H_COL],
        ];

        for i in 0..8 {
            // degree 2
            yield_constr.constraint_transition(
                update_his * (curr_row[h_i_next_field(i)] - (curr_row[h_i(i)] + vars[i])),
            );

            // degree 3
            yield_constr.constraint_transition(
                update_his
                    * (curr_row[h_i_next_field(i)]
                        - (next_row[h_i(i)]
                            + curr_row[h_i_next_quotient(i)] * FE::from_canonical_u64(1 << 32))),
            );
        }

        // set output to his during last step, 0 otherwise
        for i in 0..8 {
            // degree 2
            yield_constr.constraint(
                is_last_step
                    * (curr_row[output_i(i)]
                        - (curr_row[h_i(i)]
                            + curr_row[HASH_IDX] * FE::from_canonical_u64(1 << 32))),
            );
            yield_constr.constraint((P::ONES - is_last_step) * curr_row[output_i(i)])
        }

        // set output filter to 1 iff it's the last step, zero otherwise
        yield_constr.constraint(is_last_step * (P::ONES - curr_row[OUTPUT_FILTER]));
        yield_constr.constraint((P::ONES - is_last_step) * curr_row[OUTPUT_FILTER]);

        // message schedule to get next row's wi when next row is in phase 1

        let do_msg_schedule = next_is_phase_1;

        // s0 := (w[i-15] >>> 7) xor (w[i-15] >>> 18) xor (w[i-15] >> 3)
        for bit in 0..29 {
            let computed_bit = xor_gen(
                next_row[wi_minus_15_bit((bit + 7) % 32)],
                next_row[wi_minus_15_bit((bit + 18) % 32)],
            );
            // degree 3
            yield_constr.constraint_transition(
                do_msg_schedule * (next_row[xor_tmp_0_bit(bit)] - computed_bit),
            );

            let computed_bit = xor_gen(
                next_row[xor_tmp_0_bit(bit)],
                next_row[wi_minus_15_bit(bit + 3)],
            );
            // degree 3
            yield_constr.constraint_transition(
                do_msg_schedule * (next_row[little_s0_bit(bit)] - computed_bit),
            );
        }
        for bit in 29..32 {
            // we can ignore the second XOR in this case since it's with 0
            let computed_bit = xor_gen(
                next_row[wi_minus_15_bit((bit + 7) % 32)],
                next_row[wi_minus_15_bit((bit + 18) % 32)],
            );
            // degree 3
            yield_constr.constraint_transition(
                do_msg_schedule * (next_row[little_s0_bit(bit)] - computed_bit),
            );
        }

        // s1 := (w[i-2] >>> 17) xor (w[i-2] >>> 19) xor (w[i-2] >> 10)
        for bit in 0..22 {
            let computed_bit = xor_gen(
                next_row[wi_minus_2_bit((bit + 17) % 32)],
                next_row[wi_minus_2_bit((bit + 19) % 32)],
            );
            // degree 3
            yield_constr.constraint_transition(
                do_msg_schedule * (next_row[xor_tmp_1_bit(bit)] - computed_bit),
            );

            let computed_bit = xor_gen(
                next_row[xor_tmp_1_bit(bit)],
                next_row[wi_minus_2_bit(bit + 10)],
            );
            // degree 3
            yield_constr.constraint_transition(
                do_msg_schedule * (next_row[little_s1_bit(bit)] - computed_bit),
            );
        }
        for bit in 22..32 {
            // we can ignore the second XOR in this case since it's with 0
            let computed_bit = xor_gen(
                next_row[wi_minus_2_bit((bit + 17) % 32)],
                next_row[wi_minus_2_bit((bit + 19) % 32)],
            );
            // degree 3
            yield_constr.constraint_transition(
                do_msg_schedule * (next_row[little_s1_bit(bit)] - computed_bit),
            );
        }

        // w[i] := w[i-16] + s0 + w[i-7] + s1

        // degree 1
        let s0_field_computed = bit_decomp_32!(next_row, little_s0_bit, FE, P);
        let s1_field_computed = bit_decomp_32!(next_row, little_s1_bit, FE, P);
        let wi_minus_16_field = bit_decomp_32!(curr_row, wi_minus_15_bit, FE, P);
        let wi_minus_7_field = next_row[wi_field(8)];
        let wi = bit_decomp_32!(next_row, wi_bit, FE, P);

        // degree 2
        yield_constr.constraint_transition(
            do_msg_schedule
                * (next_row[WI_FIELD]
                    - (wi_minus_16_field
                        + s0_field_computed
                        + wi_minus_7_field
                        + s1_field_computed)),
        );
        // degree 3
        yield_constr.constraint(
            do_msg_schedule
                * (next_row[WI_FIELD]
                    - (wi + next_row[WI_QUOTIENT] * FE::from_canonical_u64(1 << 32))),
        );

        // set initial step bits to a 1 followed by NUM_STEPS_PER_HASH-1 0s
        yield_constr.constraint_first_row(P::ONES - curr_row[step_bit(0)]);
        for step in 1..NUM_STEPS_PER_HASH {
            yield_constr.constraint_first_row(curr_row[step_bit(step)]);
        }

        // inc step bits when next is not padding
        for bit in 0..NUM_STEPS_PER_HASH {
            // degree 3
            yield_constr.constraint_transition(
                next_is_not_padding
                    * (next_row[step_bit((bit + 1) % NUM_STEPS_PER_HASH)]
                        - curr_row[step_bit(bit)]),
            );
        }

        eval_bits_are_bits(curr_row, yield_constr);
    }

    fn eval_ext_circuit(
        &self,
        builder: &mut CircuitBuilder<F, D>,
        vars: StarkEvaluationTargets<D, { Self::COLUMNS }, { Self::PUBLIC_INPUTS }>,
        yield_constr: &mut RecursiveConstraintConsumer<F, D>,
    ) {

        // let curr_row = vars.local_values;
        // let next_row = vars.next_values;

        // let one_ext = builder.one_extension();
        // let zero_ext = builder.zero_extension();

        // let c = builder.sub_extension(one_ext, curr_row[HASH_IDX]);
        // yield_constr.constraint_first_row(builder, c);

        // let is_hash_start = curr_row[step_bit(0)];
        // let hash_ivs = (0..8).map(|i| builder.constant_extension(F::Extension::from_canonical_u32(HASH_IV[i]))).collect_vec();
        // for (i, &iv) in hash_ivs.iter().enumerate() {
        //     let mut c = builder.sub_extension(curr_row[h_i(i)], iv);
        //     c = builder.mul_extension(c, is_hash_start);
        //     yield_constr.constraint(builder, c);
        // }

        // let his_should_change = builder.add_extension(next_row[step_bit(64)], curr_row[step_bit(64)]);
        // for i in 0..8 {
        //     let lhs = builder.sub_extension(next_row[h_i(i)], curr_row[h_i(i)]);
        //     let rhs = builder.sub_extension(one_ext, his_should_change);
        //     let c = builder.mul_extension(lhs, rhs);
        //     yield_constr.constraint_transition(builder, c);
        // }

        // let next_is_phase_0 = builder.add_many_extension((0..16).map(|i| next_row[step_bit(i)]));
        // let next_is_phase_1 = builder.add_many_extension((16..64).map(|i| next_row[step_bit(i)]));
        // let next_is_not_padding = builder.add_many_extension(
        //     once(next_is_phase_0).chain(once(next_is_phase_1)).chain(once(next_row[step_bit(64)]))
        // );

        // let is_last_step = curr_row[step_bit(64)];
        // let transition_to_next_hash = builder.mul_extension(is_last_step, next_is_not_padding);

        // let hash_idx_diff = builder.sub_extension(next_row[HASH_IDX], curr_row[HASH_IDX]);
        // let mut c = builder.sub_extension(one_ext, hash_idx_diff);
        // c = builder.mul_extension(transition_to_next_hash, c);
        // yield_constr.constraint_transition(builder, c);

        // let mut c = builder.sub_extension(one_ext, is_last_step);
        // c = builder.mul_extension(next_is_not_padding, c);
        // c = builder.mul_extension(hash_idx_diff, c);
        // yield_constr.constraint_transition(builder, c);

        // for i in 0..16 {
        //     let decomp = bit_decomp_32_at_idx_circuit!(builder, curr_row, i, wi_bit, F);
        //     let mut c = builder.mul_const_add_extension(F::from_canonical_u64(1 << 32), curr_row[HASH_IDX], decomp);

        //     c = builder.sub_extension(decomp, curr_row[input_i(i + 1) % 16]);
        //     c = builder.mul_extension(is_hash_start, c);
        //     yield_constr.constraint(builder, c);
        // }

        // let mut c = builder.sub_extension(one_ext, curr_row[INPUT_FILTER]);
        // c = builder.mul_extension(is_hash_start, c);
        // yield_constr.constraint(builder, c);

        // let mut c = builder.sub_extension(one_ext, is_hash_start);
        // c = builder.mul_extension(c, curr_row[INPUT_FILTER]);
        // yield_constr.constraint(builder, c);

        // let next_is_not_step_0 = builder.sub_extension(one_ext, next_row[step_bit(0)]);
        // let rotate_wis = builder.mul_extension(next_is_phase_0, next_is_not_step_0);
        // for i in 0..16 {
        //     for bit in 0..32 {
        //         let mut c = builder.sub_extension(next_row[wi_bit(i, bit)], curr_row[wi_bit((i + 1) % 16, bit)]);
        //         c = builder.mul_extension(rotate_wis, c);
        //         yield_constr.constraint(builder, c);
        //     }
        // }

        // let shift_wis = next_is_phase_1;
        // for i in 0..15 {
        //     for bit in 0..32 {
        //         let mut c = builder.sub_extension(next_row[wi_bit(i, bit)], curr_row[wi_bit(i + 1, bit)]);
        //         c = builder.mul_extension(shift_wis, c);
        //         yield_constr.constraint(builder, c);
        //     }
        // }

        // let is_phase_0_or_1 = builder.add_many_extension((0..64).map(|i| curr_row[step_bit(i)]));

        // for bit in 0..32 {
        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         curr_row[e_bit(bit + 6) % 32],
        //         curr_row[e_bit(bit + 11) % 32]
        //     );
        //     let mut c = builder.sub_extension(curr_row[xor_tmp_2_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);

        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         curr_row[xor_tmp_2_bit(bit)],
        //         curr_row[e_bit(bit + 25) % 32]
        //     );

        //     let mut c = builder.sub_extension(curr_row[big_s1_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let computed_bit = builder.mul_extension(curr_row[e_bit(bit)], curr_row[f_bit(bit)]);
        //     let mut c = builder.sub_extension(curr_row[e_and_f_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);

        //     let mut computed_bit = builder.sub_extension(one_ext, curr_row[e_bit(bit)]);
        //     computed_bit = builder.mul_extension(computed_bit, curr_row[g_bit(bit)]);
        //     let mut c = builder.sub_extension(curr_row[not_e_and_g_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);

        //     let computed_bit = xor_gen_circuit(builder, curr_row[e_and_f_bit(bit)], curr_row[not_e_and_g_bit(bit)]);
        //     let mut c = builder.sub_extension(curr_row[ch_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         curr_row[a_bit(bit + 2) % 32],
        //         curr_row[a_bit(bit + 13) % 32]
        //     );
        //     let mut c = builder.sub_extension(curr_row[xor_tmp_3_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);

        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         curr_row[xor_tmp_3_bit(bit)],
        //         curr_row[a_bit((bit + 22) % 32)]
        //     );
        //     let mut c = builder.sub_extension(curr_row[big_s0_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.mul_extension(curr_row[a_bit(bit)], curr_row[b_bit(bit)]);
        //     c = builder.sub_extension(curr_row[a_and_b_bit(bit)], c);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);

        //     let mut c = builder.mul_extension(curr_row[a_bit(bit)], curr_row[c_bit(bit)]);
        //     c = builder.sub_extension(curr_row[a_and_c_bit(bit)], c);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);

        //     let mut c = builder.mul_extension(curr_row[b_bit(bit)], curr_row[c_bit(bit)]);
        //     c = builder.sub_extension(curr_row[b_and_c_bit(bit)], c);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);

        //     let computed_bit = xor_gen_circuit(builder, curr_row[a_and_b_bit(bit)], curr_row[a_and_c_bit(bit)]);
        //     let mut c = builder.sub_extension(curr_row[xor_tmp_4_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);

        //     let computed_bit = xor_gen_circuit(builder, curr_row[xor_tmp_4_bit(bit)], curr_row[b_and_c_bit(bit)]);
        //     let mut c = builder.sub_extension(curr_row[maj_bit(bit)], computed_bit);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint(builder, c);
        // }

        // for step in 0..64 {
        //     let mut c = builder.add_const_extension(curr_row[KI], -F::from_canonical_u32(ROUND_CONSTANTS[step]));
        //     c = builder.mul_extension(curr_row[step_bit(step)], c);
        //     yield_constr.constraint(builder, c);
        // }

        // let h_field = bit_decomp_32_circuit!(builder, curr_row, h_bit, F);
        // let big_s1_field = bit_decomp_32_circuit!(builder, curr_row, big_s1_bit, F);
        // let ch_field = bit_decomp_32_circuit!(builder, curr_row, ch_bit, F);
        // let wi_u32 = bit_decomp_32_at_idx_circuit!(builder, curr_row, 15, wi_bit, F);
        // let temp1_minus_ki = builder.add_many_extension(
        //     once(h_field).chain(once(big_s1_field)).chain(once(ch_field)).chain(once(wi_u32))
        // );

        // let d_field = bit_decomp_32_circuit!(builder, curr_row, d_bit, F);
        // let e_u32_next = bit_decomp_32_circuit!(builder, next_row, e_bit, F);

        // let mut c = builder.add_extension(d_field, temp1_minus_ki);
        // c = builder.add_extension(c, curr_row[KI]);
        // c = builder.sub_extension(curr_row[E_NEXT_FIELD], c);
        // c = builder.mul_extension(is_phase_0_or_1, c);
        // yield_constr.constraint(builder, c);

        // let mut c = builder.mul_const_add_extension(F::from_canonical_u64(1 << 32), curr_row[E_NEXT_QUOTIENT], e_u32_next);
        // c = builder.sub_extension(curr_row[E_NEXT_FIELD], c);
        // c = builder.mul_extension(is_phase_0_or_1, c);
        // yield_constr.constraint_transition(builder, c);

        // let s0_field = bit_decomp_32_circuit!(builder, curr_row, big_s0_bit, F);
        // let maj_field = bit_decomp_32_circuit!(builder, curr_row, maj_bit, F);
        // let temp2 = builder.add_extension(s0_field, maj_field);
        // let a_u32_next = bit_decomp_32_circuit!(builder, next_row, a_bit, F);

        // let mut c = builder.add_extension(temp1_minus_ki, curr_row[KI]);
        // c = builder.add_extension(temp2, c);
        // c = builder.sub_extension(curr_row[A_NEXT_FIELD], c);
        // c = builder.mul_extension(is_phase_0_or_1, c);
        // yield_constr.constraint(builder, c);

        // let mut c = builder.mul_const_add_extension(F::from_canonical_u64(1 << 32), curr_row[A_NEXT_QUOTIENT], a_u32_next);
        // c = builder.sub_extension(curr_row[A_NEXT_FIELD], c);
        // c = builder.mul_extension(is_phase_0_or_1, c);
        // yield_constr.constraint_transition(builder, c);

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(next_row[h_bit(bit)], curr_row[g_bit(bit)]);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint_transition(builder, c);

        //     let mut c = builder.sub_extension(next_row[g_bit(bit)], curr_row[f_bit(bit)]);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint_transition(builder, c);

        //     let mut c = builder.sub_extension(next_row[f_bit(bit)], curr_row[e_bit(bit)]);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint_transition(builder, c);

        //     let mut c = builder.sub_extension(next_row[d_bit(bit)], curr_row[c_bit(bit)]);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint_transition(builder, c);

        //     let mut c = builder.sub_extension(next_row[c_bit(bit)], curr_row[b_bit(bit)]);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint_transition(builder, c);

        //     let mut c = builder.sub_extension(next_row[b_bit(bit)], curr_row[a_bit(bit)]);
        //     c = builder.mul_extension(is_phase_0_or_1, c);
        //     yield_constr.constraint_transition(builder, c);
        // }

        // let update_his = curr_row[step_bit(63)];
        // let vars = [
        //     bit_decomp_32_circuit!(builder, next_row, a_bit, F),
        //     bit_decomp_32_circuit!(builder, next_row, b_bit, F),
        //     bit_decomp_32_circuit!(builder, next_row, c_bit, F),
        //     bit_decomp_32_circuit!(builder, next_row, d_bit, F),
        //     bit_decomp_32_circuit!(builder, next_row, e_bit, F),
        //     bit_decomp_32_circuit!(builder, next_row, f_bit, F),
        //     bit_decomp_32_circuit!(builder, next_row, g_bit, F),
        //     bit_decomp_32_circuit!(builder, next_row, h_bit, F),
        // ];

        // for i in 0..8 {
        //     let mut c = builder.add_extension(curr_row[h_i(i)], vars[i]);
        //     c = builder.sub_extension(curr_row[h_i_next_field(i)], c);
        //     c = builder.mul_extension(update_his, c);
        //     yield_constr.constraint_transition(builder, c);

        //     let mut c = builder.mul_const_add_extension(F::from_canonical_u64(1 << 32), curr_row[h_i_next_quotient(i)], next_row[h_i(i)]);
        //     c = builder.sub_extension(curr_row[h_i_next_field(i)], c);
        //     c = builder.mul_extension(update_his, c);
        //     yield_constr.constraint_transition(builder, c);
        // }

        // let is_not_last_step = builder.sub_extension(one_ext, is_last_step);

        // for i in 0..8 {
        //     let mut c = builder.mul_const_add_extension(F::from_canonical_u64(1 << 32), curr_row[HASH_IDX], curr_row[h_i(i)]);
        //     c = builder.sub_extension(curr_row[output_i(i)], c);
        //     c = builder.mul_extension(is_last_step, c);
        //     yield_constr.constraint(builder, c);

        //     let c = builder.mul_extension(is_not_last_step, curr_row[output_i(i)]);
        //     yield_constr.constraint(builder, c);
        // }

        // let mut c = builder.sub_extension(one_ext, curr_row[OUTPUT_FILTER]);
        // c = builder.mul_extension(is_last_step, c);
        // yield_constr.constraint(builder, c);

        // let c = builder.mul_extension(is_not_last_step, curr_row[OUTPUT_FILTER]);
        // yield_constr.constraint(builder, c);

        // let do_msg_schedule = next_is_phase_1;

        // for bit in 0..29 {
        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         next_row[wi_bit(0, (bit + 7) % 32)],
        //         next_row[wi_bit(0, (bit + 18) % 32)]
        //     );

        //     let mut c = builder.sub_extension(next_row[xor_tmp_0_bit(bit)], computed_bit);
        //     c = builder.mul_extension(do_msg_schedule, c);
        //     yield_constr.constraint_transition(builder, c);

        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         next_row[xor_tmp_0_bit(bit)],
        //         next_row[wi_bit(0, bit + 3)]
        //     );

        //     let mut c = builder.sub_extension(next_row[little_s0_bit(bit)], computed_bit);
        //     c = builder.mul_extension(do_msg_schedule, c);
        //     yield_constr.constraint_transition(builder, c);
        // }
        // for bit in 29..32 {
        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         next_row[wi_bit(0, (bit + 7) % 32)],
        //         next_row[wi_bit(0, (bit + 18) % 32)]
        //     );

        //     let mut c = builder.sub_extension(next_row[little_s0_bit(bit)], computed_bit);
        //     c = builder.mul_extension(do_msg_schedule, c);
        //     yield_constr.constraint_transition(builder, c);
        // }

        // for bit in 0..22 {
        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         next_row[wi_bit(13, (bit + 17) % 32)],
        //         next_row[wi_bit(13, (bit + 19) % 32)],
        //     );

        //     let mut c = builder.sub_extension(next_row[xor_tmp_1_bit(bit)], computed_bit);
        //     c = builder.mul_extension(do_msg_schedule, c);
        //     yield_constr.constraint_transition(builder, c);

        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         next_row[xor_tmp_1_bit(bit)],
        //         next_row[wi_bit(13, bit + 10)]
        //     );

        //     let mut c = builder.sub_extension(next_row[little_s1_bit(bit)], computed_bit);
        //     c = builder.mul_extension(do_msg_schedule, c);
        //     yield_constr.constraint_transition(builder, c);
        // }
        // for bit in 22..32 {
        //     let computed_bit = xor_gen_circuit(
        //         builder,
        //         next_row[wi_bit(13, (bit + 17) % 32)],
        //         next_row[wi_bit(13, (bit + 19) % 32)]
        //     );

        //     let mut c = builder.sub_extension(next_row[little_s1_bit(bit)], computed_bit);
        //     c = builder.mul_extension(do_msg_schedule, c);
        //     yield_constr.constraint_transition(builder, c);
        // }

        // let s0_field_computed = bit_decomp_32_circuit!(builder, next_row, little_s0_bit, F);
        // let s1_field_computed = bit_decomp_32_circuit!(builder, next_row, little_s1_bit, F);
        // let wi_minus_16_field_computed = bit_decomp_32_at_idx_circuit!(builder, curr_row, 0, wi_bit, F);
        // let wi_minus_7_field_computed = bit_decomp_32_at_idx_circuit!(builder, next_row, 8, wi_bit, F);
        // let wi = bit_decomp_32_at_idx_circuit!(builder, next_row, 15, wi_bit, F);

        // let mut c = builder.add_many_extension(
        //     once(wi_minus_16_field_computed).chain(once(s0_field_computed)).chain(once(wi_minus_7_field_computed)).chain(once(s1_field_computed))
        // );
        // c = builder.sub_extension(next_row[WI_FIELD], c);
        // c = builder.mul_extension(do_msg_schedule, c);
        // yield_constr.constraint_transition(builder, c);

        // let mut c = builder.mul_const_add_extension(F::from_canonical_u64(1 << 32), next_row[WI_QUOTIENT], wi);
        // c = builder.sub_extension(next_row[WI_FIELD], c);
        // c = builder.mul_extension(do_msg_schedule, c);
        // yield_constr.constraint(builder, c);

        // let c = builder.sub_extension(one_ext, curr_row[step_bit(0)]);
        // yield_constr.constraint(builder, c);

        // for step in 1..NUM_STEPS_PER_HASH {
        //     yield_constr.constraint_first_row(builder, curr_row[step_bit(step)]);
        // }

        // for bit in 0..NUM_STEPS_PER_HASH {
        //     let mut c = builder.sub_extension(next_row[step_bit((bit + 1) % NUM_STEPS_PER_HASH)], curr_row[step_bit(bit)]);
        //     c = builder.mul_extension(next_is_not_padding, c);
        //     yield_constr.constraint_transition(builder, c);
        // }

        // for bit in 0..NUM_STEPS_PER_HASH {
        //     let mut c = builder.sub_extension(one_ext, curr_row[step_bit(bit)]);
        //     c = builder.mul_extension(curr_row[step_bit(bit)], c);
        //     yield_constr.constraint(builder, c);
        // }

        // for i in 0..NUM_WIS {
        //     for bit in 0..32 {
        //         let mut c = builder.sub_extension(one_ext, curr_row[wi_bit(i, bit)]);
        //         c = builder.mul_extension(c, curr_row[wi_bit(i, bit)]);
        //         yield_constr.constraint(builder, c);
        //     }
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[little_s0_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[little_s0_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[little_s1_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[little_s1_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[a_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[a_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[b_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[b_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[c_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[c_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[d_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[d_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[e_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[e_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[f_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[f_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[g_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[g_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[h_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[h_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[big_s0_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[big_s0_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[big_s1_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[big_s1_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[not_e_and_g_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[not_e_and_g_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[e_and_f_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[e_and_f_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[ch_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[ch_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[a_and_b_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[a_and_b_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[a_and_c_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[a_and_c_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[b_and_c_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[b_and_c_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[maj_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[maj_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[xor_tmp_0_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[xor_tmp_0_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[xor_tmp_1_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[xor_tmp_1_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[xor_tmp_2_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[xor_tmp_2_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[xor_tmp_3_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[xor_tmp_3_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }

        // for bit in 0..32 {
        //     let mut c = builder.sub_extension(one_ext, curr_row[xor_tmp_4_bit(bit)]);
        //     c = builder.mul_extension(c, curr_row[xor_tmp_4_bit(bit)]);
        //     yield_constr.constraint(builder, c);
        // }
    }

    fn constraint_degree(&self) -> usize {
        3
    }
}

fn eval_bits_are_bits<F, P>(curr_row: &[P; NUM_COLS], yield_constr: &mut ConstraintConsumer<P>)
where
    F: Field,
    P: PackedField<Scalar = F>,
{
    // step bits
    for bit in 0..NUM_STEPS_PER_HASH {
        yield_constr.constraint((P::ONES - curr_row[step_bit(bit)]) * curr_row[step_bit(bit)]);
    }

    for bit in 0..32 {
        // wi
        yield_constr.constraint((P::ONES - curr_row[wi_bit(bit)]) * curr_row[wi_bit(bit)]);

        // wi minus 2
        yield_constr
            .constraint((P::ONES - curr_row[wi_minus_2_bit(bit)]) * curr_row[wi_minus_2_bit(bit)]);

        // wi minus 15
        yield_constr.constraint(
            (P::ONES - curr_row[wi_minus_15_bit(bit)]) * curr_row[wi_minus_15_bit(bit)],
        );
    }

    // s0
    for bit in 0..32 {
        yield_constr
            .constraint((P::ONES - curr_row[little_s0_bit(bit)]) * curr_row[little_s0_bit(bit)]);
    }

    // s1
    for bit in 0..32 {
        yield_constr
            .constraint((P::ONES - curr_row[little_s1_bit(bit)]) * curr_row[little_s1_bit(bit)]);
    }

    // a
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[a_bit(bit)]) * curr_row[a_bit(bit)]);
    }

    // b
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[b_bit(bit)]) * curr_row[b_bit(bit)]);
    }

    // c
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[c_bit(bit)]) * curr_row[c_bit(bit)]);
    }

    // e
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[e_bit(bit)]) * curr_row[e_bit(bit)]);
    }

    // f
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[f_bit(bit)]) * curr_row[f_bit(bit)]);
    }

    // g
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[g_bit(bit)]) * curr_row[g_bit(bit)]);
    }

    // S0
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[big_s0_bit(bit)]) * curr_row[big_s0_bit(bit)]);
    }

    // S1
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[big_s1_bit(bit)]) * curr_row[big_s1_bit(bit)]);
    }

    // (not e) and g
    for bit in 0..32 {
        yield_constr.constraint(
            (P::ONES - curr_row[not_e_and_g_bit(bit)]) * curr_row[not_e_and_g_bit(bit)],
        );
    }

    // e and f
    for bit in 0..32 {
        yield_constr
            .constraint((P::ONES - curr_row[e_and_f_bit(bit)]) * curr_row[e_and_f_bit(bit)]);
    }

    // ch
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[ch_bit(bit)]) * curr_row[ch_bit(bit)]);
    }

    // a and b
    for bit in 0..32 {
        yield_constr
            .constraint((P::ONES - curr_row[a_and_b_bit(bit)]) * curr_row[a_and_b_bit(bit)]);
    }

    // a and c
    for bit in 0..32 {
        yield_constr
            .constraint((P::ONES - curr_row[a_and_c_bit(bit)]) * curr_row[a_and_c_bit(bit)]);
    }

    // b and c
    for bit in 0..32 {
        yield_constr
            .constraint((P::ONES - curr_row[b_and_c_bit(bit)]) * curr_row[b_and_c_bit(bit)]);
    }

    // maj
    for bit in 0..32 {
        yield_constr.constraint((P::ONES - curr_row[maj_bit(bit)]) * curr_row[maj_bit(bit)]);
    }

    // tmps
    for bit in 0..29 {
        yield_constr
            .constraint((P::ONES - curr_row[xor_tmp_0_bit(bit)]) * curr_row[xor_tmp_0_bit(bit)]);
    }
    for bit in 0..22 {
        yield_constr
            .constraint((P::ONES - curr_row[xor_tmp_1_bit(bit)]) * curr_row[xor_tmp_1_bit(bit)]);
    }
    for bit in 0..32 {
        yield_constr
            .constraint((P::ONES - curr_row[xor_tmp_2_bit(bit)]) * curr_row[xor_tmp_2_bit(bit)]);

        yield_constr
            .constraint((P::ONES - curr_row[xor_tmp_3_bit(bit)]) * curr_row[xor_tmp_3_bit(bit)]);

        yield_constr
            .constraint((P::ONES - curr_row[xor_tmp_4_bit(bit)]) * curr_row[xor_tmp_4_bit(bit)]);
    }
}

#[derive(Debug, Clone, Default)]
pub struct Sha2StarkCompressor {
    inputs: Vec<([u32; 8], [u32; 8])>,
}

impl Sha2StarkCompressor {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn add_instance(&mut self, left: [u32; 8], right: [u32; 8]) {
        self.inputs.push((left, right));
    }

    /// returns the generated trace against which a proof may be generated
    pub fn generate<F: Field>(self) -> Vec<PolynomialValues<F>> {
        let max_rows = 1 << log2_ceil(self.inputs.len() * NUM_STEPS_PER_HASH);
        let mut generator = Sha2TraceGenerator::<F>::new(max_rows);
        for (left, right) in self.inputs.into_iter() {
            generator.gen_hash(left, right);
        }

        generator.into_polynomial_values()
    }
}

#[cfg(test)]
mod tests {
    use anyhow::Result;
    use plonky2::hash::hash_types::BytesHash;
    use plonky2::plonk::config::{GenericConfig, PoseidonGoldilocksConfig};
    use plonky2::util::timing::TimingTree;

    // use plonky2::util::timing::TimingTree;
    use super::*;
    use crate::prover::prove;
    use crate::verifier::verify_stark_proof;
    // use crate::config::StarkConfig;
    // use crate::prover::prove;
    // use crate::sha256_stark::generation::Sha2TraceGenerator;
    use crate::{config::StarkConfig, stark_testing::test_stark_low_degree};
    // use crate::verifier::verify_stark_proof;

    #[test]
    fn test_stark_degree() -> Result<()> {
        const D: usize = 2;
        type C = PoseidonGoldilocksConfig;
        type F = <C as GenericConfig<D>>::F;
        type S = Sha2CompressionStark<F, D>;

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
    fn test_single() -> Result<()> {
        const D: usize = 2;
        type C = PoseidonGoldilocksConfig;
        type F = <C as GenericConfig<D>>::F;
        type S = Sha2CompressionStark<F, D>;

        let mut left_input = [0u32; 8];
        let mut right_input = [0u32; 8];
        for i in 0..8 {
            left_input[i] = i as u32;
            right_input[i] = i as u32 + 8;
        }

        let mut generator = Sha2TraceGenerator::<F>::new(128);
        generator.gen_hash(left_input, right_input);

        let config = StarkConfig::standard_fast_config();
        let stark = S::new();
        let trace = generator.into_polynomial_values();
        let mut timing = TimingTree::default();
        let proof = prove::<F, C, S, D>(stark, &config, trace, [], &mut timing)?;

        verify_stark_proof(stark, proof, &config)?;

        Ok(())
    }

    #[test]
    fn test_multiple() -> Result<()> {
        const D: usize = 2;
        type C = PoseidonGoldilocksConfig;
        type F = <C as GenericConfig<D>>::F;
        type S = Sha2CompressionStark<F, D>;

        let mut compressor = Sha2StarkCompressor::new();

        let left = to_u32_array_be::<8>(BytesHash::<32>::rand().0);
        let right = to_u32_array_be::<8>(BytesHash::<32>::rand().0);
        compressor.add_instance(left, right);

        let left = to_u32_array_be::<8>(BytesHash::<32>::rand().0);
        let right = to_u32_array_be::<8>(BytesHash::<32>::rand().0);
        compressor.add_instance(left, right);

        let trace = compressor.generate();

        let config = StarkConfig::standard_fast_config();
        let stark = S::new();
        let mut timing = TimingTree::default();
        let proof = prove::<F, C, S, D>(stark, &config, trace, [], &mut timing)?;

        verify_stark_proof(stark, proof, &config)?;

        Ok(())
    }
}
