package tarballnameparsers

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/tarballname"
)

type TarballNameParser_InfoZip struct{}

func (self *TarballNameParser_InfoZip) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {

	return tarballname.ParseInfoZip(value)
}

func (self *TarballNameParser_InfoZip) Render(value *tarballname.ParsedTarballName) (string, error) {
	name := ""
	if value.Name != "" {
		name = value.Name + ""
	}
	status := ""
	if value.Status.Str != "" {
		status = "" + value.Status.Str
	}

	ext := ""
	if value.Extension != "" {
		ext = value.Extension
	}

	return fmt.Sprintf("%s%s%s%s", name, value.Version.Str, status, ext), nil
}
