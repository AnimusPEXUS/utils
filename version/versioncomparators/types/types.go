package types

import "github.com/AnimusPEXUS/utils/tarballname"

type VersionComparatorI interface {
	Compare(tarballbasename1, tarballbasename2 *tarballname.ParsedTarballName) int
}
