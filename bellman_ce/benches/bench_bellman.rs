use criterion::{black_box, criterion_group, criterion_main, Criterion};
use rand::{Rand};
use std::sync::Arc;

use criterion::measurement::Measurement;
use criterion::{BenchmarkGroup, BenchmarkId};
use ff_ce::{Field, PrimeField, PrimeFieldRepr, SqrtField};
use pairing_ce::{bls12_381::*, bn256::*, GenericCurveAffine, GenericCurveProjective, Engine};
use bellman_ce::worker::{Worker, WorkerFuture};
use bellman_ce::source::FullDensity;
use bellman_ce::multiexp::*;
use futures::*; 

// Benchmark Addition in Scalar Field
fn bench_add_ff<G: GenericCurveProjective, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    let lhs = G::rand(&mut rng);
    let rhs = G::rand(&mut rng);
    let mut lhs = black_box(lhs);
    let rhs = black_box(&rhs);
    c.bench_function("add_ff", |b| {
        b.iter(|| {
            lhs.add_assign(rhs);
        });
    });
}

// Benchmark Multiplication in Scalar Field
fn bench_mul_ff<G, M>(c: &mut BenchmarkGroup<'_, M>)
where
    G: GenericCurveProjective,
    M: Measurement,
{
    let mut rng = rand::thread_rng();
    let lhs = G::rand(&mut rng);
    let rhs_scalar = <G as GenericCurveProjective>::Scalar::rand(&mut rng);
    let mut lhs = black_box(lhs);
    let rhs_scalar = black_box(&rhs_scalar);
    c.bench_function("mul_ff", |b| {
        b.iter(|| {
            lhs.mul_assign(rhs_scalar.into_repr());
        });
    });
}

// The MSM algorithm of bellman_ce can be found here: https://github.com/matter-labs/bellman/blob/dev/src/multiexp.rs#L60
fn bench_new_multexp<G>(group: &mut BenchmarkGroup<'_, criterion::measurement::WallTime>)
where
    G: Engine,
    <G as Engine>::G1: Rand,
{
    const SAMPLES: usize = 1 << 2;

    let rng = &mut rand::thread_rng();
    let v = Arc::new((0..SAMPLES).map(|_| G::Fr::rand(rng).into_repr()).collect::<Vec<_>>());
    let g = Arc::new((0..SAMPLES).map(|_| G::G1::rand(rng).into_affine()).collect::<Vec<_>>());

    let pool = Worker::new();

    let v = black_box(v);
    let g = black_box(g);

    group.bench_function("multiexp", |b| {
        b.iter(|| {
            let result = multiexp(&pool, (g.clone(), 0), FullDensity, v.clone());
            // Call block_on to actually run the future and obtain the result.
            futures::executor::block_on(result).unwrap();
        });
    });
}

fn bench_pairing<P: pairing_ce::Engine, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("pairing", |r| {
        let a = P::G1::rand(&mut rng).into_affine();
        let b = P::G2::rand(&mut rng).into_affine();
        r.iter(|| P::pairing(a, b))
    });
}

fn bench_new_multexp_test<G>(group: &mut BenchmarkGroup<'_, criterion::measurement::WallTime>)
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
        let vec_a: Vec<_> = (0..size).map(|_|  G::Fr::rand(rng).into_repr()).collect::<Vec<_>>();

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


fn bench_bls12_381(c: &mut Criterion) {
    let mut group = c.benchmark_group("bls12_381");
    bench_add_ff::<pairing_ce::bls12_381::G1, _>(&mut group);
    bench_mul_ff::<pairing_ce::bls12_381::G1, _>(&mut group);
    bench_new_multexp_test::<Bls12>(&mut group);
    bench_pairing::<Bls12, _>(&mut group);
}

fn bench_bn256(c: &mut Criterion) {
    let mut group = c.benchmark_group("bn256");
    bench_add_ff::<pairing_ce::bn256::G1, _>(&mut group);
    bench_mul_ff::<pairing_ce::bn256::G1, _>(&mut group);
    bench_new_multexp_test::<Bn256>(&mut group);
    bench_pairing::<Bn256, _>(&mut group);
}

criterion_group!(
    benches,
    bench_bls12_381,
    bench_bn256
);
criterion_main!(benches);
