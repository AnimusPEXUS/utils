package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/AnimusPEXUS/cliapp"
)

func main() {

	var (
		results []*cliapp.GetOptResult
	)

	{
		var (
			minus_len int

			test_subjects [][]string = [][]string{
				{
					"1", "2", "3", "-4", "5", "-10=-11", "-12", "--", "-6", "-7",
				},

				// this example is from original getopt.py from wayround_i2p_utils
				{
					"a", "b", "c", "d=123", "dd=123", "-a",
					"3", "-b=3", "--c=4", "--long=5",
					"---strange=6", "--", "-e=7",
				},

				os.Args,
			}
		)

		for i, ii := range test_subjects {

			minus_len, _ = fmt.Printf("---------- Example #%02d -----------\n", i)
			fmt.Printf("%s\n", strings.Join(ii, " "))
			res := cliapp.GetOpt(ii)
			results = append(results, res)
			fmt.Printf("opts count: %d\n", len(res.Opts))
			fmt.Printf("args count: %d\n", len(res.Args))
			fmt.Println("opts:")
			for j, jj := range res.Opts {
				fmt.Printf("#%02d'%v'==", j, jj.Name)
				if jj.HasValue {
					fmt.Printf("'%v'", jj.Value)
				} else {
					fmt.Printf("none")
				}
				fmt.Println()
			}
			fmt.Println("args: ")
			for j, jj := range res.Args {
				fmt.Printf("#%02d %v\n", j, jj)
			}
			for i := 0; i != minus_len-1; i++ {
				fmt.Print("-")
			}
			fmt.Println()

		}

	}
	fmt.Println("============ \\\\\\\\\\ ============")

	{

		fmt.Println("cliapp.GetOptCheckList usage example")

		var list1 cliapp.GetOptCheckList

		list1 = append(
			list1,
			&cliapp.GetOptCheckListItem{"-4", true, false, ""},
			&cliapp.GetOptCheckListItem{"-5", true, false, ""},
			&cliapp.GetOptCheckListItem{"-10", true, true, ""},
			&cliapp.GetOptCheckListItem{"-12", true, true, ""},
		)

		res := cliapp.GetOptHaveErrors(results[0], list1, true)

		fmt.Printf("options have errors?: %v\n", res)

	}

}
