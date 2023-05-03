use criterion::{black_box, criterion_group, criterion_main, Criterion};
use rand::{Rand};
use rand::{Rng, XorShiftRng, SeedableRng};
use std::sync::Arc;
use std::any::type_name;

use criterion::measurement::Measurement;
use criterion::{BenchmarkGroup, BenchmarkId};
use ff_ce::{PrimeField};
use pairing_ce::{bls12_381::*, bn256::*, GenericCurveProjective, Engine};
use bellman_ce::worker::{Worker};
use bellman_ce::source::FullDensity;
use bellman_ce::{multiexp::*};

// Benchmark Addition in Scalar Field
fn bench_add_ff<F: PrimeField, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    // let mut rng = rand::thread_rng();
    let rng = &mut XorShiftRng::from_seed([0x5dbe6259, 0x8d313d76, 0x3237db17, 0xe5bc0654]);
    let a: F = rng.gen();
    let b: F = rng.gen(); 
    
    let mut lhs = black_box(a);
    let rhs = black_box(&b);
    c.bench_function("add_ff", |b| {
        b.iter(|| {
            F::add_assign(&mut lhs, rhs)
        });
    });
}

// Benchmark Addition in Scalar FIeld
fn bench_mul_ff<F: PrimeField, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let rng = &mut XorShiftRng::from_seed([0x5dbe6259, 0x8d313d76, 0x3237db17, 0xe5bc0654]);
    let a: F = rng.gen();
    let b: F = rng.gen(); 
    
    let mut lhs = black_box(a);
    let rhs = black_box(&b);
    c.bench_function("mul_ff", |b| {
        b.iter(|| {
            F::mul_assign(&mut lhs, rhs)
        });
    });
}

// Benchmark Addition in Elliptic Curve Group
fn bench_add_ec<G: GenericCurveProjective, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    let lhs = G::rand(&mut rng);
    let rhs = G::rand(&mut rng);
    let mut lhs = black_box(lhs);
    let rhs = black_box(&rhs);
    let group_name = type_name::<G>();
    c.bench_function(BenchmarkId::new(group_name, "add"), |b| {
        b.iter(|| {
            lhs.add_assign(rhs);
        });
    });
}

// Benchmark Multiplication in Elliptic Curve Group
fn bench_mul_ec<G, M>(c: &mut BenchmarkGroup<'_, M>)
where
    G: GenericCurveProjective,
    M: Measurement,
{
    let mut rng = rand::thread_rng();
    let lhs = G::rand(&mut rng);
    let rhs_scalar = <G as GenericCurveProjective>::Scalar::rand(&mut rng);
    let mut lhs = black_box(lhs);
    let rhs_scalar = black_box(&rhs_scalar);
    let group_name = type_name::<G>();
    c.bench_function(BenchmarkId::new(group_name, "mul"), |b| {
        b.iter(|| {
            lhs.mul_assign(rhs_scalar.into_repr());
        });
    });
}

// The MSM algorithm of bellman_ce can be found here: https://github.com/matter-labs/bellman/blob/dev/src/multiexp.rs#L60
// Currently, this uses a fork of bellman_ce, as multiexp as not exposed by bellman_ce
// Uses a Worker pool for multi threading
fn bench_msm<G>(group: &mut BenchmarkGroup<'_, criterion::measurement::WallTime>)
where
    G: Engine,
    <G as Engine>::G1: Rand,
    <G as Engine>::G2: Rand,
{
    const SAMPLES: usize = 1 << 3;
    const MAX_SIZE: usize = 3;

    let rng = &mut rand::thread_rng();
    let v = Arc::new((0..SAMPLES).map(|_| G::Fr::rand(rng).into_repr()).collect::<Vec<_>>());
    let g1 = Arc::new((0..SAMPLES).map(|_| G::G1::rand(rng).into_affine()).collect::<Vec<_>>());
    let g2 = Arc::new((0..SAMPLES).map(|_| G::G2::rand(rng).into_affine()).collect::<Vec<_>>());

    let pool = Worker::new();

    let v = black_box(v);
    let g1 = black_box(g1);
    let g2 = black_box(g2);

    for logsize in 1..=MAX_SIZE {
        // Dynamically control sample size so that big MSMs don't bench eternally
        if logsize > 20 {
            group.sample_size(10);
        }

        let size = 1 << logsize;

        group.bench_with_input(BenchmarkId::new("G1", size), &size, |b, _| {
            b.iter(|| {
                let result = multiexp(&pool, (g1.clone(), 0), FullDensity, v.clone());
                futures::executor::block_on(result).unwrap();
            });
        });

        group.bench_with_input(BenchmarkId::new("G2", size), &size, |b, _| {
            b.iter(|| {
                let result = multiexp(&pool, (g2.clone(), 0), FullDensity, v.clone());
                futures::executor::block_on(result).unwrap();
            });
        });
    }
}

// Benchmark pairing
fn bench_pairing<P: pairing_ce::Engine, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    let a = P::G1::rand(&mut rng).into_affine();
    let b = P::G2::rand(&mut rng).into_affine();
    c.bench_function("pairing", |r| {
        r.iter(|| P::pairing(a, b))
    });
}


fn bench_bls12_381(c: &mut Criterion) {
    let mut group = c.benchmark_group("bls12_381");
    bench_add_ff::<pairing_ce::bls12_381::fr::Fr, _>(&mut group);
    bench_mul_ff::<pairing_ce::bls12_381::fr::Fr, _>(&mut group);
    bench_add_ec::<pairing_ce::bls12_381::G1, _>(&mut group);
    bench_mul_ec::<pairing_ce::bls12_381::G1, _>(&mut group);
    bench_msm::<Bls12>(&mut group);
    bench_pairing::<Bls12, _>(&mut group);
}

fn bench_bn256(c: &mut Criterion) {
    let mut group = c.benchmark_group("bn256");
    bench_add_ff::<pairing_ce::bn256::fr::Fr, _>(&mut group);
    bench_mul_ff::<pairing_ce::bn256::fr::Fr,_>(&mut group);
    bench_add_ec::<pairing_ce::bn256::G1, _>(&mut group);
    bench_mul_ec::<pairing_ce::bn256::G1, _>(&mut group);
    bench_msm::<Bn256>(&mut group);
    bench_pairing::<Bn256, _>(&mut group);
}

criterion_group!(
    benches,
    bench_bls12_381,
    bench_bn256
);
criterion_main!(benches);
