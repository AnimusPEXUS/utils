package tarballnameparsers

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname"
)

type TarballNameParser_OpenSSH struct{}

func (self *TarballNameParser_OpenSSH) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {

	res, err := tarballname.ParseStrict(value)
	if err != nil {
		return nil, err
	}

	len_arr := len(res.Status.Arr)

	if !(len_arr == 0 || len_arr == 2) {
		return nil, errors.New("invalid number of elements in status")
	}

	if len_arr == 0 {
		return res, nil
	}

	{
		p_num := res.Status.Arr[1]

		version_sep := res.Version.DirtyArr[1]

		res.Status.DirtyArr = []string{}
		res.Status.Arr = []string{}
		res.Status.DirtyStr = ""
		res.Status.Str = ""

		res.Version.Arr = append(res.Version.Arr, p_num)
		res.Version.DirtyArr = append(res.Version.DirtyArr, []string{version_sep, p_num}...)
		res.Version.DirtyStr += version_sep + p_num
		res.Version.Str += "." + p_num
	}

	return res, nil
}
