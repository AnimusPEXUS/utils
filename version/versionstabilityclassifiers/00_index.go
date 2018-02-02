package versionstabilityclassifiers

import (
	"errors"

	"github.com/AnimusPEXUS/utils/version/versionstabilityclassifiers/types"
)

func Get(name string) (types.VersionStabilityClassifierI, error) {
	if t, ok := Index[name]; ok {
		return t, nil
	} else {
		return nil, errors.New("classifier not found")
	}
}

var Index = map[string](types.VersionStabilityClassifierI){}
