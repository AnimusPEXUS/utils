package tarballnameparsers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/AnimusPEXUS/utils/tarballname"
)

type TarballNameParser_OpenSSL struct{}

func (self *TarballNameParser_OpenSSL) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {

	fmt.Println("TarballNameParser_OpenSSL Parse value", value)

	res, err := tarballname.ParseStrict(value)
	if err != nil {
		return nil, err
	}

	len_arr := len(res.Status.Arr)

	if len_arr == 0 {
		return res, nil
	}

	if len_arr != 1 {
		return nil, errors.New("unacceptable status to parse tarballname")
	}

	if ok, err := regexp.MatchString(`^[a-z]{1,3}$`, res.Status.Arr[0]); err != nil {
		return nil, err
	} else {
		if !ok {
			return nil, errors.New("unacceptable status to parse tarballname")
		}
	}

	letter_versions_int := make([]int, 0)
	letter_versions_str := make([]string, 0)

	{
		stat_arr0 := res.Status.Arr[0]
		splitted_stat := strings.Split(stat_arr0, "")

		for _, i := range splitted_stat {
			ii := int(byte(i[0]) - 96)
			letter_versions_int = append(letter_versions_int, ii)
			letter_versions_str = append(letter_versions_str, strconv.Itoa(ii))
		}
	}

	fmt.Println(" status1", res.Status.Arr)
	fmt.Println(" version1", res.Version.Arr)

	{
		res.Status.DirtyArr = []string{}
		res.Status.Arr = []string{}
		res.Status.DirtyStr = ""
		res.Status.Str = ""

		res.Version.Arr = append(res.Version.Arr, letter_versions_str...)
		res.Version.DirtyArr = append(res.Version.DirtyArr, letter_versions_str...)
		res.Version.DirtyStr += "." + strings.Join(letter_versions_str, ".")
		res.Version.Str += "." + strings.Join(letter_versions_str, ".")
	}

	fmt.Println(" status2", res.Status.Arr)
	fmt.Println(" version2", res.Version.Arr)

	return res, nil
}
