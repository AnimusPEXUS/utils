package types

import (
	"github.com/AnimusPEXUS/utils/tarballname"
)

type VersionStabilityClassifierI interface {
	Check(parsed *tarballname.ParsedTarballName) (StabilityClassification, error)
	IsStable(parsed *tarballname.ParsedTarballName) (bool, error)
}
