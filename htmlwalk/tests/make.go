package main

import (
	"fmt"
	"sort"

	"github.com/AnimusPEXUS/utils/htmlwalk"
)

func main() {
	fmt.Println("1")
	r, _ := htmlwalk.NewHTMLWalk("https", "ftp.gnu.org")
	fmt.Println("2")
	t := r.Tree("/gnu/make")
	fmt.Println("3")

	ks := make([]string, 0)
	fmt.Println("4")

	for k, _ := range t {
		ks = append(ks, k)
	}

	fmt.Println("5")

	sort.Strings(ks)

	fmt.Println("6")

	for _, i := range ks {
		fmt.Println(i)
	}
}
