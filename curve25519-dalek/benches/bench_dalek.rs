#[macro_use]
extern crate criterion;

use criterion::{black_box, BenchmarkId, Criterion};
use curve25519_dalek::ristretto::RistrettoPoint;
use curve25519_dalek::scalar::Scalar;
use curve25519_dalek::traits::MultiscalarMul;
use rand::rngs::OsRng;

fn bench_add_ff(c: &mut Criterion) {
    let mut rng = OsRng;
    c.bench_function("curve25519/add_ff", |b| {
        let lhs = Scalar::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) + black_box(rhs))
    });
}

fn bench_invert_ff(c: &mut Criterion) {
    let mut rng = OsRng;
    c.bench_function("curve25519/invert", |b| {
        let lhs = Scalar::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) + black_box(rhs))
    });
}

fn bench_mul_ec(c: &mut Criterion) {
    let mut rng = OsRng;
    c.bench_function("curve25519/mul_G1", |b| {
        let lhs = RistrettoPoint::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) * black_box(rhs))
    });
}

fn bench_mul_ff(c: &mut Criterion) {
    let mut rng = OsRng;
    c.bench_function("curve25519/mul_ff", |b| {
        let lhs = Scalar::random(&mut rng);
        let rhs = Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) * black_box(rhs))
    });
}

fn bench_msm(c: &mut Criterion) {
    let mut rng = OsRng;
    for logsize in 1..=21 {
        let mut group = c.benchmark_group("curve25519/msm");
        let size = 1 << logsize;

        // Dynamically control sample size so that big MSMs don't bench eternally
        if logsize > 20 {
            group.sample_size(10);
        }

        group.bench_with_input(BenchmarkId::new("G1", size), &size, |b, &size| {
            let scalars = (0..size)
                .map(|_| Scalar::random(&mut rng))
                .collect::<Vec<_>>();
            let bases = (0..size)
                .map(|_| RistrettoPoint::random(&mut rng))
                .collect::<Vec<_>>();
            b.iter(|| RistrettoPoint::multiscalar_mul(scalars.iter(), bases.iter()))
        });
    }
}

criterion_group!(
    msm_benchmarks,
    bench_mul_ff,
    bench_mul_ec,
    bench_add_ff,
    bench_msm,
    bench_invert_ff
);

criterion_main!(msm_benchmarks);
