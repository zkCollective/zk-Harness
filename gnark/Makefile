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
	./gnark arithmetic --input="none" --operation="add" --field="bls12381"
	./gnark arithmetic --input="none" --operation="sub" --field="bls12381"
	./gnark arithmetic --input="none" --operation="mul" --field="bls12381"
	./gnark arithmetic --input="none" --operation="div" --field="bls12381"
	./gnark arithmetic --input="none" --operation="exp" --field="bls12381"
	./gnark arithmetic --input="none" --operation="add" --field="bls12377"
	./gnark arithmetic --input="none" --operation="sub" --field="bls12377"
	./gnark arithmetic --input="none" --operation="mul" --field="bls12377"
	./gnark arithmetic --input="none" --operation="div" --field="bls12377"
	./gnark arithmetic --input="none" --operation="exp" --field="bls12377"
	./gnark arithmetic --input="none" --operation="add" --field="bls24315"
	./gnark arithmetic --input="none" --operation="sub" --field="bls24315"
	./gnark arithmetic --input="none" --operation="mul" --field="bls24315"
	./gnark arithmetic --input="none" --operation="div" --field="bls24315"
	./gnark arithmetic --input="none" --operation="exp" --field="bls24315"
	./gnark arithmetic --input="none" --operation="add" --field="bw6633"
	./gnark arithmetic --input="none" --operation="sub" --field="bw6633"
	./gnark arithmetic --input="none" --operation="mul" --field="bw6633"
	./gnark arithmetic --input="none" --operation="div" --field="bw6633"
	./gnark arithmetic --input="none" --operation="exp" --field="bw6633"
	./gnark arithmetic --input="none" --operation="add" --field="bw6761"
	./gnark arithmetic --input="none" --operation="sub" --field="bw6761"
	./gnark arithmetic --input="none" --operation="mul" --field="bw6761"
	./gnark arithmetic --input="none" --operation="div" --field="bw6761"
	./gnark arithmetic --input="none" --operation="exp" --field="bw6761"