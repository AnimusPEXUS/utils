package types

import "github.com/AnimusPEXUS/utils/tarballname"

type TarballNameParserI interface {
	Parse(value string) (*tarballname.ParsedTarballName, error)
}
