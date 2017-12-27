package tarballnameparsers

import (
	"github.com/AnimusPEXUS/utils/tarballname"
)

type TarballNameParser_InfoZip struct{}

func (self *TarballNameParser_InfoZip) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {

	return tarballname.ParseInfoZip(value)
}
