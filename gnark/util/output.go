package util

import (
	"strconv"

	. "github.com/klauspost/cpuid/v2"
)

type BenchData struct {
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
	RunTime           int64

	// Currently removed to fit log structure:
	// NbInternalVariables int
	// NbCoefficients      int
	// Throughput          int

	// CPU             info
	// NbPhysicalCores int
	// ThreadsPerCore  int
	// LogicalCores    int
	// CachelineBytes  int
	// L1DataBytes     int
	// L1InstrBytes    int
	// L2DataBytes     int
	// L3DataBytes     int
	// Frequency       int
	// SupportsADX     int
	// CPUName         string
}

func (bData BenchData) Headers() []string {
	return []string{"framework", "category", "backend", "curve", "circuit", "input", "operation", "nbConstraints", "nbSecret", "nbPublic", "ram(mb)", "time(ms)", "nbPhysicalCores", "nbLogicalCores", "cpu"}
}

func (bData BenchData) Values() []string {

	return []string{
		bData.Framework,
		bData.Category,
		bData.Backend,
		bData.Curve,
		bData.Circuit,
		bData.Input,
		bData.Operation,
		strconv.Itoa(int(bData.NbConstraints)),
		strconv.Itoa(int(bData.NbSecretVariables)),
		strconv.Itoa(int(bData.NbPublicVariables)),
		strconv.Itoa(int(bData.MaxRAM)),
		strconv.Itoa(int(bData.RunTime)),
		strconv.Itoa(CPU.PhysicalCores),
		strconv.Itoa(CPU.LogicalCores),
		CPU.BrandName,

		// strconv.Itoa(int(bData.NbInternalVariables)),
		// strconv.Itoa(bData.NbCoefficients),
		// strconv.Itoa(bData.Throughput),
		// strconv.Itoa(bData.Throughput / CPU.LogicalCores),

		// strconv.Itoa(CPU.ThreadsPerCore),
		// strconv.Itoa(CPU.CacheLine),
		// strconv.Itoa(CPU.Cache.L1D),
		// strconv.Itoa(CPU.Cache.L1I),
		// strconv.Itoa(CPU.Cache.L2),
		// strconv.Itoa(CPU.Cache.L3),
		// strconv.Itoa(int(CPU.Hz / 1000000)),
		// fmt.Sprintf("%v", CPU.Supports(ADX) && CPU.Supports(BMI2)),
		// fmt.Sprintf("%v", amd64_adx),
	}
}
