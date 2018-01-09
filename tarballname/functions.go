package tarballname

import (
	"errors"
	"path"
	"regexp"
	"sort"
	"strings"
)

var (
	ACCEPTABLE_TARBALL_EXTENSIONS                          []string
	ACCEPTABLE_SOURCE_NAME_EXTENSIONS_REV_SORTED_BY_LENGTH []string
	KNOWN_SIGNING_EXTENSIONS                               []string
	ALL_DELIMITERS                                         []string
	STATUS_DELIMITERS                                      []string

	TARBALL_NAME_SEPARATION *regexp.Regexp = regexp.MustCompile(
		`[a-zA-Z]+|\d+|[\.\-\_\~\+]`,
	)
)

func init() {

	// NOTE: do not change order in this array: it's should be sorted in order
	//       of downloading preference
	ACCEPTABLE_TARBALL_EXTENSIONS = []string{
		".tar.xz",
		".tar.lzma",
		".tar.bz2",
		".tar.gz",
		".txz",
		".tlzma",
		".tbz2",
		".tbz",
		".tgz",
		".7z",
		".zip",
		".jar",
		".tar",
	}

	ACCEPTABLE_SOURCE_NAME_EXTENSIONS_REV_SORTED_BY_LENGTH =
		make([]string, len(ACCEPTABLE_TARBALL_EXTENSIONS))

	copy(
		ACCEPTABLE_SOURCE_NAME_EXTENSIONS_REV_SORTED_BY_LENGTH,
		ACCEPTABLE_TARBALL_EXTENSIONS,
	)

	sort.Sort(
		acceptable_tarball_ext_sorter(
			ACCEPTABLE_SOURCE_NAME_EXTENSIONS_REV_SORTED_BY_LENGTH,
		),
	)

	KNOWN_SIGNING_EXTENSIONS = []string{
		".sign", ".asc",
	}

	ALL_DELIMITERS = []string{".", "_", "-", "~"}

	STATUS_DELIMITERS = append(ALL_DELIMITERS, "+")
}

type acceptable_tarball_ext_sorter []string

func (a acceptable_tarball_ext_sorter) Len() int {
	return len(a)
}

func (a acceptable_tarball_ext_sorter) Swap(i, j int) {
	a[i], a[j] =
		a[j], a[i]
}

func (a acceptable_tarball_ext_sorter) Less(i, j int) bool {
	return len(a[i]) > len(a[j])
}

func IsPossibleTarballName(filename string) bool {
	for _, i := range ACCEPTABLE_SOURCE_NAME_EXTENSIONS_REV_SORTED_BY_LENGTH {
		if strings.HasSuffix(filename, i) {
			return true
		}
	}
	return false
}

func IsPossibleTarballNameErr(filename string) error {
	if !IsPossibleTarballName(filename) {
		return errors.New("file name doesn't looks like tarball's")
	}
	return nil
}

type (
	SlicedName []string

	Slice              [2]uint64
	SlicesSlice        []Slice
	SinglesOrMultiples SlicesSlice

	SinglesAndMultiples struct {
		Singles   SinglesOrMultiples
		Multiples SinglesOrMultiples
	}

	MapOfSinglesAndMultiples (map[string](*SinglesAndMultiples))
	MapOfSinglesOrMultiples  (map[string](*SinglesOrMultiples))

	VersionFinderFunction func(name_sliced SlicedName) (*Slice, bool)

	VersionSplitterFunction func(
		name_sliced SlicedName,
		most_possible_version Slice,
	) *ParsedVersion

	VersionsSelectorFunction func(tarballbasename string) (
		VersionFinderFunction,
		VersionSplitterFunction,
	)
)

// TODO: need to find better solution, or write separate module for
//       functions like this
func IsDecimal(val string) bool {

	matched, err := regexp.MatchString(`^\d+$`, val)

	return err == nil && matched
}

// TODO: Move this function to separate module
func StripList(data []string, what_to_remove []string) []string {

	var ret []string

	ret = append(ret[:0], data...)

	for _, i := range what_to_remove {

		for true {
			if len(ret) == 0 {
				break
			}
			if ret[0] == i {
				ret = ret[1:]
			} else {
				break
			}
		}

		for true {
			if len(ret) == 0 {
				break
			}
			if ret[len(ret)-1] == i {
				ret = ret[:len(ret)-1]
			} else {
				break
			}
		}
	}

	return ret
}

func RemoveItemsList(data []string, what_to_remove []string) []string {

	var ret []string

	ret = append(ret[:0], data...)

	len_ret := len(ret)

removal:
	for ii, i := range ret {
		ii = len_ret - (ii + 1)
		i = ret[ii]

		for _, j := range what_to_remove {
			if i == j {
				ret = append(
					ret[:ii],
					ret[ii+1:]...,
				)
				continue removal
			}
		}
	}

	return ret
}

func FindPossibleCharedSinglesAndMultiplesSub0(
	slices SlicesSlice,
	version_started int64,
	version_ended int64,
) SlicesSlice {
	return append(
		slices,
		[2]uint64{
			uint64(version_started),
			uint64(version_ended) + 1,
		},
	)
}

func FindPossibleCharedSinglesAndMultiples(
	name_sliced SlicedName,
	separator string,
) *SinglesAndMultiples {

	var (
		slices          SlicesSlice
		version_started int64 = -1
		version_ended   int64 = -1
		index           int64 = -1
		ret             *SinglesAndMultiples
	)

	ret = new(SinglesAndMultiples)

	if len(([]rune)(separator)) != 1 {
		panic("`separator' text lenght must be exactly 1")
	}

	for _, i := range name_sliced {

		index++

		if IsDecimal(i) {
			if version_started == -1 {
				version_started = index
			}

			version_ended = index
		} else {
			if version_started != -1 {
				if i != separator {
					slices = FindPossibleCharedSinglesAndMultiplesSub0(
						slices,
						version_started,
						version_ended,
					)

					version_started = -1
				}
			}
		}
	}

	if version_started != -1 {
		slices = FindPossibleCharedSinglesAndMultiplesSub0(
			slices,
			version_started,
			version_ended,
		)
	}

	for _, i := range slices {
		if i[1]-i[0] == 1 {
			ret.Singles = append(ret.Singles, i)
		} else if i[1]-i[0] > 1 {
			ret.Multiples = append(ret.Multiples, i)
		} else {
			panic("programming error")
		}
	}

	return ret
}

func FindAllVersionsAndSingles(
	name_sliced SlicedName,
) MapOfSinglesAndMultiples {

	ret := make(MapOfSinglesAndMultiples)

	for _, i := range ALL_DELIMITERS {
		ret[i] = FindPossibleCharedSinglesAndMultiples(name_sliced, i)
	}

	return ret
}

func defaultVersionFinderSub0(
	v1 MapOfSinglesOrMultiples,
) (
	*Slice,
	bool,
) {

	var (
		found bool = false
		ret   *Slice
	)

	ret = new(Slice)

search0:
	for _, i := range ALL_DELIMITERS {

		switch len(*v1[i]) {
		case 0:
			continue
		case 1:
			(*ret) = (*v1[i])[0]
			found = true
			break search0
		default:
			current_delimiter_group := v1[i]

			maximum_length := uint64(0)

			for _, j := range *current_delimiter_group {
				l := j[1] - j[0]
				if l > maximum_length {
					maximum_length = l
				}
			}

			if maximum_length != 0 {
				lists_to_compare := make(SinglesOrMultiples, 0)

				for _, j := range *current_delimiter_group {
					l := j[1] - j[0]
					if l == maximum_length {
						lists_to_compare = append(lists_to_compare, j)
					}
				}

				switch len(lists_to_compare) {
				case 0:
					//pass
				case 1:
					ret = &lists_to_compare[0]
					found = true
					break search0
				default:
					ret = &lists_to_compare[0]

					for _, j := range lists_to_compare {
						if j[0] < ret[0] {
							ret = &j
						}
					}

					found = true
					break search0
				}
			}
		}
	}
	if !found {
		ret = nil
	}
	return ret, found
}

func DefaultVersionFinder(
	name_sliced SlicedName,
) (
	*Slice,
	bool,
) {
	var (
		found bool = false

		possible_versions_and_singles_grouped  MapOfSinglesAndMultiples
		possible_versions_grouped_by_delimeter MapOfSinglesOrMultiples
		possible_singles_grouped_by_delimeter  MapOfSinglesOrMultiples

		ret *Slice
	)

	ret = new(Slice)

	possible_versions_and_singles_grouped = FindAllVersionsAndSingles(
		name_sliced,
	)

	possible_versions_grouped_by_delimeter = make(MapOfSinglesOrMultiples)
	possible_singles_grouped_by_delimeter = make(MapOfSinglesOrMultiples)

	for _, i := range ALL_DELIMITERS {
		possible_versions_grouped_by_delimeter[i] =
			&(possible_versions_and_singles_grouped[i].Multiples)

		possible_singles_grouped_by_delimeter[i] =
			&(possible_versions_and_singles_grouped[i].Singles)

	}

	ret, found = defaultVersionFinderSub0(
		possible_versions_grouped_by_delimeter,
	)
	if !found {
		ret, found = defaultVersionFinderSub0(
			possible_singles_grouped_by_delimeter,
		)
	}

	if !found {
		ret = nil
	}

	return ret, found
}

func defaultVersionSplitterSub0(parsed_version *ParsedVersion) {

	parsed_version.Arr =
		RemoveItemsList(parsed_version.DirtyArr, ALL_DELIMITERS)

	parsed_version.Str = strings.Join(parsed_version.Arr, ".")
	parsed_version.DirtyStr = strings.Join(parsed_version.DirtyArr, "")

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
) *ParsedVersion {
	var (
		ret *ParsedVersion
	)

	ret = new(ParsedVersion)

	ret.DirtyArr = append(
		ret.DirtyArr[:0],
		strings.Split(name_sliced[1], "")...,
	)

	defaultVersionSplitterSub0(ret)

	return ret
}

// func DefaultVersionsFunctionSelector(
// 	tarballbasename string,
// ) (
// 	VersionFinderFunction,
// 	VersionSplitterFunction,
// ) {
// 	var (
// 		version_finder_function   VersionFinderFunction   = DefaultVersionFinder
// 		version_splitter_function VersionSplitterFunction = DefaultVersionSplitter
// 	)
//
// 	res, err := regexp.MatchString(`^(un)?zip\d+.*$`, tarballbasename)
//
// 	if err == nil && res == true {
// 		version_finder_function = InfoZipVersionFinder
// 		version_splitter_function = InfoZipVersionSplitter
// 	}
//
// 	return version_finder_function, version_splitter_function
//
// }

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

func ParseEx(
	full_path_or_basename string,
	acceptable_extensions []string,
	versions_selector_function VersionsSelectorFunction,
) (*ParsedTarballName, error) {
	var (
		ret *ParsedTarballName

		version_finder_function   VersionFinderFunction   = DefaultVersionFinder
		version_splitter_function VersionSplitterFunction = DefaultVersionSplitter

		most_possible_version *Slice
		found                 bool = false

		version_splitted *ParsedVersion

		name_sliced []string

		basename          string
		extension         string
		without_extension string
	)

	ret = new(ParsedTarballName)

	basename = path.Base(full_path_or_basename)

	extension = ""
	for _, ii := range acceptable_extensions {
		if strings.HasSuffix(basename, ii) {
			extension = ii
			break
		}
	}

	// if len(extension) == 0 && !allow_non_extension {
	// 	return nil, errors.New("not a tarball extension")
	// }

	version_finder_function, version_splitter_function =
		versions_selector_function(basename)

	without_extension = basename[:len(basename)-len(extension)]

	name_sliced = TARBALL_NAME_SEPARATION.FindAllString(without_extension, -1)

	most_possible_version, found = version_finder_function(name_sliced)

	if !found {
		ret = nil
		return nil, errors.New("not found version information in tarball name")
	}

	ret.Basename = basename

	ret.Name = strings.Join(name_sliced[:most_possible_version[0]], "")

	ret.Extension = extension

	for _, i := range ALL_DELIMITERS {
		for true {
			if strings.HasPrefix(ret.Name, i) {
				ret.Name = ret.Name[1:]
			} else {
				break
			}
		}

		for true {
			if strings.HasSuffix(ret.Name, i) {
				ret.Name = ret.Name[:len(ret.Name)-1]
			} else {
				break
			}
		}
	}

	version_splitted =
		version_splitter_function(name_sliced, *most_possible_version)

	ret.Version = *version_splitted

	ret.Status.DirtyArr = name_sliced[most_possible_version[1]:]

	ret.Status.DirtyArr = StripList(ret.Status.DirtyArr, STATUS_DELIMITERS)

	ret.Status.Arr = append(ret.Status.Arr[:0], ret.Status.DirtyArr...)

	ret.Status.Arr = RemoveItemsList(ret.Status.Arr, STATUS_DELIMITERS)

	ret.Status.Arr = StripList(ret.Status.Arr, STATUS_DELIMITERS)

	ret.Status.Str = strings.Join(ret.Status.Arr, ".")

	ret.Status.DirtyStr = strings.Join(ret.Status.DirtyArr, "")

	return ret, nil
}

// func Parse(full_path_or_basename string) (*ParsedTarballName, error) {
// 	ret, err := ParseEx(
// 		full_path_or_basename,
// 		ACCEPTABLE_TARBALL_EXTENSIONS,
// 		DefaultVersionsFunctionSelector,
// 	)
// 	return ret, err
// }

func ParseStrict(full_path_or_basename string) (*ParsedTarballName, error) {
	ret, err := ParseEx(
		full_path_or_basename,
		ACCEPTABLE_TARBALL_EXTENSIONS,
		StrictVersionsFunctionSelector,
	)
	return ret, err
}

func ParseInfoZip(full_path_or_basename string) (*ParsedTarballName, error) {
	ret, err := ParseEx(
		full_path_or_basename,
		ACCEPTABLE_TARBALL_EXTENSIONS,
		InfoZipVersionsFunctionSelector,
	)
	return ret, err
}
