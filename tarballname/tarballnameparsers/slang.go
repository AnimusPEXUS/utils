package tarballnameparsers

import (
	"strconv"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/versionorstatus"
)

type TarballNameParser_Slang struct{}

func (self *TarballNameParser_Slang) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {

	res, err := tarballname.ParseStrict(value)
	if err != nil {
		return nil, err
	}

	len_arr := len(res.Status.StrSlice())

	if len_arr == 0 {
		return res, nil
	}

	{
		letter := res.Status.StrSlice()[0]

		var letterb byte = byte(letter[0])

		ver_num := int(letterb)
		ver_num_str := strconv.Itoa(ver_num)

		res.Status =
			versionorstatus.NewParsedVersionOrStatusFromStringSlice([]string{}, "")

		res.Version =
			versionorstatus.NewParsedVersionOrStatusFromStringSlice(
				append(res.Version.StrSlice(), ver_num_str),
				".",
			)
	}

	return res, nil
}
