#![allow(non_snake_case)]

use blstrs::{Bls12, G2Prepared, Gt};
use blstrs::{G1Affine, G1Projective, G2Projective, Scalar};
use criterion::*;
use group::ff::Field;
use group::{Curve, Group};
use pairing_lib::{MillerLoopResult, MultiMillerLoop, PairingCurveAffine};
use std::ops::Add;

fn bench_mul(c: &mut Criterion) {
    let mut rng = rand::thread_rng();

    c.bench_function("bls12_381/mul_ff", |b| {
        let lhs = Scalar::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) * black_box(rhs))
    });

    c.bench_function("bls12_381/mul_G1", |b| {
        let lhs = G1Projective::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) * black_box(rhs))
    });

    c.bench_function("bls12_381/mul_G2", |b| {
        let lhs = G2Projective::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) * black_box(rhs))
    });

    c.bench_function("bls12_381/mul_Gt", |b| {
        let lhs = Gt::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) * black_box(rhs))
    });
}

fn bench_add(c: &mut Criterion) {
    let mut rng = rand::thread_rng();

    c.bench_function("bls12_381/add_ff", |b| {
        let lhs = Scalar::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) + black_box(rhs))
    });

    c.bench_function("bls12_381/add_G1", |r| {
        let a = G1Projective::random(&mut rng);
        let b = G1Projective::random(&mut rng);
        r.iter(|| a.add(b))
    });

    c.bench_function("bls12_381/add_G2", |r| {
        let a = G2Projective::random(&mut rng);
        let b = G2Projective::random(&mut rng);
        r.iter(|| a.add(b))
    });

    c.bench_function("bls12_381/add_Gt", |r| {
        let a = Gt::random(&mut rng);
        let b = Gt::random(&mut rng);
        r.iter(|| a.add(b))
    });
}

fn bench_msm(c: &mut Criterion) {
    let mut rng = rand::thread_rng();

    let mut group = c.benchmark_group("bls12_381/msm");
    for logsize in 1..=21 {
        // Dynamically control sample size so that big MSMs don't bench eternally
        if logsize > 20 {
            group.sample_size(10);
        }

        let size = 1 << logsize;
        let vec_a: Vec<_> = (0..size).map(|_| Scalar::random(&mut rng)).collect();
        // G1 benchmarks
        let vec_B_G1: Vec<_> = (0..size).map(|_| G1Projective::random(&mut rng)).collect();
        group.bench_with_input(BenchmarkId::new("G1", size), &size, |b, _| {
            b.iter(|| G1Projective::multi_exp(&vec_B_G1, &vec_a));
        });

        // G2 benchmarks
        let vec_B_G2: Vec<_> = (0..size).map(|_| G2Projective::random(&mut rng)).collect();
        group.bench_with_input(BenchmarkId::new("G2", size), &size, |b, _| {
            b.iter(|| G2Projective::multi_exp(&vec_B_G2, &vec_a));
        });
    }

    group.finish()
}

fn bench_invert(c: &mut Criterion) {
    let mut rng = rand::thread_rng();
    c.bench_function("bls12_381/invert", |b| {
        let a = Scalar::random(&mut rng);
        b.iter(|| a.invert().unwrap());
    });
}

fn bench_pairing(c: &mut Criterion) {
    let mut rng = rand::thread_rng();
    c.bench_function("bls12_381/pairing", |r| {
        let a = G1Projective::random(&mut rng).to_affine();
        let b = G2Projective::random(&mut rng).to_affine();
        r.iter(|| a.pairing_with(&b))
    });
}

fn bench_multi_pairing(c: &mut Criterion) {
    let mut rng = rand::thread_rng();
    let mut group = c.benchmark_group("bls12_381/msm/Gt");
    for d in 1..=10 {
        let size = 1 << d;
        let mut v: Vec<(G1Affine, G2Prepared)> = Vec::new();
        for _ in 0..size {
            let g1 = G1Affine::from(G1Projective::random(&mut rng));
            let g2 = G2Prepared::from(G2Projective::random(&mut rng).to_affine());
            v.push((g1, g2));
        }

        let mut v_ref: Vec<(&G1Affine, &G2Prepared)> = Vec::new();
        for i in 0..size {
            v_ref.push((&v[i].0, &v[i].1));
        }

        group.bench_with_input(BenchmarkId::from_parameter(size), &d, |b, _| {
            b.iter(|| Bls12::multi_miller_loop(&v_ref).final_exponentiation())
        });
    }
}

criterion_group!(
    blstrs_benchmarks,
    bench_mul,
    bench_add,
    bench_msm,
    bench_invert,
    bench_pairing,
    bench_multi_pairing,
);

criterion_main!(blstrs_benchmarks);
