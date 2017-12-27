package tarballnameparsers

import (
	"github.com/AnimusPEXUS/utils/tarballname"
)

type TarballNameParser_Std struct{}

func (self *TarballNameParser_Std) ParseName(value string) (
	*tarballname.ParsedTarballName,
	error,
) {

	return tarballname.ParseStrict(value)
}
