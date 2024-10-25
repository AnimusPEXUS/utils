package types

import "github.com/AnimusPEXUS/utils/tarballname"

type TarballNameParserI interface {
	Parse(value string) (
		*tarballname.ParsedTarballName,
		error,
	)

	// TODO: maybe rendering should be separated to own facility..
	Render(value *tarballname.ParsedTarballName) (string, error)
}
