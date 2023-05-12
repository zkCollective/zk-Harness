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

type BenchDataArithmetic struct {
	Framework string
	Category  string
	Curve     string
	Field     string
	Operation string
	Input     string
	MaxRAM    uint64
	Count     int
	RunTime   int64
}

func (bDataArith BenchDataArithmetic) Headers() []string {
	return []string{"framework", "category", "curve", "field", "operation", "input", "ram", "time(ns)", "nbPhysicalCores", "nbLogicalCores", "count", "cpu"}
}

func (bDataArith BenchDataArithmetic) Values() []string {
	return []string{
		bDataArith.Framework,
		bDataArith.Category,
		bDataArith.Curve,
		bDataArith.Field,
		bDataArith.Operation,
		bDataArith.Input,
		strconv.Itoa(int(bDataArith.MaxRAM)),
		strconv.Itoa(int(bDataArith.RunTime)),
		strconv.Itoa(CPU.PhysicalCores),
		strconv.Itoa(CPU.LogicalCores),
		strconv.Itoa(int(bDataArith.Count)),
		CPU.BrandName,
	}
}

type BenchDataCurve struct {
	Framework string
	Category  string
	Curve     string
	Operation string
	Input     string
	MaxRAM    uint64
	Count     int
	RunTime   int64
}

func (bDataCurve BenchDataCurve) Headers() []string {
	return []string{"framework", "category", "curve", "operation", "input", "ram", "time(ms)", "nbPhysicalCores", "nbLogicalCores", "count", "cpu"}
}

func (bDataCurve BenchDataCurve) Values() []string {
	return []string{
		bDataCurve.Framework,
		bDataCurve.Category,
		bDataCurve.Curve,
		bDataCurve.Operation,
		bDataCurve.Input,
		strconv.Itoa(int(bDataCurve.MaxRAM)),
		strconv.Itoa(int(bDataCurve.RunTime)),
		strconv.Itoa(CPU.PhysicalCores),
		strconv.Itoa(CPU.LogicalCores),
		strconv.Itoa(int(bDataCurve.Count)),
		CPU.BrandName,
	}
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
	return []string{"framework", "category", "backend", "curve", "circuit", "input", "operation", "nbConstraints", "nbSecret", "nbPublic", "ram", "time(ms)", "proofSize", "nbPhysicalCores", "nbLogicalCores", "count", "cpu"}
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
		strconv.Itoa(CPU.PhysicalCores),
		strconv.Itoa(CPU.LogicalCores),
		strconv.Itoa(int(bDataCirc.Count)),
		CPU.BrandName,
	}
}
