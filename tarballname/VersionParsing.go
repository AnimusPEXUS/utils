package tarballname

import (
	"github.com/AnimusPEXUS/utils/versionorstatus"
)

type (
	VersionFinderFunction func(name_sliced SlicedName) (*Slice, bool)

	VersionSplitterFunction func(
		name_sliced SlicedName,
		most_possible_version Slice,
	) *versionorstatus.ParsedVersionOrStatus

	VersionsSelectorFunction func(tarballbasename string) (
		VersionFinderFunction,
		VersionSplitterFunction,
	)
)

// func defaultVersionSplitterSub0(parsed_version *versionorstatus.ParsedVersionOrStatus) {
//
// 	parsed_version.Arr = RemoveItemsList(parsed_version.DirtyArr, ALL_DELIMITERS)
// 	parsed_version.Str = strings.Join(parsed_version.Arr, ".")
// 	parsed_version.DirtyStr = strings.Join(parsed_version.DirtyArr, "")
//
// }

func DefaultVersionSplitter(
	name_sliced SlicedName,
	most_possible_version Slice,
) *versionorstatus.ParsedVersionOrStatus {

	values := make([]string, 0)

	for _, j := range name_sliced[most_possible_version[0]:most_possible_version[1]] {
		values = append(values, j)
	}

	values = RemoveItemsList(values, ALL_DELIMITERS)

	sep := "."
	if len(values) > 1 {
		sep = values[1]
	}

	ret := versionorstatus.NewParsedVersionOrStatusFromStringSlice(values, sep)

	return ret
}

func InfoZipVersionFinder(
	name_sliced SlicedName,
) (
	*Slice,
	bool,
) {
	return &Slice{1, 2}, true
}

func InfoZipVersionSplitter(
	name_sliced SlicedName,
	most_possible_version Slice,
) *versionorstatus.ParsedVersionOrStatus {

	name_sliced1 := name_sliced[1]

	values := []string{name_sliced1[:1], name_sliced1[1:]}

	ret := versionorstatus.NewParsedVersionOrStatusFromStringSlice(values, "")

	return ret
}

func StrictVersionsFunctionSelector(
	tarballbasename string,
) (VersionFinderFunction, VersionSplitterFunction) {
	return DefaultVersionFinder, DefaultVersionSplitter
}

func InfoZipVersionsFunctionSelector(
	tarballbasename string,
) (VersionFinderFunction, VersionSplitterFunction) {
	return InfoZipVersionFinder, InfoZipVersionSplitter
}
