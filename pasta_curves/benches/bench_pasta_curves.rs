use criterion::{
    black_box, criterion_group, criterion_main, measurement::Measurement, BenchmarkGroup,
    BenchmarkId, Criterion,
};
use group::{ff::Field, Curve, Group};
use halo2_proofs::arithmetic::{best_fft, best_multiexp};
use pasta_curves::{arithmetic::CurveAffine, pallas, vesta};

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

fn bench_msm<C: CurveAffine, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();

    for logsize in 1..=21 {
        let size = 1 << logsize;

        // Dynamically control sample size so that big MSMs don't bench eternally
        if logsize > 20 {
            c.sample_size(10);
        }

        c.bench_with_input(BenchmarkId::new("msm", size), &size, |b, &size| {
            let scalars = (0..size)
                .map(|_| C::Scalar::random(&mut rng))
                .collect::<Vec<_>>();
            let bases = (0..size)
                .map(|_| C::Curve::random(&mut rng).to_affine())
                .collect::<Vec<_>>();
            b.iter(|| best_multiexp(&scalars, &bases))
        });
    }
}

fn bench_fft<Scalar: Field, M: Measurement>(c: &mut BenchmarkGroup<'_, M>) {
    let mut rng = rand::thread_rng();

    for logsize in 1..=21 {
        let degree = 1 << logsize;

        // Dynamically control sample size so that big FFTs don't bench eternally
        if logsize > 20 {
            c.sample_size(10);
        }

        c.bench_with_input(BenchmarkId::new("fft", degree), &degree, |b, &degree| {
            let mut scalars = (0..degree)
                .map(|_| Scalar::random(&mut rng))
                .collect::<Vec<_>>();
            let omega = Scalar::random(&mut rng);
            b.iter(|| best_fft(&mut scalars, omega, logsize as u32))
        });
    }
}

fn bench_pallas(c: &mut Criterion) {
    let mut group = c.benchmark_group("pallas");
    bench_add_ff::<pallas::Point, _>(&mut group);
    bench_mul_ff::<pallas::Point, _>(&mut group);
    bench_invert::<pallas::Point, _>(&mut group);
    bench_add_ec::<pallas::Point, _>(&mut group);
    bench_dbl_ec::<pallas::Point, _>(&mut group);
    bench_mul_ec::<pallas::Point, _>(&mut group);
    bench_msm::<pallas::Affine, _>(&mut group);
    bench_fft::<pallas::Scalar, _>(&mut group);
    group.finish();
}

fn bench_vesta(c: &mut Criterion) {
    let mut group = c.benchmark_group("vesta");
    bench_add_ff::<vesta::Point, _>(&mut group);
    bench_mul_ff::<vesta::Point, _>(&mut group);
    bench_invert::<vesta::Point, _>(&mut group);
    bench_add_ec::<vesta::Point, _>(&mut group);
    bench_dbl_ec::<vesta::Point, _>(&mut group);
    bench_mul_ec::<vesta::Point, _>(&mut group);
    bench_msm::<vesta::Affine, _>(&mut group);
    bench_fft::<vesta::Scalar, _>(&mut group);
    group.finish();
}

criterion_group! {
    name = pasta_curves_benchmarks;
    config = Criterion::default();
    targets = bench_pallas, bench_vesta
}

criterion_main!(pasta_curves_benchmarks);
