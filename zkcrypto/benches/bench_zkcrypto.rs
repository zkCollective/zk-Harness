use criterion::{
    black_box, criterion_group, criterion_main, measurement::Measurement, BenchmarkGroup,
    BenchmarkId, Criterion,
};
use group::{ff::Field, Curve, Group};
use pairing::{MillerLoopResult, MultiMillerLoop};

fn bench_add_ff<G: Group, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("add_ff", |b| {
        let lhs = G::Scalar::random(&mut rng);
        let rhs = G::Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) + black_box(rhs))
    });
}

fn bench_mul_ff<G: Group, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("mul_ff", |b| {
        let lhs = G::Scalar::random(&mut rng);
        let rhs = G::Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) * black_box(rhs))
    });
}

fn bench_invert<G: Group, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("invert", |b| {
        let a = G::Scalar::random(&mut rng);
        b.iter(|| a.invert().unwrap());
    });
}

fn bench_add_ec<G: Group, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("add_G1", |r| {
        let a = G::random(&mut rng);
        let b = G::random(&mut rng);
        r.iter(|| a.add(b))
    });
}

fn bench_dbl_ec<G: Group, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("dbl_G1", |r| {
        let a = G::random(&mut rng);
        r.iter(|| a.double())
    });
}

fn bench_mul_ec<G: Group, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("mul_G1", |b| {
        let lhs = G::random(&mut rng);
        let rhs = G::Scalar::random(&mut rng);
        b.iter(|| black_box(lhs) * black_box(rhs))
    });
}

fn bench_pairing<P: pairing::Engine, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();
    c.bench_function("pairing", |r| {
        let a = P::G1::random(&mut rng).to_affine();
        let b = P::G2::random(&mut rng).to_affine();
        r.iter(|| P::pairing(&a, &b))
    });
}

fn bench_pairing_product<P: pairing::Engine + MultiMillerLoop, M: Measurement>(
    c: &mut BenchmarkGroup<'_, M>,
) {
    let mut rng = rand::thread_rng();
    for d in 1..=10 {
        let size = 1 << d;
        let mut v: Vec<(P::G1Affine, P::G2Prepared)> = Vec::new();
        for _ in 0..size {
            let g1 = P::G1::random(&mut rng).to_affine();
            let g2 = P::G2Prepared::from(P::G2::random(&mut rng).to_affine());
            v.push((g1, g2));
        }

        let mut v_ref: Vec<(&P::G1Affine, &P::G2Prepared)> = Vec::new();
        for i in 0..size {
            v_ref.push((&v[i].0, &v[i].1));
        }

        c.bench_with_input(BenchmarkId::new("msm/Gt", size), &d, |b, _| {
            b.iter(|| P::multi_miller_loop(&v_ref).final_exponentiation())
        });
    }
}

fn bench_bls12_381(c: &mut Criterion) {
    let mut group = c.benchmark_group("bls12_381");
    bench_add_ff::<bls12_381::G1Projective, _>(&mut group);
    bench_mul_ff::<bls12_381::G1Projective, _>(&mut group);
    bench_invert::<bls12_381::G1Projective, _>(&mut group);
    bench_add_ec::<bls12_381::G1Projective, _>(&mut group);
    bench_dbl_ec::<bls12_381::G1Projective, _>(&mut group);
    bench_mul_ec::<bls12_381::G1Projective, _>(&mut group);
    bench_pairing::<bls12_381::Bls12, _>(&mut group);
    bench_pairing_product::<bls12_381::Bls12, _>(&mut group);
    group.finish();
}

fn bench_jubjub(c: &mut Criterion) {
    let mut group = c.benchmark_group("jubjub");
    bench_add_ff::<jubjub::ExtendedPoint, _>(&mut group);
    bench_mul_ff::<jubjub::ExtendedPoint, _>(&mut group);
    bench_invert::<jubjub::ExtendedPoint, _>(&mut group);
    bench_add_ec::<jubjub::ExtendedPoint, _>(&mut group);
    bench_dbl_ec::<jubjub::ExtendedPoint, _>(&mut group);
    bench_mul_ec::<jubjub::ExtendedPoint, _>(&mut group);
    group.finish();
}

criterion_group! {
    name = zkcrypto_benchmarks;
    config = Criterion::default();
    targets = bench_bls12_381, bench_jubjub
}

criterion_main!(zkcrypto_benchmarks);
