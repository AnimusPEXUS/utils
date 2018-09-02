package tarballnameparsers

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/versionorstatus"
)

type TarballNameParser_FreeImage struct{}

func (self *TarballNameParser_FreeImage) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {

	ret, err := tarballname.ParseStrict(value)
	if err != nil {
		return nil, err
	}

	ver := ret.Version.DirtyString()

	if len(ver) >= 4 {
		ver_sl := []string{
			ver[0:1],
			ver[1:3],
			ver[3:],
		}

		new_version := versionorstatus.NewParsedVersionOrStatusFromStringSlice(ver_sl, ".")

		ret.Version = new_version
	}

	return ret, nil
}

func (self *TarballNameParser_FreeImage) Render(value *tarballname.ParsedTarballName) (string, error) {
	name := ""
	if value.Name != "" {
		name = value.Name + ""
	}
	status := ""
	if value.Status.StrSliceString("") != "" {
		status = "" + value.Status.StrSliceString("")
	}

	ext := ""
	if value.Extension != "" {
		ext = value.Extension
	}

	vstr := value.Version.StrSliceString("")

	return fmt.Sprintf("%s%s%s%s", name, vstr, status, ext), nil
}
