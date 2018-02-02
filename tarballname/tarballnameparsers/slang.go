package tarballnameparsers

import (
	"strconv"

	"github.com/AnimusPEXUS/utils/tarballname"
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

	len_arr := len(res.Status.Arr)

	if len_arr == 0 {
		return res, nil
	}

	{
		letter := res.Status.Arr[0]

		var letterb byte = byte(letter[0])

		ver_num := int(letterb)
		ver_num_str := strconv.Itoa(ver_num)

		res.Status.DirtyArr = []string{}
		res.Status.Arr = []string{}
		res.Status.DirtyStr = ""
		res.Status.Str = ""

		res.Version.Arr = append(res.Version.Arr, ver_num_str)
		res.Version.DirtyArr = append(res.Version.DirtyArr, []string{".", ver_num_str}...)
		res.Version.DirtyStr += "." + ver_num_str
		res.Version.Str += "." + ver_num_str
	}

	return res, nil
}
