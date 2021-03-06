package tarballnamefilterfunctions

import (
	"errors"
	"path"
	"regexp"

	"github.com/AnimusPEXUS/utils/softwareunit"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers"
	"github.com/AnimusPEXUS/utils/textlist"
)

var TarballNameFilterFunctions textlist.FilterFunctions

func init() {

	TarballNameFilterFunctions = make(textlist.FilterFunctions)

	TarballNameFilterFunctions["basename-re"] = BasenameReFilterFunction
	TarballNameFilterFunctions["status-re"] = StatusFilterFunction

}

func BasenameReFilterFunction(
	parameter string,
	case_sensitive bool,
	value_to_match string,
	data map[string]interface{},
) (bool, error) {

	value_to_match = path.Base(value_to_match)
	rexp, err := regexp.Compile(parameter)
	if err != nil {
		return false, err
	}

	return rexp.MatchString(value_to_match), nil
}

func StatusFilterFunction(
	parameter string,
	case_sensitive bool,
	value_to_match string,
	data map[string]interface{},
) (bool, error) {

	info, ok := data["pkg_info"].(*softwareunit.PackageInfo)
	if !ok {
		return false,
			errors.New(
				"StatusFilterFunction requires data[\"pkg_info\"] of type *softwareunit.PackageInfo",
			)
	}

	parser, err := tarballnameparsers.Get(info.TarballFileNameParser)
	if err != nil {
		return false, err
	}

	parse_res, err := parser.Parse(value_to_match)
	if err != nil {
		return false, err
	}

	rexp, err := regexp.Compile(parameter)
	if err != nil {
		return false, err
	}

	return rexp.MatchString(parse_res.Status.StrSliceString("")), nil
}
