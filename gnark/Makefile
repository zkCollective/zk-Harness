SHELL = zsh

directory = github.com/tumberger/zk-compilers/gnark

.PHONY = test

all: test-toy

test-toy:
	go test $(directory)/circuits/toy/cubic
	go test $(directory)/circuits/toy/exponentiate