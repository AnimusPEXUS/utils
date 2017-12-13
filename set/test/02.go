package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/set"
)

func main() {
	a := set.NewSet("0", "1", "1", "2", "3")

	fmt.Println(a.Len())

	for i, j := range a.List() {
		fmt.Println(i, j)
	}
}
