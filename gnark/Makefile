SHELL = zsh

directory = github.com/tumberger/zk-compilers/gnark

.PHONY = test

all: test-toy

test-toy:
	go test $(directory)/circuits/toy/cubic
	go test $(directory)/circuits/toy/expo
	go test $(directory)/circuits/toy/exponentiate

test-prf:
	go test $(directory)/circuits/prf/mimc

bench-arithmetic-test:
	./gnark arithmetic --input="none" --operation="add"
	./gnark arithmetic --input="none" --operation="sub"
	./gnark arithmetic --input="none" --operation="mul"
	./gnark arithmetic --input="none" --operation="div"
	./gnark arithmetic --input="none" --operation="exp"