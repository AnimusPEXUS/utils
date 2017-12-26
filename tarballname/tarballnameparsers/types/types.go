package types

import (
	"github.com/AnimusPEXUS/utils/tarballname"
)

type DevelopmentStatus uint

type TarballNameParserI interface {
	ParseName(value string) (*ParseResult, error)
}

const (
	Development DevelopmentStatus = iota
	PreAlpha
	Alpha
	Beta
	RC
	RTM
	GA
	Gold
)

func (self DevelopmentStatus) String() string {
	switch self {
	default:
		return "unknown"
	case Development:
		return "development"
	case PreAlpha:
		return "prealpha"
	case Alpha:
		return "alpha"
	case Beta:
		return "beta"
	case RC:
		return "RC"
	case RTM:
		return "RTM"
	case GA:
		return "GA"
	case Gold:
		return "Gold"
	}
}

type ParseResult struct {
	Name        string
	HaveVersion bool
	Version     tarballname.ParsedVersion
	Status      tarballname.ParsedStatus
	HaveStatus  bool
	// DevelopmentStatus DevelopmentStatus // TODO: make function
	HaveBuildId bool
	BuildId     string
}
