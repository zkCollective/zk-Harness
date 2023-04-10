pub const TREE_DEPTH: usize = 5;
pub const TREE_WIDTH: usize = 1 << (TREE_DEPTH - 1);
pub fn level_width(level: usize) -> usize {
	1 << (TREE_DEPTH - 1 - level)
}

// see `mod.rs` for an explanation of what this is for
pub fn get_level_end_shift(i: usize, level: usize) -> usize {
	TREE_WIDTH - level_width(level) + 3 + i
}

pub const WORDS_PER_HASH: usize = 8;

// 17 hashes - 16 for the leaves, one for the root
pub const NUM_PUBLIC_INPUTS: usize = (TREE_WIDTH + 1) * WORDS_PER_HASH;
pub const NUM_COLS: usize = LAST_COL + 1;

pub const LEAVES_PI_WORDS_START: usize = 0;
pub fn pi_leaf_i_word(i: usize, word: usize) -> usize {
	LEAVES_PI_WORDS_START + i * WORDS_PER_HASH + word
}

pub const ROOT_PI_WORDS_START: usize = TREE_WIDTH * WORDS_PER_HASH;
pub fn pi_root_word(word: usize) -> usize {
	ROOT_PI_WORDS_START + word
}

pub const INPUT_FILTER: usize = 0;
pub const OUTPUT_FILTER: usize = INPUT_FILTER + 1;
pub const PC: usize = OUTPUT_FILTER + 1;
pub const LEVEL_DONE_FLAG: usize = PC + 1;
pub const HALF_LEVEL_WIDTH: usize = LEVEL_DONE_FLAG + 1;
pub const HASH_IDX: usize = HALF_LEVEL_WIDTH + 1;

pub const LEVEL_FLAGS_START: usize = HASH_IDX + 1;
pub fn level_flag(level: usize) -> usize {
	LEVEL_FLAGS_START + level
}

pub const VALS_START: usize = LEVEL_FLAGS_START + TREE_DEPTH - 1;
pub fn val_i_word(i: usize, word: usize) -> usize {
	VALS_START + i * 8 + word
}

pub const HASH_INPUT_0_START: usize = VALS_START + TREE_WIDTH * WORDS_PER_HASH;
pub fn hash_input_0_word(word: usize) -> usize {
	HASH_INPUT_0_START + word
}

pub const HASH_INPUT_1_START: usize = HASH_INPUT_0_START + WORDS_PER_HASH;
pub fn hash_input_1_word(word: usize) -> usize {
	HASH_INPUT_1_START + word
}

pub const HASH_OUTPUT_START: usize = HASH_INPUT_1_START + WORDS_PER_HASH;
pub fn hash_output_word(word: usize) -> usize {
	HASH_OUTPUT_START + word
}

pub const LAST_COL: usize = HASH_OUTPUT_START + WORDS_PER_HASH - 1;
