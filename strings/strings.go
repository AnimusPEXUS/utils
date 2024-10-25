package strings

/*

Explanations:

https://github.com/golang/go/issues/13075
https://groups.google.com/forum/#!topic/golang-nuts/S7U1iRmuydg

-> gopherbot locked and limited conversation to collaborators on Oct 27, 2016 <-


                                I disagree!


strings.Split(strings.Join([]string{}, "/"), "/") should return []string{}, not
[]string{""}

*/

import (
	"regexp"
	oristrings "strings"
)

const ALL_SPACES_RE = `^\s*$`

var ALL_SPACES_RE_C = regexp.MustCompile(ALL_SPACES_RE)

func Split(s, sep string) []string {
	if len(s) == 0 {
		return []string{}
	} else {
		return oristrings.Split(s, sep)
	}
}

// this will return true with len(value) != 0 and all it's runes are spaces
func StringOfSpaces(value string) bool {
	if len(value) == 0 {
		return false
	}

	return ALL_SPACES_RE_C.MatchString(value)
}
