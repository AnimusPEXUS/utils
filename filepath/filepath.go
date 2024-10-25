package filepath

import (
	"strings"

	augstrings "github.com/AnimusPEXUS/utils/strings"
)

// NOTE: this package was made by inertion after using python3 for a long time,
//       but looks like path package in golang works better than python one,
//       so looks like this package is going to be deleted

var (
	S_SEP string = "/"
	D_SEP string = "//"
)

func ReplaceDoubleSepsWithSingleOnes(value string) string {
	return strings.Replace(value, D_SEP, S_SEP, -1)
}

func TrimRightSeps(value string) string {
	return strings.TrimRight(value, S_SEP)
}

func TrimLeftSeps(value string) string {
	return strings.TrimLeft(value, S_SEP)
}

func TrimSeps(value string) string {
	return strings.Trim(value, S_SEP)
}

func StringIsSeps(value string) bool {
	if len(value) == 0 {
		return false
	}
	for i := 0; i != len(value); i++ {
		if value[i] != S_SEP[0] {
			return false
		}
	}
	return true
}

func Join(segments ...string) string {

	abs := false

	for _, i := range segments {
		if len(i) == 0 {
			continue
		}
		if i[0] == S_SEP[0] {
			abs = true
			break
		} else {
			break
		}
	}

	ret_l := []string{}

	for _, i := range segments {

		i_c := i

		if StringIsSeps(i_c) {
			continue
		}

		i_c = ReplaceDoubleSepsWithSingleOnes(i_c)
		i_c = TrimSeps(i_c)

		if len(i_c) != 0 {
			ret_l = append(ret_l, Split(i_c)...)
		} else {
			ret_l = append(ret_l, "")
		}
	}

	ret := ReplaceDoubleSepsWithSingleOnes(strings.Join(ret_l, S_SEP))

	if abs && ret[0] != S_SEP[0] {
		ret = S_SEP + ret
	}

	return ret
}

func Split(value string) []string {
	return augstrings.Split(ReplaceDoubleSepsWithSingleOnes(value), S_SEP)
}

func SelectByPreferredExtension(files []string, exts []string) string {

	for _, i := range exts {
		for _, j := range files {
			if strings.HasSuffix(j, i) {
				return j
			}
		}
	}

	return ""
}
