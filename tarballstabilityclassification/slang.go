package tarballstabilityclassification

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["slang"] = &ClassifierSlang{}
}

type ClassifierSlang struct {
}

func (self *ClassifierSlang) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	return types.Release, nil
}

func (self *ClassifierSlang) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
