package main

import (
	"fmt"

	"github.com/AnimusPEXUS/cliapp"
)

func test_func1(
	getopt_result *cliapp.GetOptResult,
	available_options cliapp.GetOptCheckList,
	depth_level []string,
	subnode *cliapp.AppCmdNode,
	rootnode *cliapp.AppCmdNode,
	arg0 string,
	pass_data *interface{},
) *cliapp.AppResult {
	fmt.Printf("test_func1\n")

	a := getopt_result.GetLastNamedRetOptItem("--testopt")

	if a != nil {

		fmt.Printf("--testopt recived")
		if a.HasValue {
			fmt.Printf(", value: %v\n", a.Value)
		} else {
			fmt.Printf(" without value\n")
		}
	} else {
		fmt.Printf("--testopt not recived\n")
	}

	return &cliapp.AppResult{0, "No Errors in func1", false}
}

func test_func2(
	getopt_result *cliapp.GetOptResult,
	available_options cliapp.GetOptCheckList,
	depth_level []string,
	subnode *cliapp.AppCmdNode,
	rootnode *cliapp.AppCmdNode,
	arg0 string,
	pass_data *interface{},
) *cliapp.AppResult {
	fmt.Printf("test_func2\n")

	return &cliapp.AppResult{0, "No Errors in func2", false}
}

func main() {

	test1 := new(cliapp.AppCmdNode)
	test1.Name = "test1"
	test1.ShortDescription = "short description 1"
	test1.Callable = test_func1
	test1.AvailableOptions = cliapp.GetOptCheckList{
		&cliapp.GetOptCheckListItem{
			"--testopt",
			false,
			true,
			"Some simple descrioption for non-required option",
		},
		&cliapp.GetOptCheckListItem{
			"--testreqopt",
			true,
			true,
			"Some simple descrioption for required option",
		},
	}

	test2 := new(cliapp.AppCmdNode)
	test2.Name = "test2"
	test2.ShortDescription = "short description 2"
	test2.Callable = test_func2

	test := new(cliapp.AppCmdNode)

	test.Name = "testapp"
	test.Version = "0.0"
	test.ShortDescription = "short info on command"
	test.Description = `demonstrative description
demonstrative description demonstrative description demonstrative description
 demonstrative description demonstrative description demonstrative description
  demonstrative description demonstrative description demonstrative description
   demonstrative description demonstrative description
    demonstrative description demonstrative description`
	//test.VersionInfo = "неведомая хрень. удалить?"
	test.License = "GPLv3"
	test.Developers = []string{"AnimusPEXUS"}
	test.Date = "today"
	test.ManPages = []string{
		"printf(1)", "asprintf(3)", "dprintf(3)", "puts(3)",
		"scanf(3)", "setlocale(3)", "wcrtomb(3)", "wprintf(3)",
		"locale(5)",
	}
	test.URIs = []string{
		"https://github.com/AnimusPEXUS/cliapp/tree/master/tests/app1",
	}

	test.SubCmds = append(test.SubCmds, test1, test2)

	cliapp.RunApp(test, nil)

}
