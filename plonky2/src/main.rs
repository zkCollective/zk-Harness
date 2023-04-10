#![feature(generic_const_exprs)]
use std::env;
use anyhow::{Result, Ok};
use log::{Level, LevelFilter};
use plonky2::hash::hash_types::RichField;
use plonky2::iop::witness::{PartialWitness, WitnessWrite};
use plonky2::plonk::circuit_builder::CircuitBuilder;
use plonky2::plonk::config::{GenericConfig, PoseidonGoldilocksConfig, AlgebraicHasher, Hasher};
use plonky2::plonk::prover::prove;
use plonky2::timed;
use plonky2::util::timing::TimingTree;
use plonky2_field::extension::Extendable;
use plonky2_sha256::circuit::{array_to_bits, make_circuits};
use sha2::{Digest, Sha256};
use plonky2::plonk::circuit_data::{
    CircuitConfig, CommonCircuitData, VerifierCircuitTarget, VerifierOnlyCircuitData,
};
use plonky2::plonk::proof::ProofWithPublicInputs;

pub fn prove_sha256(msg: &[u8]) -> Result<()> {
    let mut hasher = Sha256::new();
    hasher.update(msg);
    let hash = hasher.finalize();
    println!("Hash: {:#04X}", hash);

    let msg_bits = array_to_bits(msg);
    let len = msg.len() * 8;
    println!("block count: {}", (len + 65 + 511) / 512);
    const D: usize = 2;
    type C = PoseidonGoldilocksConfig;
    type F = <C as GenericConfig<D>>::F;

    let mut config = CircuitConfig::standard_recursion_config();
    config.num_wires = 60;
    config.num_routed_wires = 60;

    let mut builder = CircuitBuilder::<F, D>::new(config);
    let targets = make_circuits(&mut builder, len as u64);
    let mut pw = PartialWitness::new();


    let mut timing = TimingTree::new("generate pw", Level::Info);
    timed!(timing, "witness generation", {
        for i in 0..len {
            pw.set_bool_target(targets.message[i], msg_bits[i]);
        }
    });
    timing.print();

   
    let expected_res = array_to_bits(hash.as_slice());
    for i in 0..expected_res.len() {
        if expected_res[i] {
            builder.assert_one(targets.digest[i].target);
        } else {
            builder.assert_zero(targets.digest[i].target);
        }
    }

    println!(
        "Constructing proof with {} gates",
        builder.num_gates(),
    );
    let data = builder.build::<C>();
    let mut timing = TimingTree::new("prove", Level::Debug);

    let original_proof = prove(&data.prover_only, &data.common, pw, &mut timing).unwrap();
    timing.print();

    let timing = TimingTree::new("verify", Level::Debug);
    data.verify(original_proof.clone()).unwrap();
    timing.print();

    println!(
        "Proof size: {}",
        original_proof.to_bytes().len(),
    );

    let mut config = CircuitConfig::standard_recursion_config();
    config.fri_config.rate_bits = 6;
    config.fri_config.num_query_rounds = 14;
    recursive_proof::<F, C, C, D>(original_proof, data.verifier_only, data.common, &config)?;

    Ok(())
}

fn recursive_proof<
    F: RichField + Extendable<D>,
    C: GenericConfig<D, F = F>,
    InnerC: GenericConfig<D, F = F>,
    const D: usize,
>(
    inner_proof: ProofWithPublicInputs<F, InnerC, D>,
    inner_vd: VerifierOnlyCircuitData<InnerC, D>,
    inner_cd: CommonCircuitData<F, D>,
    config: &CircuitConfig,
) -> Result<()>
where
    InnerC::Hasher: AlgebraicHasher<F>,
    [(); C::Hasher::HASH_SIZE]:,
{
    let mut builder = CircuitBuilder::<F, D>::new(config.clone());
    let mut pw = PartialWitness::new();
    let timing = TimingTree::new("generate recursive pw", Level::Info);
    let pt = builder.add_virtual_proof_with_pis::<InnerC>(&inner_cd);
    pw.set_proof_with_pis_target(&pt, &inner_proof);

    let inner_data = VerifierCircuitTarget {
        constants_sigmas_cap: builder.add_virtual_cap(inner_cd.config.fri_config.cap_height),
        circuit_digest: builder.add_virtual_hash(),
    };
    pw.set_cap_target(
        &inner_data.constants_sigmas_cap,
        &inner_vd.constants_sigmas_cap,
    );
    pw.set_hash_target(inner_data.circuit_digest, inner_vd.circuit_digest);
    timing.print();

    builder.register_public_inputs(inner_data.circuit_digest.elements.as_slice());
    for i in 0..builder.config.fri_config.num_cap_elements() {
        builder.register_public_inputs(&inner_data.constants_sigmas_cap.0[i].elements);
    }
    builder.verify_proof::<InnerC>(&pt, &inner_data, &inner_cd);

    // builder.print_gate_counts(0);
    println!(
        "Constructing recursive proof with {} gates",
        builder.num_gates(),
    );

    let data = builder.build::<C>();

    let mut timing = TimingTree::new("prove recursive", Level::Debug);
    let proof = prove(&data.prover_only, &data.common, pw, &mut timing)?;
    timing.print();
    
    let timing = TimingTree::new("verify recursive proof", Level::Debug);
    data.verify(proof.clone())?;
    println!(
        "Recursive proof size: {}",
        proof.to_bytes().len(),
    );
    timing.print();

    Ok(())
}

fn main() -> Result<()> {
    // Initialize logging
    let mut builder = env_logger::Builder::from_default_env();
    builder.format_timestamp(None);
    builder.filter_level(LevelFilter::Debug);
    builder.try_init()?;

    let args: Vec<String> = env::args().collect();
    let message_size = args[1].parse::<usize>().unwrap();
    if message_size < 1024 {
        println!("message size: {} B", message_size);
    } else {
        println!("message size: {} KB", message_size / 1024);
    }

    let mut msg = vec![0; message_size];
    for i in 0..message_size - 1 {
        msg[i] = 0 as u8;
    }

    prove_sha256(&msg)
}
