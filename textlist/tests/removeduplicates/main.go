package main

import (
	"fmt"
)

var example []string = []string{
	"1", "2", "2", "3",
}

func main() {

	for _, i := range example {
		fmt.Printf("%v ", i)
	}

	fmt.Println()

	res := textlist.RemoveDuplicatedString(example)

	for _, i := range res {
		fmt.Printf("%v ", i)
	}

	fmt.Println()

}
