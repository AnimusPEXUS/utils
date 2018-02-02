package tarballnameparsers

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

func Get(name string) (types.TarballNameParserI, error) {
	if t, ok := Index[name]; !ok {
		return nil, errors.New("name parser not found")
	} else {
		return t(), nil
	}
}

var Index = map[string](func() types.TarballNameParserI){
	"std": func() types.TarballNameParserI {
		return new(TarballNameParser_Std)
	},

	"infozip": func() types.TarballNameParserI {
		return new(TarballNameParser_InfoZip)
	},

	"openssh": func() types.TarballNameParserI {
		return new(TarballNameParser_OpenSSH)
	},

	"slang": func() types.TarballNameParserI {
		return new(TarballNameParser_Slang)
	},
}
