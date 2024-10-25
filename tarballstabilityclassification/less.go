package tarballstabilityclassification

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["less"] = &ClassifierLess{}
}

type ClassifierLess struct {
}

func (self *ClassifierLess) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	// Less does not some logical version scheme to destinguish beta from stable
	// so less's site have to be checked periodically
	// http://www.greenwoodsoftware.com/less/
	for _, i := range []string{
		"487",
	} {
		if parsed.Version.StrSlice()[0] == i {
			return types.Release, nil
		}
	}

	return types.Beta, nil
}

func (self *ClassifierLess) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
