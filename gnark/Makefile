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
	./gnark arithmetic --input="none" --operation="add" --count=20
	./gnark arithmetic --input="none" --operation="sub" --count=20
	./gnark arithmetic --input="none" --operation="mul" --count=20
	./gnark arithmetic --input="none" --operation="div" --count=20
	./gnark arithmetic --input="none" --operation="exp" --count=20
	./gnark arithmetic --input="none" --operation="add" --field="bls12381" --count=20
	./gnark arithmetic --input="none" --operation="sub" --field="bls12381" --count=20
	./gnark arithmetic --input="none" --operation="mul" --field="bls12381" --count=20
	./gnark arithmetic --input="none" --operation="div" --field="bls12381" --count=20
	./gnark arithmetic --input="none" --operation="exp" --field="bls12381" --count=20
	./gnark arithmetic --input="none" --operation="add" --field="bls12377" --count=20
	./gnark arithmetic --input="none" --operation="sub" --field="bls12377" --count=20
	./gnark arithmetic --input="none" --operation="mul" --field="bls12377" --count=20
	./gnark arithmetic --input="none" --operation="div" --field="bls12377" --count=20
	./gnark arithmetic --input="none" --operation="exp" --field="bls12377" --count=20
	./gnark arithmetic --input="none" --operation="add" --field="bls24315" --count=20
	./gnark arithmetic --input="none" --operation="sub" --field="bls24315" --count=20
	./gnark arithmetic --input="none" --operation="mul" --field="bls24315" --count=20
	./gnark arithmetic --input="none" --operation="div" --field="bls24315" --count=20
	./gnark arithmetic --input="none" --operation="exp" --field="bls24315" --count=20
	./gnark arithmetic --input="none" --operation="add" --field="bw6633" --count=20
	./gnark arithmetic --input="none" --operation="sub" --field="bw6633" --count=20
	./gnark arithmetic --input="none" --operation="mul" --field="bw6633" --count=20
	./gnark arithmetic --input="none" --operation="div" --field="bw6633" --count=20
	./gnark arithmetic --input="none" --operation="exp" --field="bw6633" --count=20
	./gnark arithmetic --input="none" --operation="add" --field="bw6761" --count=20
	./gnark arithmetic --input="none" --operation="sub" --field="bw6761" --count=20
	./gnark arithmetic --input="none" --operation="mul" --field="bw6761" --count=20
	./gnark arithmetic --input="none" --operation="div" --field="bw6761" --count=20
	./gnark arithmetic --input="none" --operation="exp" --field="bw6761" --count=20