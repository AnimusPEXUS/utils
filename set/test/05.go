package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/set"
)

func main() {

	d := set.Difference(
		set.NewSet(0, 1, 2, 3, 4, 5, 6, 7, 8, 9),
		set.NewSet(1),
		set.NewSet(7),
	)

	/*
		d := goset.NewSet(0, 1, 2, 3, 4, 5, 6, 7, 8, 9).Difference(
			goset.NewSet(1),
			goset.NewSet(7),
		)
	*/

	fmt.Println(d.Len())

	for i, j := range d.List() {
		fmt.Println(i, j)
	}
}
