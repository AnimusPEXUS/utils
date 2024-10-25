package tarballstabilityclassification

import (
	"errors"
	"strings"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["libevent"] = &ClassifierLibEvent{}
}

type ClassifierLibEvent struct {
}

func (self *ClassifierLibEvent) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if len(parsed.Status.StrSlice()) == 0 {
		return types.Development, nil
	}

	switch strings.ToLower(parsed.Status.StrSlice()[0]) {
	default:
		return types.Development, nil
	case "stable":
		return types.Release, nil
	case "rc":
		return types.RC, nil
	case "beta":
		return types.Beta, nil
	}

	return types.Development, errors.New("unexpected return")
}

func (self *ClassifierLibEvent) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
