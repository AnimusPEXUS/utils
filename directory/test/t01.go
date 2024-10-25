package main

import "github.com/AnimusPEXUS/utils/directory"
import "fmt"

func main() {
	t := directory.NewTree()

	f, _ := t.MkFile("test.txt", nil)

	s, _ := t.TreeString()

	fmt.Println(s)

	ps, _ := f.PathString()

	fmt.Println("f path", ps)

	f1, err := t.GetByPath([]string{"test.txt"}, false, true, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("f1", f1)

	f2, err := t.GetByPath([]string{"", "test.txt"}, false, true, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("f2", f2)

	res, err := t.FindFile("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}
