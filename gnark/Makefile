SHELL = zsh

directory = github.com/tumberger/zk-compilers/gnark

.PHONY = test

all: test-toy test-prf

test-toy:
	go test $(directory)/circuits/toy/cubic
	go test $(directory)/circuits/toy/exponentiate
	go test $(directory)/circuits/toy/emulate

test-prf:
	go test $(directory)/circuits/prf/mimc
	go test $(directory)/circuits/prf/sha256