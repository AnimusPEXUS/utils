package tarballstabilityclassification

import (
	"strings"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["lynx"] = &ClassifierLynx{}
}

type ClassifierLynx struct {
}

func (self *ClassifierLynx) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if len(parsed.Status.StrSlice()) == 0 {
		return types.Development, nil
	}

	lower := strings.ToLower(parsed.Status.StrSlice()[0])

	if strings.Contains(lower, "rel") {
		return types.Release, nil
	}

	if strings.Contains(lower, "pre") {
		return types.RC, nil
	}

	if strings.Contains(lower, "dev") {
		return types.Development, nil
	}

	return types.Development, nil
}

func (self *ClassifierLynx) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
