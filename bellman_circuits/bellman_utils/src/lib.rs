use std::{fs::File};
use std::fs;
use bellman::groth16::Proof;
use bls12_381::Bls12;

pub fn measure_size_in_bytes(proof: &Proof<Bls12>) -> usize {
    // TODO: Should we serialize the proof in another format?
    // Serialize data and save it to a temporary file
    let temp_file_path = "temp_file.bin";

    let mut file = File::create(temp_file_path).expect("Could not create temp file");
    proof.write(&mut file).expect("Could not write proof to file");

    // Measure the size of the file
    let file_size: usize = fs::metadata(&temp_file_path).expect("Cannot read the size of temp file").len() as usize;

    // Convert file size to MB
    let size_in_mb = file_size;

    // Remove the temporary file
    fs::remove_file(&temp_file_path).expect("Cannot remove temp file");

    return size_in_mb;
}