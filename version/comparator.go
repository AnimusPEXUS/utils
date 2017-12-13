package version

import (
	"errors"

	"github.com/AnimusPEXUS/utils/tarballname/tarballnameparsers/types"
)

func Compare(one, two []int) int {

	one_c := make([]int, 0)
	for _, i := range one {
		one_c = append(one_c, i)
	}

	two_c := make([]int, 0)
	for _, i := range two {
		two_c = append(two_c, i)
	}

	for len(one_c) < len(two_c) {
		one_c = append(one_c, 0)
	}

	for len(two_c) < len(one_c) {
		two_c = append(two_c, 0)
	}

	for i := 0; i != len(one); i++ {
		if one_c[i] > two_c[i] {
			return 1
		}

		if one_c[i] < two_c[i] {
			return -1
		}
	}

	return 0
}

func SortByVersion(names []string, nameparser types.TarballNameParserI) error {
	for i := 0; i != len(names)-1; i++ {
		for j := i + 1; j != len(names); j++ {
			pi, err := nameparser.ParseName(names[i])
			if err != nil {
				return err
			}
			pj, err := nameparser.ParseName(names[j])
			if err != nil {
				return err
			}

			if pi.Name != pj.Name {
				return errors.New("by version sort name dismuch")
			}

			res := Compare(pi.VersionInt(), pj.VersionInt())
			// fmt.Println("comp res", pi.VersionInt(), res, pj.VersionInt())
			if res == 1 {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	return nil
}
