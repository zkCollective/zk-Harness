#include <stdio.h>
#include <stdlib.h>
#include "alt_bn128.hpp"
#include <time.h>
#include <chrono>

using namespace AltBn128;

__uint128_t g_lehmer64_state = 0xAAAAAAAAAAAAAAAALL;

// Fast random generator
// https://lemire.me/blog/2019/03/19/the-fastest-conventional-random-number-generator-that-can-pass-big-crush/

uint64_t lehmer64() {
  g_lehmer64_state *= 0xda942042e4dd58b5LL;
  return g_lehmer64_state >> 64;
}

int main(int argc, char **argv) {

    int X = atoi(argv[1]);
    int N = atoi(argv[2]);

    uint8_t *scalars = new uint8_t[X*32];
    G1PointAffine *bases = new G1PointAffine[X];

    // random scalars
    for (int i=0; i<X*4; i++) {
        *((uint64_t *)(scalars + i*8)) = lehmer64();
    }

    G1.copy(bases[0], G1.one());
    G1.copy(bases[1], G1.one());

    for (int i=2; i<X; i++) {
        G1.add(bases[i], bases[i-1], bases[i-2]);
    }

    double start, end;
    double cpu_time_used;

    G1Point p1;

#ifdef COUNT_OPS
    G1.resetCounters();
#endif

    std::chrono::steady_clock::time_point chrono_begin = std::chrono::steady_clock::now();
    for (long long i=0; i<N; i++) {
        G1.multiMulByScalar(p1, bases, (uint8_t *)scalars, 32, X);
    }
    std::chrono::steady_clock::time_point chrono_end = std::chrono::steady_clock::now();
    printf("%lld", (std::chrono::duration_cast<std::chrono::nanoseconds> (chrono_end - chrono_begin).count())/N);

#ifdef COUNT_OPS
    G1.printCounters();
#endif
}
