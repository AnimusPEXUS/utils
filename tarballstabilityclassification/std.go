package tarballstabilityclassification

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["std"] = &ClassifierStd{}
}

type ClassifierStd struct {
}

func (self *ClassifierStd) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if parsed.Status.Str != "" {
		return types.Development, nil
	}

	return types.Release, nil
}

func (self *ClassifierStd) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
