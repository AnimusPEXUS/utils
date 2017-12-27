package tarballnamefilterfunctions

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/AnimusPEXUS/aipsetup/basictypes"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers"
	"github.com/AnimusPEXUS/utils/textlist"
)

var TarballNameFilterFunctions textlist.FilterFunctions

func init() {

	TarballNameFilterFunctions = make(textlist.FilterFunctions)

	// for k, v := range versionfilterfunctions.VersionFilterFunctions {
	// 	TarballNameFilterFunctions[k] = v
	// }

	TarballNameFilterFunctions["tarball-status-re"] = StatusFilterFunction

}

func StatusFilterFunction(
	parameter string,
	case_sensitive bool,
	value_to_match string,
	data map[string]interface{},
) (bool, error) {

	fmt.Println("StatusFilterFunction", parameter, case_sensitive, value_to_match)

	info, ok := data["pkg_info"].(*basictypes.PackageInfo)
	if !ok {
		return false,
			errors.New(
				"StatusFilterFunction requires data[\"pkg_info\"] of type *basictypes.PackageInfo",
			)
	}

	// TODO: I don't like this. If somebody knows how to get those functions here
	// without (interface{})s and import loop problems - let Me know.
	// DetermineTarballsBuildInfo :=
	// 	data["DetermineTarballsBuildInfo"].(func(string) (map[string]*basictypes.PackageInfo, error))
	//
	// var info *basictypes.PackageInfo
	//
	// {
	// 	res, err := DetermineTarballsBuildInfo(value_to_match)
	// 	if err != nil {
	// 		return false, err
	// 	}
	//
	// 	if len(res) != 1 {
	// 		return false, errors.New(
	// 			"can't correctly determine info for given tarballnmame",
	// 		)
	// 	}
	//
	// 	for _, info = range res {
	// 	}
	// }

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

	fmt.Println("status in question", parse_res.Status.Str)

	return rexp.MatchString(parse_res.Status.Str), nil
}
