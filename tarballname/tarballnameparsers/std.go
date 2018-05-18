package tarballnameparsers

import (
	"github.com/AnimusPEXUS/utils/tarballname"
)

type TarballNameParser_Std struct{}

func (self *TarballNameParser_Std) Parse(value string) (
	*tarballname.ParsedTarballName,
	error,
) {
	return tarballname.ParseStrict(value)
}

func (self *TarballNameParser_Std) Render(value *tarballname.ParsedTarballName) (string, error) {
	return value.Render(true)
}
