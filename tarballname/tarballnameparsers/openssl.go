package tarballnameparsers

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/tarballname"
)

type TarballNameParser_OpenSSL struct{}

func (self *TarballNameParser_OpenSSL) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {
	return tarballname.ParseStrict(value)
}

func (self *TarballNameParser_OpenSSL) Render(value *tarballname.ParsedTarballName) (string, error) {
	name := ""
	if value.Name != "" {
		name = value.Name + "-"
	}
	status := ""
	if value.Status.StrSliceString("") != "" {
		status = value.Status.StrSliceString("")
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
