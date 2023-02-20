package util

import (
	"strconv"

	. "github.com/klauspost/cpuid/v2"
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
	ProofSize         int
	MaxRAM            uint64
	RunTime           int64
}

func (bDataCirc BenchDataCircuit) Headers() []string {
	return []string{"framework", "category", "backend", "curve", "circuit", "input", "operation", "nbConstraints", "nbSecret", "nbPublic", "proofSize", "ram(mb)", "time(ms)", "nbPhysicalCores", "nbLogicalCores", "cpu"}
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
		strconv.Itoa(int(bDataCirc.ProofSize)),
		strconv.Itoa(int(bDataCirc.MaxRAM)),
		strconv.Itoa(int(bDataCirc.RunTime)),
		strconv.Itoa(CPU.PhysicalCores),
		strconv.Itoa(CPU.LogicalCores),
		CPU.BrandName,
	}
}

type BenchDataArithmetic struct {
	Framework string
	Category  string
	Field     string // native / non-native
	Order     int
	Operation string
	Input     string
	MaxRAM    uint64
	RunTime   int64
}

func (bDataArith BenchDataArithmetic) Headers() []string {
	return []string{"framework", "category", "field", "p(bitlength)", "operation", "input", "ram(mb)", "time(ns)", "nbPhysicalCores", "nbLogicalCores", "cpu"}
}

func (bDataArith BenchDataArithmetic) Values() []string {
	return []string{
		bDataArith.Framework,
		bDataArith.Category,
		bDataArith.Field,
		strconv.Itoa(int(bDataArith.Order)),
		bDataArith.Operation,
		bDataArith.Input,
		strconv.Itoa(int(bDataArith.MaxRAM)),
		strconv.Itoa(int(bDataArith.RunTime)),
		strconv.Itoa(CPU.PhysicalCores),
		strconv.Itoa(CPU.LogicalCores),
		CPU.BrandName,
	}
}
