pub const NUM_PIS: usize = 0;
pub const NUM_COLS: usize = LAST_COL + 1;
pub const NUM_STEPS_PER_HASH: usize = 65;

pub const HASH_IDX: usize = 0;
pub const STEP_BITS_START: usize = HASH_IDX + 1;
pub fn step_bit(i: usize) -> usize {
    STEP_BITS_START + i
}

pub const INPUT_START: usize = STEP_BITS_START + NUM_STEPS_PER_HASH;
pub fn input_i(i: usize) -> usize {
    INPUT_START + i
}
pub const INPUT_FILTER: usize = INPUT_START + 16;

pub const WI_BITS_START: usize = INPUT_FILTER + 1;
pub fn wi_bit(bit: usize) -> usize {
    WI_BITS_START + bit
}

pub const WI_MINUS_2_START: usize = WI_BITS_START + 32;
pub fn wi_minus_2_bit(bit: usize) -> usize {
    WI_MINUS_2_START + bit
}

pub const WI_MINUS_15_START: usize = WI_MINUS_2_START + 32;
pub fn wi_minus_15_bit(bit: usize) -> usize {
    WI_MINUS_15_START + bit
}

pub const NUM_WIS_FIELD: usize = 13;
pub const WIS_FIELD_START: usize = WI_MINUS_15_START + 32;
pub fn wi_field(i: usize) -> usize {
    match i {
        15 | 13 | 0 => panic!("invalid index into field-encoded wis"),
        1..=12 => WIS_FIELD_START + i - 1, 
        14 => WIS_FIELD_START + i - 2,
        _ => unreachable!()
    }
}

pub const XOR_TMP_0_START: usize = WIS_FIELD_START + NUM_WIS_FIELD;
pub fn xor_tmp_0_bit(bit: usize) -> usize {
    XOR_TMP_0_START + bit
}

pub const XOR_TMP_1_START: usize = XOR_TMP_0_START + 29;
pub fn xor_tmp_1_bit(bit: usize) -> usize {
    XOR_TMP_1_START + bit
}

pub const XOR_TMP_2_START: usize = XOR_TMP_1_START + 22;
pub fn xor_tmp_2_bit(bit: usize) -> usize {
    XOR_TMP_2_START + bit
}

pub const XOR_TMP_3_START: usize = XOR_TMP_2_START + 32;
pub fn xor_tmp_3_bit(bit: usize) -> usize {
    XOR_TMP_3_START + bit
}

pub const XOR_TMP_4_START: usize = XOR_TMP_3_START + 32;
pub fn xor_tmp_4_bit(bit: usize) -> usize {
    XOR_TMP_4_START + bit
}

pub const LITTLE_S0_START: usize = XOR_TMP_4_START + 32;
pub fn little_s0_bit(bit: usize) -> usize {
    LITTLE_S0_START + bit
}

pub const LITTLE_S1_START: usize = LITTLE_S0_START + 32;
pub fn little_s1_bit(bit: usize) -> usize {
    LITTLE_S1_START + bit
}

pub const KI: usize = LITTLE_S1_START + 32;
pub const WI_FIELD: usize = KI + 1;
pub const WI_QUOTIENT: usize = WI_FIELD + 1;

pub const A_START: usize = WI_QUOTIENT + 1;
pub fn a_bit(bit: usize) -> usize {
    A_START + bit
}

pub const B_START: usize = A_START + 32;
pub fn b_bit(bit: usize) -> usize {
    B_START + bit
}

pub const C_START: usize = B_START + 32;
pub fn c_bit(bit: usize) -> usize {
    C_START + bit
}

pub const D_COL: usize = C_START + 32;

pub const E_START: usize = D_COL + 1;
pub fn e_bit(bit: usize) -> usize {
    E_START + bit
}

pub const F_START: usize = E_START + 32;
pub fn f_bit(bit: usize) -> usize {
    F_START + bit
}

pub const G_START: usize = F_START + 32;
pub fn g_bit(bit: usize) -> usize {
    G_START + bit
}

pub const H_COL: usize = G_START + 32;

pub const BIG_S0_START: usize = H_COL + 1;
pub fn big_s0_bit(bit: usize) -> usize {
    BIG_S0_START + bit
}

pub const BIG_S1_START: usize = BIG_S0_START + 32;
pub fn big_s1_bit(bit: usize) -> usize {
    BIG_S1_START + bit
}

pub const NOT_E_AND_G_START: usize = BIG_S1_START + 32;
pub fn not_e_and_g_bit(bit: usize) -> usize {
    NOT_E_AND_G_START + bit
}

pub const E_AND_F_START: usize = NOT_E_AND_G_START + 32;
pub fn e_and_f_bit(bit: usize) -> usize {
    E_AND_F_START + bit
}

pub const CH_START: usize = E_AND_F_START + 32;
pub fn ch_bit(bit: usize) -> usize {
    CH_START + bit
}

pub const A_AND_B: usize = CH_START + 32;
pub fn a_and_b_bit(bit: usize) -> usize {
    A_AND_B + bit
}

pub const A_AND_C: usize = A_AND_B + 32;
pub fn a_and_c_bit(bit: usize) -> usize {
    A_AND_C + bit
}

pub const B_AND_C: usize = A_AND_C + 32;
pub fn b_and_c_bit(bit: usize) -> usize {
    B_AND_C + bit
}

pub const MAJ_START: usize = B_AND_C + 32;
pub fn maj_bit(bit: usize) -> usize {
    MAJ_START + bit
}

pub const BIG_SO_FIELD: usize = MAJ_START + 32;
pub const BIG_S1_FIELD: usize = BIG_SO_FIELD + 1;
pub const CH_FIELD: usize = BIG_S1_FIELD + 1;
pub const MAJ_FIELD: usize = CH_FIELD + 1;

pub const A_NEXT_FIELD: usize = MAJ_FIELD + 1;
pub const E_NEXT_FIELD: usize = A_NEXT_FIELD + 1;

pub const A_NEXT_QUOTIENT: usize = E_NEXT_FIELD + 1;
pub const E_NEXT_QUOTIENT: usize = A_NEXT_QUOTIENT + 1;

pub const HIS_START: usize = E_NEXT_QUOTIENT + 1;
pub fn h_i(i: usize) -> usize {
    HIS_START + i
}

pub const HIS_NEXT_FIELD_START: usize = HIS_START + 8;
pub fn h_i_next_field(i: usize) -> usize {
    HIS_NEXT_FIELD_START + i
}

pub const HIS_NEXT_QUOTIENT_START: usize = HIS_NEXT_FIELD_START + 8;
pub fn h_i_next_quotient(i: usize) -> usize {
    HIS_NEXT_QUOTIENT_START + i
}

pub const OUTPUT_COLS_START: usize = HIS_NEXT_QUOTIENT_START + 8;
pub fn output_i(i: usize) -> usize {
    OUTPUT_COLS_START + i
}

pub const OUTPUT_FILTER: usize = OUTPUT_COLS_START + 8;

pub const LAST_COL: usize = OUTPUT_FILTER;
