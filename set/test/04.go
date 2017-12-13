package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/set"
)

func main() {

	a := set.NewSet("0", "1", "1", "2", "3", "4", "5")
	b := set.NewSet("6", "7", "7", "2", "9", "10", "11")

	c := a.Intersection(b)

	fmt.Println(c.Len())

	for i, j := range c.List() {
		fmt.Println(i, j)
	}
}
