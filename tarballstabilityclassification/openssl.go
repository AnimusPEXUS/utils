package tarballstabilityclassification

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["openssl"] = &ClassifierOpenSSL{}
}

type ClassifierOpenSSL struct {
}

func (self *ClassifierOpenSSL) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	return types.Release, nil
}

func (self *ClassifierOpenSSL) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
