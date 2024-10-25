package tarballstabilityclassification

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["gcc"] = &ClassifierGCC{}
}

type ClassifierGCC struct {
}

func (self *ClassifierGCC) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if parsed.Status.StrSliceString("") != "" {
		return types.Development, nil
	}

	version, err := parsed.Version.IntSlice()
	if err != nil {
		return types.Development, err
	}

	if version[0] < 5 {
		return types.Release, nil
	}

	if version[1] == 0 {
		if version[2] == 0 {
			return types.Alpha, nil
		} else {
			return types.Beta, nil
		}
	}

	return types.Release, nil
}

func (self *ClassifierGCC) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
