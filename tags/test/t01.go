package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/tags"
)

func main() {
	t := tags.New([]string{"parent_project:gnome", "system", "fundamental"})
	for k, v := range t.Map() {
		fmt.Printf("  '%s'\n", k)
		for _, i := range v {
			fmt.Printf("    '%s'\n", i)
		}
	}
}
