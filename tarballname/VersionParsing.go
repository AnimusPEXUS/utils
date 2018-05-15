package tarballname

import (
	"strings"

	"github.com/AnimusPEXUS/utils/versionorstatus"
)

func defaultVersionSplitterSub0(parsed_version *versionorstatus.ParsedVersion) {

	parsed_version.Arr = RemoveItemsList(parsed_version.DirtyArr, ALL_DELIMITERS)
	parsed_version.Str = strings.Join(parsed_version.Arr, ".")
	parsed_version.DirtyStr = strings.Join(parsed_version.DirtyArr, "")

}

func DefaultVersionSplitter(
	name_sliced SlicedName,
	most_possible_version Slice,
) *versionorstatus.ParsedVersion {
	var (
		ret *versionorstatus.ParsedVersion
	)

	ret = new(versionorstatus.ParsedVersion)

	ret.DirtyArr = append(ret.DirtyArr[:0], ret.DirtyArr[:0]...)

	for _, j := range name_sliced[most_possible_version[0]:most_possible_version[1]] {
		ret.DirtyArr = append(ret.DirtyArr, j)
	}

	defaultVersionSplitterSub0(ret)

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
) *versionorstatus.ParsedVersion {

	ret := new(versionorstatus.ParsedVersion)

	name_sliced1 := name_sliced[1]

	splitted_version := []string{name_sliced1[:1], name_sliced1[1:]}

	ret.DirtyArr = append(ret.DirtyArr, splitted_version...)

	defaultVersionSplitterSub0(ret)

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
