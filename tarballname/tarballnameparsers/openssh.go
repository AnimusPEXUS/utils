package tarballnameparsers

import (
	"fmt"
	"strings"

	"github.com/AnimusPEXUS/utils/tarballname"
)

type TarballNameParser_OpenSSH struct{}

func (self *TarballNameParser_OpenSSH) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {
	return tarballname.ParseStrict(value)
}

func (self *TarballNameParser_OpenSSH) Render(value *tarballname.ParsedTarballName) (string, error) {
	name := ""
	if value.Name != "" {
		name = value.Name + "-"
	}
	status := ""
	if value.Status.StrSliceString("") != "" {
		str_arr := value.Status.StrSlice()
		if len(str_arr) > 0 {
			status += strings.ToLower(str_arr[0]) + str_arr[1]
			if len(str_arr) > 2 {
				status += "-" + strings.Join(str_arr[2:], ".")
			}
		}
	}

	ext := ""
	if value.Extension != "" {
		ext = value.Extension
	}

	vstr, err := value.Version.IntSliceString(".")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s%s%s", name, vstr, status, ext), nil
}
