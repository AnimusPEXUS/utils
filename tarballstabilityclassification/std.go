package tarballstabilityclassification

import (
	"strings"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballstabilityclassification/types"
)

func init() {
	Index["std"] = &ClassifierStd{}
}

type ClassifierStd struct {
}

// NOTE: as aipsetup now enforces use of stability check modules upon
//       updating tarballs repository, there must be a module, which allways
//       returns 'Release' value. 'std' MUST be this module!

func (self *ClassifierStd) Check(parsed *tarballname.ParsedTarballName) (
	types.StabilityClassification,
	error,
) {

	// switch parsed.Status.DirtyStr {
	// default:
	// 	return types.Development, nil
	// case "":
	// 	fallthrough
	// case "src":
	// 	fallthrough
	// case "source":
	// 	fallthrough
	// case "release":
	// 	return types.Release, nil
	// }

	switch parsed.Status.DirtyString() {
	case "alpha":
		return types.Alpha, nil
	case "beta":
		return types.Beta, nil
	case "dev":
		return types.Development, nil
	}

	if strings.HasPrefix(parsed.Status.DirtyString(), "rc") {
		return types.RC, nil
	}

	return types.Release, nil
}

func (self *ClassifierStd) IsStable(parsed *tarballname.ParsedTarballName) (bool, error) {
	cr, err := self.Check(parsed)
	if err != nil {
		return false, err
	}
	return types.IsStable(cr), nil
}
