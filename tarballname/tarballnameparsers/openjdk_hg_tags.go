// this parser is not for tarball name, but for openjdk tags of mercurial
// repositories

package tarballnameparsers

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/versionorstatus"
)

var (
	// "10" in OPENJDK_HG_RE_10 does not stends for JDK10, it is an index, like
	//  index line in Basic language, as I don't know how many times
	//  openjdk development changed it's tag name formatting.

	// last openjdk tag which used OPENJDK_HG_RE_10 is jdk9-b94
	OPENJDK_HG_RE_10 = regexp.MustCompile(`^(\w+)((?:\d+\.?)+)\-b(\d+)$`)

	// first openjdk tag using OPENJDK_HG_RE_20 is jdk-9+95
	OPENJDK_HG_RE_20 = regexp.MustCompile(`^(\w+)\-((?:\d+\.?)+)\+(\d+)$`)
)

type TarballNameParser_OpenJDK_Mercurial_Tags_Convertor struct {
}

func (self *TarballNameParser_OpenJDK_Mercurial_Tags_Convertor) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {
	var (
		name                 string
		openjdk_main_version *versionorstatus.ParsedVersionOrStatus
		build_number         int
	)

	// NOTE: I'm not using OPENJDK_HG_RE_10 and OPENJDK_HG_RE_20 as loop values,
	//       because with OpenJDK development incompetence, they can easily make
	//       some new strange format.

	if OPENJDK_HG_RE_10.MatchString(value) {
		var err error

		r := OPENJDK_HG_RE_10.FindStringSubmatch(value)
		name = r[1]

		openjdk_main_version = versionorstatus.NewParsedVersionOrStatusFromString(
			r[2],
			".",
		)

		build_number, err = strconv.Atoi(r[3])
		if err != nil {
			return nil, err
		}

	} else if OPENJDK_HG_RE_20.MatchString(value) {

		var err error

		r := OPENJDK_HG_RE_20.FindStringSubmatch(value)

		name = r[1]

		openjdk_main_version = versionorstatus.NewParsedVersionOrStatusFromString(
			r[2],
			".",
		)
		if err != nil {
			return nil, err
		}

		build_number, err = strconv.Atoi(r[3])
		if err != nil {
			return nil, err
		}

	} else {
		return nil, errors.New("couldn't parse OpenJDK tag")
	}

	{
		old_ver, err := openjdk_main_version.IntSlice()
		if err != nil {
			return nil, err
		}

		if len(old_ver) == 1 {
			old_ver = append(old_ver, []int{0, 0}...)
		}

		old_ver = append(old_ver, build_number)

		new_ver := versionorstatus.NewParsedVersionOrStatusFromIntSlice(old_ver, ".")

		openjdk_main_version = new_ver
	}

	ret := new(tarballname.ParsedTarballName)

	sta := versionorstatus.NewParsedVersionOrStatusFromIntSlice([]int{}, ".")

	ret.Name = name
	ret.Version = openjdk_main_version
	ret.Status = sta

	return ret, nil
}

func (self *TarballNameParser_OpenJDK_Mercurial_Tags_Convertor) Render(
	value *tarballname.ParsedTarballName,
) (
	string,
	error,
) {
	//	return fmt.Sprintf("%s-%s.tar.bz2", value.Name, value.Version.StrSliceString(".")), nil
	return "", errors.New("openjdk_hg_tags does not support rendering")
}
