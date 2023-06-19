package util

import (
	"strconv"
)

type HeadersProvider interface {
	Headers() []string
}

type ValuesProvider interface {
	Values() []string
}

type BenchDataCircuit struct {
	Framework         string
	Category          string
	Backend           string
	Curve             string
	Circuit           string
	Input             string
	Operation         string
	NbConstraints     int
	NbSecretVariables int
	NbPublicVariables int
	MaxRAM            uint64
	Count             int
	RunTime           int64
	ProofSize         int
}

func (bDataCirc BenchDataCircuit) Headers() []string {
	return []string{"framework", "category", "backend", "curve", "circuit", "input", "operation", "nbConstraints", "nbSecret", "nbPublic", "ram", "time(ms)", "proofSize", "count"}
}

func (bDataCirc BenchDataCircuit) Values() []string {
	return []string{
		bDataCirc.Framework,
		bDataCirc.Category,
		bDataCirc.Backend,
		bDataCirc.Curve,
		bDataCirc.Circuit,
		bDataCirc.Input,
		bDataCirc.Operation,
		strconv.Itoa(int(bDataCirc.NbConstraints)),
		strconv.Itoa(int(bDataCirc.NbSecretVariables)),
		strconv.Itoa(int(bDataCirc.NbPublicVariables)),
		strconv.Itoa(int(bDataCirc.MaxRAM)),
		strconv.Itoa(int(bDataCirc.RunTime)),
		strconv.Itoa(int(bDataCirc.ProofSize)),
		strconv.Itoa(int(bDataCirc.Count)),
	}
}

type BenchDataRecursion struct {
	Framework          string
	Category           string
	InnerBackend       string
	OuterBackend       string
	InnerCurve         string
	OuterCurve         string
	Circuit            string
	Input              string
	Operation          string
	InnerNbConstraints int
	NbConstraints      int
	NbSecretVariables  int
	NbPublicVariables  int
	MaxRAM             uint64
	RunTime            int64
	ProofSize          int
	Count              int
}

func (bDataCirc BenchDataRecursion) Headers() []string {
	return []string{"framework", "category", "innerBackend", "outerBackend", "innerCurve", "outerCurve", "circuit", "input", "operation", "innerNbConstraints", "outerNbConstraints", "nbSecret", "nbPublic", "ram", "time(ms)", "proofSize", "count"}
}

func (bDataCirc BenchDataRecursion) Values() []string {
	return []string{
		bDataCirc.Framework,
		bDataCirc.Category,
		bDataCirc.InnerBackend,
		bDataCirc.OuterBackend,
		bDataCirc.InnerCurve,
		bDataCirc.OuterCurve,
		bDataCirc.Circuit,
		bDataCirc.Input,
		bDataCirc.Operation,
		strconv.Itoa(int(bDataCirc.InnerNbConstraints)),
		strconv.Itoa(int(bDataCirc.NbConstraints)),
		strconv.Itoa(int(bDataCirc.NbSecretVariables)),
		strconv.Itoa(int(bDataCirc.NbPublicVariables)),
		strconv.Itoa(int(bDataCirc.MaxRAM)),
		strconv.Itoa(int(bDataCirc.RunTime)),
		strconv.Itoa(int(bDataCirc.ProofSize)),
		strconv.Itoa(int(bDataCirc.Count)),
	}
}
