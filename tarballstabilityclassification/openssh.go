package tarballstabilityclassification

import (
	"regexp"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["openssh"] = &ClassifierOpenSSH{}
}

type ClassifierOpenSSH struct {
}

func (self *ClassifierOpenSSH) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if ok, err := regexp.MatchString(`^[Pp]\d+$`, parsed.Status.DirtyStr); err != nil {
		return types.Development, err
	} else {
		if !ok {
			return types.Development, nil
		}
	}

	return types.Release, nil
}

func (self *ClassifierOpenSSH) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
