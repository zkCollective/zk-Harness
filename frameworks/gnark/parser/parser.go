package parser

import (
	"errors"
	"strings"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/pkg/profile"
	"github.com/zkCollective/zk-Harness/frameworks/gnark/circuits"
)

type Config struct {
	Circuit      *string
	CircuitSize  *int
	Algo         *string
	Profile      *string
	Count        *int
	Curve        *string
	InputPath    *string
	Operation    *string
	OuterBackend *string
	OutputPath   *string
}

func NewConfig() *Config {
	return &Config{
		Circuit:      new(string),
		CircuitSize:  new(int),
		Algo:         new(string),
		Profile:      new(string),
		Count:        new(int),
		Curve:        new(string),
		InputPath:    new(string),
		Operation:    new(string),
		OuterBackend: new(string),
		OutputPath:   new(string),
	}
}

var (
	InnerCurveID ecc.ID
	CurveID      ecc.ID
	P            func(p *profile.Profile)
	C            circuits.BenchCircuit
)

var config Config

func ParseFlags(config *Config) error {
	if *config.CircuitSize <= 0 {
		return errors.New("circuit size must be >= 0")
	}
	if *config.Count <= 0 {
		return errors.New("bench count must be >= 0")
	}

	switch *config.Algo {
	case "compile", "setup", "witness", "prove", "verify":
	default:
		return errors.New("invalid algo")
	}

	switch *config.Profile {
	case "none":
	case "trace":
		P = profile.TraceProfile
	case "cpu":
		P = profile.CPUProfile
	case "mem":
		P = profile.MemProfile
	default:
		return errors.New("invalid profile")
	}

	curves := ecc.Implemented()
	for _, id := range curves {
		if *config.Curve == strings.ToLower(id.String()) {
			CurveID = id
		}
	}
	if CurveID == ecc.UNKNOWN {
		return errors.New("invalid curve")
	}

	var ok bool
	C, ok = circuits.BenchCircuits[*config.Circuit]
	if !ok {
		return errors.New("unknown circuit")
	}

	return nil
}

func ParseFlagsMemory(config *Config) error {

	if *config.CircuitSize <= 0 {
		return errors.New("circuit size must be >= 0")
	}
	if *config.Count <= 0 {
		return errors.New("bench count must be >= 0")
	}

	curves := ecc.Implemented()
	for _, id := range curves {
		if *config.Curve == strings.ToLower(id.String()) {
			CurveID = id
		}
	}
	if CurveID == ecc.UNKNOWN {
		return errors.New("invalid curve")
	}

	var ok bool
	C, ok = circuits.BenchCircuits[*config.Circuit]
	if !ok {
		return errors.New("unknown circuit")
	}

	return nil
}

func AssertNoError(err error) {
	if err != nil {
		panic(err)
	}
}
