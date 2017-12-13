package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/set"
)

func main() {

	a := set.NewSetString()
	for _, i := range []string{"a", "b", "strings"} {
		a.Add(i)
	}

	fmt.Println(a.ListStrings())
}
