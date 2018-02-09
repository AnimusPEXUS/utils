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
	if value.Status.Str != "" {
		if len(value.Status.Arr) > 0 {
			status += strings.ToLower(value.Status.Arr[0]) + value.Status.Arr[1]
			if len(value.Status.Arr) > 2 {
				status += "-" + strings.Join(value.Status.Arr[2:], ".")
			}
		}
	}

	ext := ""
	if value.Extension != "" {
		ext = value.Extension
	}

	return fmt.Sprintf("%s%s%s%s", name, value.Version.Str, status, ext), nil
}
