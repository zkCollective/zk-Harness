directory = github.com/tumberger/zk-compilers/gnark

.PHONY = test

all: test-toy test-prf

test-toy:
	go test $(directory)/circuits/toy/cubic
	go test $(directory)/circuits/toy/expo
	go test $(directory)/circuits/toy/exponentiate

test-prf:
	go test $(directory)/circuits/prf/mimc
	go test $(directory)/circuits/prf/sha256