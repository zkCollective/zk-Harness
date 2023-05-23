use criterion::{black_box, criterion_group, criterion_main, Criterion};
use pasta_curves::group::Group;
use std::any::type_name;

use criterion::measurement::Measurement;
use criterion::{BenchmarkGroup, BenchmarkId};
use ff::{PrimeField, Field};
use pairing::Engine;
use halo2curves::{bn256, pairing, CurveExt};
use rand_core::OsRng;

// Benchmark Addition in Scalar Field
fn bench_add_ff<F: PrimeField, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let a: F = F::random(OsRng);
    let b: F = F::random(OsRng); 
    
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
    let a: F = F::random(OsRng);
    let b: F = F::random(OsRng); 
    let mut lhs = black_box(a);
    let rhs = black_box(&b);
    c.bench_function("mul_ff", |b| {
        b.iter(|| {
            F::mul_assign(&mut lhs, rhs)
        });
    });
}

// Benchmark Addition in Elliptic Curve Group
fn bench_add_ec<G: CurveExt, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let lhs = G::random(OsRng);
    let rhs = G::random(OsRng);
    let mut lhs = black_box(lhs);
    let rhs = black_box(&rhs);
    let group_name = type_name::<G>();
    c.bench_function(BenchmarkId::new(group_name, "add"), |b| {
        b.iter(|| {
            lhs.add_assign(rhs);
        });
    });
}

// TODO - Multiplication in EC

// Sampling Scalars for MSM
fn get_data(n: usize) -> (Vec<bn256::G1Affine>, Vec<bn256::Fr>) {
    const MAX_N: usize = 1 << 22;
    assert!(n <= MAX_N);

    let (points, scalars): (Vec<_>, Vec<_>) = (0..n)
        .map(|_| {
            let point = bn256::G1Affine::random(OsRng);
            let scalar = bn256::Fr::random(OsRng);
            (point, scalar)
        })
        .unzip();

    (points, scalars)
}


// The MSM algorithm of halo2_curve can be found here: https://github.com/privacy-scaling-explorations/halo2curves/pull/29
// FIXME - Fixed to BN256
fn bench_msm(group: &mut BenchmarkGroup<'_, criterion::measurement::WallTime>) {
    // const SAMPLES: usize = 1 << 3;
    const MAX_SIZE: usize = 20;

    // let (min_k, max_k) = (4, 20);
    let (points, scalars) = get_data(1 << MAX_SIZE);

    for logsize in 1..=MAX_SIZE {
        // Dynamically control sample size so that big MSMs don't bench eternally
        if logsize > 20 {
            group.sample_size(10);
        }

        let size = 1 << logsize;

        let scalars = &scalars[..size];
        let points = &points[..size];
        let mut r1 = bn256::G1::identity();

        group.bench_with_input(BenchmarkId::new("G1", size), &size, |b, _| {
            b.iter(|| {
                bn256::msm::MSM::evaluate(scalars, points, &mut r1);
            });
        });
    }
}

// Benchmark pairing
// FIXME - Fixed to BN256
fn bench_pairing< M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let g1 = bn256::G1::generator();
    let g2 = bn256::G2::generator();
    c.bench_function("pairing", |r| {
        r.iter(|| bn256::Bn256::pairing(&bn256::G1Affine::from(g1), &bn256::G2Affine::from(g2)))
    });
}

fn bench_bn256(c: &mut Criterion) {
    let mut group = c.benchmark_group("bn256");
    bench_add_ff::<bn256::Fr, _>(&mut group);
    bench_mul_ff::<bn256::Fr,_>(&mut group);
    bench_add_ec::<bn256::G1, _>(&mut group);
    bench_msm(&mut group);
    bench_pairing::<_>(&mut group);
}

criterion_group!(
    benches,
    bench_bn256
);
criterion_main!(benches);
