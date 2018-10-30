package tarballnameparsers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/versionorstatus"
)

type TarballNameParser_ReleaseVersion struct{}

func (self *TarballNameParser_ReleaseVersion) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {

	spl := strings.Split(value, "/")

	if len(spl) != 2 {
		return nil, errors.New("invalid input data")
	}

	if spl[0] != "release" {
		return nil, errors.New("invalid input data")
	}

	ver := versionorstatus.NewParsedVersionOrStatusFromString(spl[1], ".")

	r := &tarballname.ParsedTarballName{
		Name:               spl[0],
		Version:            ver,
		Status:             versionorstatus.NewParsedVersionOrStatusFromString("", "."),
		OriginalInputValue: value,
	}

	return r, nil
}

func (self *TarballNameParser_ReleaseVersion) Render(value *tarballname.ParsedTarballName) (string, error) {
	return fmt.Sprintf("%s/%s", value.Name, value.Version.StrSliceString(".")), nil
}
