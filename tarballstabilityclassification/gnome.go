package tarballstabilityclassification

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["gnome"] = &ClassifierGnome{}
}

type ClassifierGnome struct {
}

func (self *ClassifierGnome) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if parsed.Status.Str != "" {
		return types.Development, nil
	}

	version, err := parsed.Version.ArrInt()
	if err != nil {
		return types.Development, err
	}

	if len(version) < 2 {
		return types.Development,
			errors.New("version numbers array too short")
	}

	if version[1]%2 != 0 {
		return types.Development, nil
	}

	return types.Release, nil
}

func (self *ClassifierGnome) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
