package srstagparsers

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers"
	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

// principle is the same as far as i can tell.
// but this is not the same in the end.
// trying to make it relatively easy to complitley separate tags and tarballs
// parsers if this indeed will be nacessary.
var Index = tarballnameparsers.Index

func Get(name string) (types.TarballNameParserI, error) {
	if t, ok := Index[name]; !ok {
		return nil, errors.New("tag parser not found")
	} else {
		return t(), nil
	}
}
