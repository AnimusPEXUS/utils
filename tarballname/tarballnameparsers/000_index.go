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
	"std": &TarballNameParser_Std{},

	// TODO: probably this shold be removed too (like openssl, openssh and slang)
	//       as VersionComparator facility can handle problems because of which
	//       separate name parsers was implimented. neather the less, infozip
	//       have very (and too much) special version number parsing problem.
	"infozip": &TarballNameParser_InfoZip{},

	// requires special tarballname renderer
	"openssh": &TarballNameParser_OpenSSH{},

	// requires special tarballname renderer
	"openssl": &TarballNameParser_OpenSSL{},

	// "slang": &TarballNameParser_Slang{},

	// OpenJDK team isn't smart enough to embrace standard version numbering
	"openjdk_hg_tags": &TarballNameParser_OpenJDK_Mercurial_Tags_Convertor{},
}
