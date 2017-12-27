package types

import "github.com/AnimusPEXUS/utils/tarballname"

type TarballNameParserI interface {
	ParseName(value string) (*tarballname.ParsedTarballName, error)
}
