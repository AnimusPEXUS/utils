package types

import (
	"github.com/AnimusPEXUS/utils/tarballname"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

type VersionComparatorI interface {
	Compare(
		tarballbasename1, tarballbasename2 *tarballname.ParsedTarballName,
	) (
		int,
		error,
	)
	Sort(tarballbasename2 []*tarballname.ParsedTarballName) error
	SortStrings(tarballbasenames_s []string, parser types.TarballNameParserI) error
}
