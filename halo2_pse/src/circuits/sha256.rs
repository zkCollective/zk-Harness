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
    pub sha_data: usize,
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
        match self.sha_data {
            0..=64 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 2]
            ).unwrap(),
            65..=128 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 3]
            ).unwrap(),
            129..=256 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 5]
            ).unwrap(),
            257..=512 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 9]
            ).unwrap(),
            513..=1024 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 17]
            ).unwrap(),
            1025..=2048 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 33]
            ).unwrap(),
            2049..=4096 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 65]
            ).unwrap(),
            4097..=8192 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 129]
            ).unwrap(),
            8193..=16384 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 257]
            ).unwrap(),
            16385..=32768 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 513]
            ).unwrap(),
            32769..=65536 => Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
                &[BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 1025]
            ).unwrap(),
            _ => panic!("unexpected sha data"),
        };
        
        // Should be as following:
        // 
        //Sha256::digest(table16_chip, layouter.namespace(|| "'sha one'"),
        //    self.sha_data.clone().into_boxed_slice().as_ref().try_into().unwrap()).unwrap();
        Ok(())
    }
}

fn hex_to_binary_words(input: &str) -> Vec<BlockWord> {
    let bytes = hex::decode(input).expect("Invalid hex input");
    let binary_string = bytes.iter().map(|byte| format!("{:08b}", byte)).collect::<String>();

    let padded_binary_string = format!("{:0>width$}", binary_string, width = (binary_string.len() + 31) / 32 * 32);

    let words: Vec<BlockWord> = padded_binary_string
        .chars()
        .collect::<Vec<char>>()
        .chunks(32)
        .map(|chunk| {
            let word_str: String = chunk.iter().collect();
            let value = u32::from_str_radix(&word_str, 2).unwrap();
            BlockWord(Value::known(value))
        })
        .collect();
    println!("words length {:?}", words.len());
    words
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Sha256Data {
    PreImage: String,
    Hash: String,
}

pub fn get_sha256_data (
    input_file_str: String
) -> (u32, usize ){
    // FIXME this function should return Vec<BlockWord>
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
    //let sha_data = hex_to_binary_words(&preimage);
    return (k, preimage_length);
}

#[cfg(test)]
mod tests {
    use halo2_proofs::circuit::Value;

    use super::*;

    #[test]
    fn test_exponentiation_circuit() {
        use halo2_proofs::dev::MockProver;
        let k = 17;
        let sha_data = vec![BlockWord(Value::known(0b01111000100000000000000000000000)); 16 * 1];
        let circuit = Sha256Circuit {
            sha_data: sha_data,
        };
        let public_input = vec![];
        let prover = MockProver::run(k, &circuit, public_input.clone()).unwrap();
        prover.assert_satisfied();
    }
}
