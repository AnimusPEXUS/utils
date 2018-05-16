package tarballstabilityclassification

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["strict"] = &ClassifierStrict{}
}

type ClassifierStrict struct {
}

func (self *ClassifierStrict) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	switch parsed.Status.DirtyString() {
	default:
		return types.Development, nil
	case "":
		fallthrough
	case "src":
		fallthrough
	case "source":
		fallthrough
	case "release":
		return types.Release, nil
	}

	return types.Development, errors.New("impossible error")
}

func (self *ClassifierStrict) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
