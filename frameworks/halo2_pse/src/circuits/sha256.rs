use halo2_gadgets::sha256::{BlockWord, Sha256, Table16Chip, Table16Config};
use halo2_proofs::{
    circuit::{Layouter, SimpleFloorPlanner, Value},
    plonk::{
        Circuit, 
        ConstraintSystem, Error, 
    }, 
    halo2curves::bn256::Fr
};
use std::convert::TryInto;
use serde::{Deserialize, Serialize};
use serde_json;

#[derive(Default, Clone)]
pub struct Sha256Circuit {
    // FIXME that should be Vec<BlockWord>
    pub sha_data: Vec<BlockWord>,
}

impl Circuit<Fr> for Sha256Circuit {
    type Config = Table16Config;
    type FloorPlanner = SimpleFloorPlanner;

    fn without_witnesses(&self) -> Self {
        Self::default()
    }

    fn configure(meta: &mut ConstraintSystem<Fr>) -> Self::Config {
        Table16Chip::configure(meta)
    }

    fn synthesize(
        &self,
        config: Self::Config,
        mut layouter: impl Layouter<Fr>,
    ) -> Result<(), Error> {
        Table16Chip::load(config.clone(), &mut layouter)?;
        let table16_chip = Table16Chip::construct(config);

        // FIXME Remove that code
        //match self.sha_data {
        //    0..=64 => 
        //        Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16]
        //    ).unwrap(),
        //    65..=128 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 3]
        //    ).unwrap(),
        //    129..=256 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 5]
        //    ).unwrap(),
        //    257..=512 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 9]
        //    ).unwrap(),
        //    513..=1024 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 17]
        //    ).unwrap(),
        //    1025..=2048 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 33]
        //    ).unwrap(),
        //    2049..=4096 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 65]
        //    ).unwrap(),
        //    4097..=8192 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 129]
        //    ).unwrap(),
        //    8193..=16384 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 257]
        //    ).unwrap(),
        //    16385..=32768 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 513]
        //    ).unwrap(),
        //    32769..=65536 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //        &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 1025]
        //    ).unwrap(),
        //    _ => panic!("unexpected sha data"),
        //};
        
        // Should be as following:
        // 
        Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
            self.sha_data.clone().into_boxed_slice().as_ref().try_into().unwrap()).unwrap();
        Ok(())
    }
}

fn hex_to_binary_words(input: &str) -> Vec<BlockWord> {
    let mut block_words: Vec<BlockWord> = Vec::new();

    // Convert the input string to bytes
    let bytes_data = hex::decode(input).unwrap();

    // Split the bytes into chunks of 4 bytes (32 bits)
    let chunk_size = 4;
    for chunk in bytes_data.chunks_exact(chunk_size) {
        // Create the u32 value from the chunk
        let value = u32::from_be_bytes(chunk.try_into().unwrap());

        // Create the BlockWord and push it to the vector
        let block_word = BlockWord(Value::known(value));
        block_words.push(block_word);
    }
    if block_words.len() < 16 {
        println!("Input is less than 16 words, padding with zeros is required");
        let padding = 16 - block_words.len();
        for _ in 0..padding {
            block_words.push(BlockWord(Value::known(0)));
        }
    }
    block_words
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Sha256Data {
    PreImage: String,
    Hash: String,
}

pub fn get_sha256_data (
    input_file_str: String
) -> (u32, Vec<BlockWord>){
    let data: Sha256Data = serde_json::from_str(&input_file_str).expect("Cannot read json string");
    let preimage_length: usize = data.PreImage.len();
    println!("Preimage length {}", preimage_length);
    let k = match preimage_length {
        0..=2048 => 17,
        2049..=4096 => 18,
        4097..=8192 => 19,
        8193..=16384 => 20,
        16385..=32768 => 21,
        32769..=65536 => 22,
        _ => 1,
    };
    let preimage = data.PreImage;
    let _hash = data.Hash;
    // FIXME fix that function
    let sha_data = hex_to_binary_words(&preimage);
    return (k, sha_data);
}

#[cfg(test)]
mod tests {
    use halo2_proofs::circuit::Value;

    use super::*;

    #[test]
    fn test_exponentiation_circuit() {
        use halo2_proofs::dev::MockProver;
        let (k, sha_data) = get_sha256_data(String::from("{ \"PreImage\": \"28ca152c94e1db7f7d892d27b5b674dd414028635c3c0321289f9afb0eee906a\", \"Hash\": \"c80263908bbc7bece8d340575547cf920206178f36f37ea787abe81ef2e24043\" }"));
        let circuit = Sha256Circuit {
            sha_data: sha_data,
        };
        let public_input = vec![];
        let prover = MockProver::run(k, &circuit, public_input.clone()).unwrap();
        prover.assert_satisfied();
    }
}