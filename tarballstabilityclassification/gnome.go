package tarballstabilityclassification

import (
	"errors"
	"regexp"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

var TESTING_9X_RE = regexp.MustCompile(`^[98]\d+$`)

func init() {
	Index["gnome"] = &ClassifierGnome{}
}

type ClassifierGnome struct {
}

func (self *ClassifierGnome) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	if parsed.Status.String() != "" {
		return types.Development, nil
	}

	version, err := parsed.Version.IntSlice()
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

	for _, i := range parsed.Version.StrSlice() {
		if TESTING_9X_RE.MatchString(i) {
			return types.Development, err
		}
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
