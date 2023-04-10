# TIME=/usr/bin/time
TIME=/usr/local/bin/gtime
which $TIME

RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 64
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 128
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 256
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 512
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 1024
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 2048
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 4096
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 8192
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 16384
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 32768
RUSTFLAGS=-Ctarget-cpu=native $TIME -f "Peak memory: %M kb CPU usage: %P" cargo +nightly run --release --package plonky2_sha256 --bin plonky2_sha256 -- 65536
