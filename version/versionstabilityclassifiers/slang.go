package versionstabilityclassifiers

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification"
)

func init() {
	Index["slang"] = &ClassifierSlang{}
}

type ClassifierSlang struct {
}

func (self *ClassifierSlang) Check(parsed *tarballname.ParsedTarballName) (
	tarballstabilityclassification.StabilityClassification,
	error,
) {

	return tarballstabilityclassification.Release, nil
}

func (self *ClassifierSlang) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return tarballstabilityclassification.IsStable(cr), nil
}
