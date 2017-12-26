package tarballnameparsers

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

type TarballNameParser_Std struct{}

func (self *TarballNameParser_Std) ParseName(value string) (
	*types.ParseResult,
	error,
) {

	result, err := tarballname.Parse(value)
	if err != nil {
		return nil, err
	}

	ret := new(types.ParseResult)

	ret.HaveVersion = true
	ret.Version = result.Version

	ret.HaveStatus = false
	ret.Status = result.Status

	ret.HaveBuildId = false

	ret.Name = result.Name

	return ret, nil
}
