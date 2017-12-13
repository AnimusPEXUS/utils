package main

import (
	"fmt"

	augfilepath "github.com/AnimusPEXUS/filepath"
)

var SPLIT_EXAMPLES []string = []string{
	"/usr/path//something/else",
	"/usr/path//something/else/",
	"usr/path//something/else",
	"",
	"/",
	"/usr/path//some𐘓thing/else/",
}

var JOIN_EXAMPLES [][]string = [][]string{
	{},
	{""},
	{"/", "usr", "/path"},
	{"usr", "/path"},
	{"usr", "/", "/path"},
	{"usr", "/", "//", "/path"},
	{"/", "usr", "/", "//", "/path"},
	{"/", "//", "usr", "/", "//", "/path"},
	{"/", "/usr", "/", "//", "/path"},
	{"/", "//", "/usr", "/", "//", "/path"},
	{"", "/", "//", "/usr", "/", "//", "/path"},
	{"", "", "/", "//", "/usr", "/", "//", "/path"},
	{"", "", "usr", "", "path"},
}

func formatStringList(s []string) string {
	var ret string
	for i, ii := range s {
		ret += "'" + ii + "'"
		if i < len(s)-1 {
			ret += ", "
		}
	}
	ret = "[" + ret + "]"
	return ret
}

func main() {

	fmt.Println("Split Examples")
	{
		for ii, i := range SPLIT_EXAMPLES {

			fmt.Printf("Example %d\n", ii)
			fmt.Printf(" '%s'\n", i)

			res := augfilepath.Split(i)
			JOIN_EXAMPLES = append(JOIN_EXAMPLES, res)

			fmt.Println("  ", formatStringList(res), "len:", len(res))
			fmt.Println()
		}
	}

	fmt.Println("-------------------------------")

	fmt.Println("Join Examples")
	{

		for ii, i := range JOIN_EXAMPLES {

			fmt.Printf("Example %d\n", ii)
			fmt.Println(" ", formatStringList(i), "len:", len(i))
			fmt.Printf("  '%s'\n", augfilepath.Join(i))
			fmt.Println()
		}
	}
}
