OS_ARCH=$(uname -a)
TIME=/usr/local/bin/gtime
if [[ $OS_ARCH == *"Darwin"*"arm64"* ]]; then
  TIME=/opt/homebrew/bin/gtime
fi
which $TIME
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 1   # 1 SHAs = 2 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 2   # 2 SHAs = 4 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 4   # 4 SHAs = 8 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 8   # 8 SHAs = 16 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 16  # 16 SHAs = 32 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 32  # 32 SHAs = 64 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 64  # 64 SHAs = 128 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 128 # 128 SHAs = 256 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 256 # 256 SHAs = 512 Blocks
RUSTFLAGS=-Awarnings $TIME -f "peak mem %M, avg cpu %P" cargo +nightly run --release --example time_sha256_compression -- 512 # 512 SHAs = 1024 Blocks
