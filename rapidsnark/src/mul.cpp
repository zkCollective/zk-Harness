#include <stdio.h>
#include <stdlib.h>
#include "fr.hpp"
#include <chrono>

int main(int argc, char **argv) {

    int N = atoi(argv[1]);

    Fr_init();

    FrElement a;
    a.type = Fr_LONGMONTGOMERY;
    for (int i=0; i<Fr_N64; i++) {
        a.longVal[i] = 0xAAAAAAAA;
    }

    std::chrono::steady_clock::time_point begin = std::chrono::steady_clock::now();
    for (long long i=0; i<N; i++) {
        Fr_mul(&a, &a, &a);
    }
    std::chrono::steady_clock::time_point end = std::chrono::steady_clock::now();
    printf("%lld", (std::chrono::duration_cast<std::chrono::nanoseconds> (end - begin).count())/N);
}
