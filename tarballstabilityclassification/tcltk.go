package tarballstabilityclassification

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["tcltk"] = &ClassifierTclTk{}
}

type ClassifierTclTk struct {
}

func (self *ClassifierTclTk) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if parsed.Status.String() != "src" {
		return types.Development, nil
	}

	return types.Release, nil
}

func (self *ClassifierTclTk) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
