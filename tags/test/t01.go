package main

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/tags"
)

func Print(tag *tags.Tags) {

	for k, v := range tag.Map() {
		fmt.Printf("  '%s'\n", k)
		for _, i := range v {
			fmt.Printf("    '%s'\n", i)
		}
	}
	fmt.Println("------------------------")
}

func main() {
	t := tags.New(
		[]string{
			"parent_project:gnome",
			"system",
			"fundamental",
		},
	)
	Print(t)

	t.Add("what?", "this")
	Print(t)

	t.DeleteGroup("")
	Print(t)
}
