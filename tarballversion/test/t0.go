package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/tarballversion"
)

func main() {
	test_bases := []string{
		"b-0.tar",
		"b-0.0.tar",
		"b-0.0.0.tar",
		"b-1.tar",
		"b-2.0.tar",
		"b-2.1.tar",
		"b-2.2.0.tar",
		"b-2.2.tar",
		"b-2.2.0.tar.xz",
		"b-2.2.1.tar",
		"b-2.2.2.tar",
		"b-2.3.tar",
	}

	ver_tree, err := version.NewVersionTree("b", "std")
	if err != nil {
		panic(err)
	}

	for _, i := range test_bases {
		ver_tree.Add(i)
	}

	fmt.Println(ver_tree.Basenames([]string{".tar.xz"}))
}
