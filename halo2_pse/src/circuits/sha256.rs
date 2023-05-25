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
        Sha256::digest(table16_chip.clone(), layouter.namespace(|| "'sha one'"),
            self.sha_data.clone().into_boxed_slice().as_ref().try_into().unwrap());
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
) -> (u32, Vec<BlockWord> ){
    let data: Sha256Data = serde_json::from_str(&input_file_str).expect("Cannot read json string");
    let preimage_length: usize = data.PreImage.len();
    println!("Preimage length {}", preimage_length);
    let k = match preimage_length {
        0..=64 => 17,
        65..=128 => 17,
        129..=256 => 17,
        257..=512 => 17,
        513..=1024 => 17,
        1025..=2048 => 17,
        2049..=4096 => 17,
        4097..=8192 => 17,
        8193..=16384 => 17,
        16385..=32768 => 17,
        32769..=65536 => 17,
        _ => 1,
    };
    let preimage = data.PreImage;
    let _hash = data.Hash;
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
