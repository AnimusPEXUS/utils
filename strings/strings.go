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
	oristrings "strings"
)

func Split(s, sep string) []string {
	if len(s) == 0 {
		return []string{}
	} else {
		return oristrings.Split(s, sep)
	}
}
