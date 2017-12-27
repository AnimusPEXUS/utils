package versionfilterfunctions

import (
	"errors"
	"strconv"
	"strings"

	"github.com/AnimusPEXUS/aipsetup/basictypes"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers"
	"github.com/AnimusPEXUS/utils/textlist"
	"github.com/AnimusPEXUS/utils/version"
)

var VersionFilterFunctions textlist.FilterFunctions

func VersionCheck(
	function string,
	parameter string,
	value_to_match string,
	data map[string]interface{},
) (bool, error) {

	info, ok := data["info"].(*basictypes.PackageInfo)
	if !ok {
		return false, errors.New("VersionCheck requires data[\"info\"]")
	}

	p, err := tarballnameparsers.Get(info.TarballFileNameParser)
	if err != nil {
		return false, err
	}

	res, err := p.Parse(value_to_match)
	if err != nil {
		return false, err
	}

	// fmt.Println("VersionCheck", function, parameter, value_to_match)

	ret := false

	vtm_i_array, err := res.Version.ArrInt()
	if err != nil {
		return false, err
	}

	param_i_array := make([]int, 0)
	for _, i := range strings.Split(parameter, ".") {
		i_i, err := strconv.Atoi(i)
		if err != nil {
			return false, err
		}
		param_i_array = append(param_i_array, i_i)
	}

	switch function {
	default:
		return false, errors.New("invalid version comparison function")
	case "<":
		ret = version.Compare(vtm_i_array, param_i_array) == -1
	case "<=":
		r := version.Compare(vtm_i_array, param_i_array)
		ret = r == -1 || r == 0
	case "==":
		ret = version.Compare(vtm_i_array, param_i_array) == 0
	case ">=":
		r := version.Compare(vtm_i_array, param_i_array)
		ret = r == 0 || r == 1
	case ">":
		ret = version.Compare(vtm_i_array, param_i_array) == 1
	case "!=":
		ret = version.Compare(vtm_i_array, param_i_array) != 0
	}

	return ret, nil
}

func init() {
	VersionFilterFunctions = make(textlist.FilterFunctions)

	// for k, v := range VersionFilterFunctions {
	// 	VersionFilterFunctions[k] = v
	// }

	VersionFilterFunctions["version-<"] = func(
		parameter string,
		case_sensitive bool,
		value_to_match string,
		data map[string]interface{},
	) (bool, error) {
		return VersionCheck("<", parameter, value_to_match, data)
	}
	VersionFilterFunctions["version-<="] = func(
		parameter string,
		case_sensitive bool,
		value_to_match string,
		data map[string]interface{},
	) (bool, error) {
		return VersionCheck("<=", parameter, value_to_match, data)
	}
	VersionFilterFunctions["version-=="] = func(
		parameter string,
		case_sensitive bool,
		value_to_match string,
		data map[string]interface{},
	) (bool, error) {
		return VersionCheck("==", parameter, value_to_match, data)
	}
	VersionFilterFunctions["version->="] = func(
		parameter string,
		case_sensitive bool,
		value_to_match string,
		data map[string]interface{},
	) (bool, error) {
		return VersionCheck(">=", parameter, value_to_match, data)
	}
	VersionFilterFunctions["version->"] = func(
		parameter string,
		case_sensitive bool,
		value_to_match string,
		data map[string]interface{},
	) (bool, error) {
		return VersionCheck(">", parameter, value_to_match, data)
	}
	VersionFilterFunctions["version-!="] = func(
		parameter string,
		case_sensitive bool,
		value_to_match string,
		data map[string]interface{},
	) (bool, error) {
		return VersionCheck("!=", parameter, value_to_match, data)
	}
}
