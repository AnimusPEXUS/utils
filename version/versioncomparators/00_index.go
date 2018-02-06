package versioncomparators

import (
	"errors"

	"github.com/AnimusPEXUS/utils/version/versioncomparators/types"
)

func Get(name string) (types.VersionComparatorI, error) {
	if t, ok := Index[name]; ok {
		return t, nil
	} else {
		return nil, errors.New("version comparator not found")
	}
}

var Index = map[string]types.VersionComparatorI{}
