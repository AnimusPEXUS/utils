package tarballnameparsers

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

func Get(name string) (types.TarballNameParserI, error) {
	if t, ok := Index[name]; !ok {
		return nil, errors.New("name parser not found")
	} else {
		return t, nil
	}
}

var Index = map[string]types.TarballNameParserI{
	"std":     &TarballNameParser_Std{},
	"infozip": &TarballNameParser_InfoZip{},
	"openssh": &TarballNameParser_OpenSSH{},
	"openssl": &TarballNameParser_OpenSSL{},
	"slang":   &TarballNameParser_Slang{},
}
