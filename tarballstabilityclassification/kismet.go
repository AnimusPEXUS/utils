package tarballstabilityclassification

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["kismet"] = &ClassifierKismet{}
}

type ClassifierKismet struct {
}

func (self *ClassifierKismet) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if len(parsed.Status.Arr) == 0 {
		return types.Development, nil
	}

	switch parsed.Status.Arr[0] {
	default:
		return types.Development, nil
	case "R":
		return types.Release, nil
	case "RC":
		return types.RC, nil
	}

	return types.Development, errors.New("burn in hell!")
}

func (self *ClassifierKismet) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
