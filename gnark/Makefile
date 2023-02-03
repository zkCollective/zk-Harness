SHELL = zsh

directory = github.com/tumberger/zk-compilers/gnark

.PHONY = test

all: test-toy

test-toy:
	go test $(directory)/toy/cubic
	go test $(directory)/toy/exponentiate