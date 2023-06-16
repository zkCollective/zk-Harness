use rand::rngs::OsRng;
use halo2_proofs::{
    plonk::{keygen_vk, keygen_pk, create_proof, verify_proof, ProvingKey, VerifyingKey, Circuit},
    poly::{
        kzg::{
            commitment::{ParamsKZG, KZGCommitmentScheme}, 
            multiopen::{ProverGWC, VerifierGWC}, 
            strategy::AccumulatorStrategy
        }
    }, 
    halo2curves::bn256::{Bn256, G1Affine, Fr}, 
    transcript::{
        Blake2bRead, Blake2bWrite, Challenge255, TranscriptReadBuffer, TranscriptWriterBuffer,
    },
};
use halo2_proofs::SerdeFormat;
use criterion::{
    measurement::Measurement, BenchmarkGroup,
};
use serde::Serialize;
use std::{env, fs::File};
use std::io::{BufReader, Read, Write};
use std::process;
use psutil;
use std::fs;
use clap::Parser;

pub const DEFAULT_SERDE_FORMAT: SerdeFormat = SerdeFormat::RawBytesUnchecked;

#[derive(Serialize)]
struct Results {
    initial_rss: u64,
    setup_rss: u64,
    proof_rss: u64,
    verify_rss: u64,
    proof_size: usize,
}


#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
pub struct BinaryArgs {
    #[arg(short, long)]
    pub input: String,

    #[arg(short, long)]
    pub phase: String,

    #[arg(short, long)]
    pub params: Option<String>,

    #[arg(short, long)]
    pub pk: Option<String>,

    #[arg(short, long)]
    pub vk: Option<String>,

    #[arg(short, long)]
    pub proof: Option<String>,
}


pub fn setup_circuit<C :Circuit<halo2_proofs::halo2curves::bn256::Fr> + Clone>(
    k: u32,
    circuit: C
) -> (ParamsKZG<Bn256>, VerifyingKey<G1Affine>, ProvingKey<G1Affine>) {
    let params = ParamsKZG::<Bn256>::setup(k, OsRng);
    let vk = keygen_vk(&params, &circuit).unwrap();
    let pk = keygen_pk(&params, vk.clone(), &circuit).unwrap();
    return (params, vk, pk);
}

pub fn prove_circuit<C :Circuit<halo2_proofs::halo2curves::bn256::Fr> + Clone>(
    params: &ParamsKZG<Bn256>,
    pk: &ProvingKey<G1Affine>,
    circuit: C,
    public_input: &[&[halo2_proofs::halo2curves::bn256::Fr]]
) -> Vec<u8> {
    let mut transcript = Blake2bWrite::<_, _, Challenge255<_>>::init(vec![]);
    create_proof::<
        KZGCommitmentScheme<Bn256>,
        ProverGWC<'_, Bn256>,
        Challenge255<G1Affine>,
        _,
        Blake2bWrite<Vec<u8>, G1Affine, Challenge255<_>>,
        _,
    >(
        &params,
        &pk,
        &[circuit],
        &[public_input],
        OsRng,
        &mut transcript,
    )
    .expect("prover should not fail");
    return transcript.finalize();
}

pub fn verify_circuit<C: Circuit<halo2_proofs::halo2curves::bn256::Fr> + Clone>(
    params: &ParamsKZG<Bn256>,
    vk: &VerifyingKey<G1Affine>,
    proof: &[u8],
    public_input: &[&[halo2_proofs::halo2curves::bn256::Fr]]
) {
    let strategy = AccumulatorStrategy::new(&params);
    let mut transcript = Blake2bRead::<_, _, Challenge255<_>>::init(&proof[..]);
    verify_proof::<KZGCommitmentScheme<_>, VerifierGWC<_>, _, _, _>(
        &params,
        &vk,
        strategy,
        &[public_input],
        &mut transcript,
    ).unwrap();
}

pub fn bench_circuit<M: Measurement, C: Circuit<halo2_proofs::halo2curves::bn256::Fr> + Clone>(
    c: &mut BenchmarkGroup<'_, M>,
    k: u32,
    circuit: C,
    public_input: &[&[halo2_proofs::halo2curves::bn256::Fr]]
) {

    let mut params: Option<ParamsKZG<Bn256>> = None;
    let mut pk: Option<ProvingKey<G1Affine>> = None;
    let mut vk: Option<VerifyingKey<G1Affine>> = None;


    c.bench_function("setup", |b| {
        b.iter(|| { 
            params = Some(ParamsKZG::<Bn256>::setup(k, OsRng));
            vk = Some(keygen_vk(params.as_ref().unwrap(), &circuit.clone()).expect("keygen_vk should not fail"));
            pk = Some(keygen_pk(params.as_ref().unwrap(), vk.clone().unwrap(), &circuit).expect("keygen_pk should not fail"));
        });
    });

    let mut proof: Option<Vec<u8>> = None;

    c.bench_function("prove", |b| {
         b.iter(|| { 
            let mut transcript = Blake2bWrite::<_, _, Challenge255<_>>::init(vec![]);
            create_proof::<
                KZGCommitmentScheme<Bn256>,
                ProverGWC<'_, Bn256>,
                Challenge255<G1Affine>,
                _,
                Blake2bWrite<Vec<u8>, G1Affine, Challenge255<_>>,
                _,
            >(
                params.as_ref().unwrap(),
                pk.as_ref().unwrap(),
                &[circuit.clone()],
                &[public_input],
                OsRng,
                &mut transcript,
            )
            .expect("prover should not fail");
            proof = Some(transcript.finalize());
         });
    });

    c.bench_function("verify", |b| {
        let proof = proof.clone().unwrap();
         b.iter(|| { 
            //println!("proof len: {}", proof.len());
            let strategy = AccumulatorStrategy::new(params.as_ref().unwrap(),);
            let mut transcript = Blake2bRead::<_, _, Challenge255<_>>::init(&proof[..]);
            let _strategy = verify_proof::<KZGCommitmentScheme<_>, VerifierGWC<_>, _, _, _>(
                params.as_ref().unwrap(),
                pk.as_ref().unwrap().get_vk(),
                strategy,
                &[public_input],
                &mut transcript,
            ).unwrap();
        })
    });
}

pub fn read_file_contents(file_name: String) -> String {
    let mut file = File::open(file_name).expect("Cannot load file");
    let mut file_str = String::new();
    file.read_to_string(&mut file_str).expect("Cannot read file");
    return file_str;
}

pub fn read_file_from_env_var(env_var_name: String) -> String {
    let input_file = env::var(env_var_name.clone()).unwrap_or_else(|_| {
        println!("Please set the {} environment variable to point to the input file", env_var_name);
        process::exit(1);
    });
    return read_file_contents(input_file);
}

/// Return current RSS memory in bytes
pub fn get_memory() -> u64 {
    let current_process = psutil::process::Process::current().expect("Cannot get current process");
    let mem = current_process.memory_info().expect("Cannot get memory info");
    return mem.rss();
}

pub fn save_results(
    initial_rss: u64,
    setup_rss: u64,
    proof_rss: u64,
    verify_rss: u64,
    proof_size: usize,
    output_file_str: String,
) {
    // Create a Results struct with the provided data
    let results = Results {
        initial_rss,
        setup_rss,
        proof_rss,
        verify_rss,
        proof_size,
    };

    // Serialize the Results struct to JSON
    let serialized = serde_json::to_string(&results).expect("Could not serialize results");

    // Create a file and write the JSON data
    std::fs::write(
        output_file_str,
        serialized
    ).expect("Could not write to file");
}

pub fn measure_size_in_bytes(data: &[u8]) -> usize {
    // TODO: Should we serialize the proof in another format?
    // Serialize data and save it to a temporary file
    let temp_file_path = "temp_file.bin";
    std::fs::write(
        temp_file_path,
        data
    ).expect("Could not write to temp file");

    // Measure the size of the file
    let file_size: usize = fs::metadata(&temp_file_path).expect("Cannot read the size of temp file").len() as usize;

    // Convert file size to MB
    let size_in_mb = file_size;

    // Remove the temporary file
    fs::remove_file(&temp_file_path).expect("Cannot remove temp file");

    return size_in_mb;
}

pub fn save_params(params_file: String, params: &ParamsKZG<Bn256>) {
    let mut params_buf = Vec::new();
    params.write_custom(&mut params_buf, DEFAULT_SERDE_FORMAT).unwrap();

    let mut params_file = File::create(params_file).unwrap();
    params_file.write_all(&params_buf[..]).unwrap();
}

pub fn load_params(params_file: String) -> ParamsKZG<Bn256> {
    let f = File::open(params_file).unwrap();
    return ParamsKZG::<Bn256>::read_custom::<_>(&mut BufReader::new(f), DEFAULT_SERDE_FORMAT).unwrap();
}

pub fn save_vk(vk_file: String, vk: &VerifyingKey<G1Affine>) {
    let mut vk_file = File::create(vk_file).unwrap();
    vk.write(&mut vk_file, DEFAULT_SERDE_FORMAT).unwrap();
}

pub fn load_vk<C: Circuit<Fr>>(vk_file: String) -> VerifyingKey<G1Affine> {
    let f = File::open(vk_file).unwrap();
    return VerifyingKey::<G1Affine>::read::<_, C>(&mut BufReader::new(f), DEFAULT_SERDE_FORMAT).unwrap();
}

pub fn save_pk(pk_file: String, pk: &ProvingKey<G1Affine>) {
    let mut pk_file = File::create(pk_file).unwrap();
    pk.write(&mut pk_file, DEFAULT_SERDE_FORMAT).unwrap();
}

pub fn load_pk<C: Circuit<Fr>>(pk_file: String) -> ProvingKey<G1Affine> {
    let f = File::open(pk_file).unwrap();
    return ProvingKey::<G1Affine>::read::<_, C>(&mut BufReader::new(f), DEFAULT_SERDE_FORMAT).unwrap();
}

pub fn save_proof(proof_file: String, proof: &[u8]) {
    let mut proof_file = File::create(proof_file).unwrap();
    proof_file.write_all(proof).unwrap();
}

pub fn load_proof(proof_file: String) -> Vec<u8> {
    let f = File::open(proof_file).unwrap();
    return BufReader::new(f).bytes().map(|x| x.unwrap()).collect();
}

pub fn f_setup<C: Circuit<Fr> + Clone>(k: u32, circuit: C, params_file: String, vk_file: String, pk_file: String) {
    let (params, vk, pk) = setup_circuit(k, circuit.clone());
    save_params(params_file, &params);
    save_vk(vk_file, &vk);
    save_pk(pk_file, &pk);
}

pub fn f_prove<C: Circuit<Fr> + Clone>(circuit: C, params_file: String, pk_file: String, public_input: &[&[Fr]], proof_file: String) {
    let params = load_params(params_file);
    let pk = load_pk::<C>(pk_file);
    let proof = prove_circuit(&params, &pk, circuit, public_input);
    save_proof(proof_file, &proof);
}

pub fn f_verify<C: Circuit<Fr> + Clone>(params_file: String, vk_file: String, proof_file: String, public_input: &[&[Fr]]) {
    let params = load_params(params_file);
    let vk = load_vk::<C>(vk_file);
    let proof = load_proof(proof_file);
    verify_circuit::<C>(&params, &vk, &proof, public_input);
}