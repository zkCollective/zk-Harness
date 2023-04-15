#[macro_use]
extern crate criterion;

use ark_ec::pairing::Pairing;
use ark_ec::{CurveGroup, ScalarMul};
use ark_ff::{FftField, Field};
use ark_poly::univariate::DensePolynomial;
use ark_poly::{DenseUVPolynomial, EvaluationDomain, GeneralEvaluationDomain};
use ark_std::test_rng;
use ark_std::UniformRand;
use criterion::measurement::Measurement;
use criterion::BenchmarkGroup;
use criterion::{BenchmarkId, Criterion};

fn bench_msm<G: CurveGroup, M: Measurement>(c: &mut BenchmarkGroup<'_, M>, group_name: &str) {
    let rng = &mut test_rng();

    for logsize in 1..=21 {
        let size = 1 << logsize;

        // Dynamically control sample size so that big MSMs don't bench eternally
        if logsize > 20 {
            c.sample_size(10);
        }

        let scalars = (0..size)
            .map(|_| G::ScalarField::rand(rng))
            .collect::<Vec<_>>();
        let gs = (0..size)
            .map(|_| G::rand(rng).into_affine())
            .collect::<Vec<_>>();
        c.bench_with_input(
            BenchmarkId::new(format!("msm/{}", group_name), size),
            &logsize,
            |b, _| b.iter(|| G::msm(&gs, &scalars)),
        );
    }
}

fn bench_multi_pairing<P: Pairing, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let rng = &mut test_rng();
    for logsize in 1..=18 {
        let size = 1 << logsize;
        let g1s = (0..size)
            .map(|_| P::G1::rand(rng).into_affine())
            .collect::<Vec<_>>();
        let g2s = (0..size)
            .map(|_| P::G2::rand(rng).into_affine())
            .collect::<Vec<_>>();
        c.bench_with_input(BenchmarkId::new("msm/Gt", size), &logsize, |b, _| {
            b.iter(|| P::multi_pairing(&g1s, &g2s))
        });
    }
}

fn bench_sum_of_products<F: Field, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let rng = &mut test_rng();
    c.bench_function("msm/ff", |b| {
        const SIZE: usize = 256;
        let lhs: [F; SIZE] = (0..SIZE)
            .map(|_| F::rand(rng))
            .collect::<Vec<_>>()
            .try_into()
            .unwrap();
        let rhs: [F; SIZE] = (0..SIZE)
            .map(|_| F::rand(rng))
            .collect::<Vec<_>>()
            .try_into()
            .unwrap();
        b.iter(|| F::sum_of_products(&lhs, &rhs))
    });
}

fn bench_mul<G: ScalarMul, M: Measurement>(c: &mut BenchmarkGroup<'_, M>, group_name: &str) {
    let rng = &mut test_rng();
    c.bench_function(format!("mul_{}", group_name), |b| {
        const SIZE: usize = 256;
        let lhs = G::rand(rng);
        let rhs = G::ScalarField::rand(rng);
        b.iter(|| lhs * &rhs)
    });
}

fn bench_pairing<P: Pairing, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("pairing", |r| {
        let a = P::G1::rand(&mut rng).into_affine();
        let b = P::G2::rand(&mut rng).into_affine();
        r.iter(|| P::pairing(a, b))
    });
}

fn bench_fft<F: FftField, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    for logsize in 1..=21 {
        let degree = 1 << logsize;
        let domain = GeneralEvaluationDomain::<F>::new(degree).unwrap();
        c.bench_with_input(BenchmarkId::new("fft", degree), &logsize, |b, _| {
            let a = DensePolynomial::<F>::rand(degree, &mut rng)
                .coeffs()
                .to_vec();
            b.iter(|| domain.fft(&a))
        });
    }
}

fn bench_bls12_381(c: &mut Criterion) {
    use ark_bls12_381::{Bls12_381, Fr, G1Projective, G2Projective};
    type Gt = ark_ec::pairing::PairingOutput<Bls12_381>;

    let mut group = c.benchmark_group("bls12_381");
    bench_msm::<G1Projective, _>(&mut group, "G1");
    bench_msm::<G2Projective, _>(&mut group, "G2");
    bench_mul::<Gt, _>(&mut group, "Gt");
    bench_multi_pairing::<Bls12_381, _>(&mut group);
    bench_pairing::<Bls12_381, _>(&mut group);
    bench_sum_of_products::<Fr, _>(&mut group);
    bench_fft::<Fr, _>(&mut group);
    group.finish();
}

fn bench_bls12_377(c: &mut Criterion) {
    use ark_bls12_377::{Bls12_377, Fr, G1Projective, G2Projective};
    type Gt = ark_ec::pairing::PairingOutput<Bls12_377>;

    let mut group = c.benchmark_group("bls12_377");
    bench_msm::<G1Projective, _>(&mut group, "G1");
    bench_msm::<G2Projective, _>(&mut group, "G2");
    bench_mul::<Gt, _>(&mut group, "Gt");
    bench_multi_pairing::<Bls12_377, _>(&mut group);
    bench_pairing::<Bls12_377, _>(&mut group);
    bench_sum_of_products::<Fr, _>(&mut group);
    bench_fft::<Fr, _>(&mut group);
    group.finish();
}

fn bench_curve25519(c: &mut Criterion) {
    use ark_curve25519::{EdwardsProjective as G, Fr};
    let mut group = c.benchmark_group("curve25519");
    bench_msm::<G, _>(&mut group, "G1");
    bench_sum_of_products::<Fr, _>(&mut group);
    bench_fft::<Fr, _>(&mut group);
    group.finish();
}

fn bench_secp256k1(c: &mut Criterion) {
    use ark_secp256k1::{Fr, Projective as G};
    let mut group = c.benchmark_group("secp256k1");
    bench_msm::<G, _>(&mut group, "G1");
    bench_sum_of_products::<Fr, _>(&mut group);
    bench_fft::<Fr, _>(&mut group);
    group.finish();
}

fn bench_pallas(c: &mut Criterion) {
    use ark_pallas::{Fr, Projective as G};
    let mut group = c.benchmark_group("pallas");
    bench_msm::<G, _>(&mut group, "G1");
    bench_sum_of_products::<Fr, _>(&mut group);
    bench_fft::<Fr, _>(&mut group);
    group.finish();
}

fn bench_vesta(c: &mut Criterion) {
    use ark_pallas::{Fr, Projective as G};
    let mut group = c.benchmark_group("vesta");
    bench_msm::<G, _>(&mut group, "G1");
    bench_sum_of_products::<Fr, _>(&mut group);
    bench_fft::<Fr, _>(&mut group);
    group.finish();
}

criterion_group!(
    benches,
    bench_bls12_381,
    bench_bls12_377,
    bench_curve25519,
    bench_secp256k1,
    bench_pallas,
    bench_vesta
);
criterion_main!(benches);
