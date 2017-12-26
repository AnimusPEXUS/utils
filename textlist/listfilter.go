package textlist

import (
	"errors"
	"regexp"
	"strings"

	"github.com/AnimusPEXUS/utils/set"
)

/*
create functions for FilterList().
parameter - parameter defined in filter and passed to function
case_sensitive - if function must take case into account
value_to_match - value which function have to check
data - can be user to pass some additional data to functions
*/
type FilterFunctions map[string]func(
	parameter string,
	case_sensitive bool,
	value_to_match string,
	data map[string]interface{},
) (bool, error)

type Filters []*FilterItem

// structure for one of items of filter text parse result
type FilterItem struct {
	// if [function result]+[NotFunc] == true.
	// if Add == true, then value which is chacked by this filter item,
	// considered to be added to result,
	// else, if Add == false, - item should be removed from result
	Add bool

	// apply boolean not to function result
	NotFunc bool

	// name of function, which FilterList have to use
	Func string

	// if function have to be case sensetive
	CaseSensitive bool

	// some functioning data, to be passed to function
	FuncParam string
}

func ParseFilterTextLinesMust(text []string) Filters {
	ret, err := ParseFilterTextLines(text)
	if err != nil {
		panic(err)
	}
	return ret
}

func ParseFilterTextMust(text string) Filters {
	ret, err := ParseFilterText(text)
	if err != nil {
		panic(err)
	}
	return ret
}

func ParseFilterText(text string) (Filters, error) {
	return ParseFilterTextLines(strings.Split(text, "\n"))
}

func ParseFilterTextLines(text []string) (Filters, error) {
	ret := make(Filters, 0)

	for _, i := range text {

		if i == "" ||
			regexp.MustCompile(`^\s+$`).Match([]byte(i)) ||
			(len(i) != 0 && string(i[0]) == `#`) {
			continue
		}

		splitted_line := strings.SplitN(i, " ", 3)

		if len(splitted_line) != 3 {
			return ret, errors.New("invalid filter text")
		}

		new_item := &FilterItem{
			Add:       splitted_line[0] == "+",
			Func:      splitted_line[1],
			FuncParam: splitted_line[2],
		}

		if strings.HasPrefix(new_item.Func, "!") {
			new_item.NotFunc = true
			new_item.Func = new_item.Func[1:]
		}

		if strings.HasSuffix(new_item.Func, "!") {
			new_item.CaseSensitive = true
			new_item.Func = new_item.Func[:len(new_item.Func)-1]
		}

		ret = append(ret, new_item)

	}

	return ret, nil
}

// TODO: im suspecting this is a garbage. commentedout on 25 Dec 2017
// var StdFunctions = FilterFunctions{}
// func FilterListStd(in_list []string, filters Filters) ([]string, error) {
// 	return FilterList(in_list, filters, StdFunctions)
// }

// Filters subject date passed by in_list, with filter set passed by filters.
// functions should contain functions asked by filters.
// data - additional data to pass to functions
func FilterList(
	in_list []string,
	filters Filters,
	functions FilterFunctions,
	data map[string]interface{},
) (
	[]string,
	error,
) {

	// fmt.Println("FilterList")
	// fmt.Println("in_list", in_list)
	// fmt.Println("filters", filters)
	// fmt.Println("functions", functions)

	out_list := set.NewSetString()

	for _, i := range in_list {
		out_list.Add(i)
	}

	for _, i := range filters {
		if _, ok := functions[i.Func]; !ok {
			return out_list.ListStrings(),
				errors.New("requested function not found: " + i.Func)
		}
	}

	for _, filter := range filters {
		funct := functions[filter.Func]

		for _, line := range in_list {
			matched, err := funct(
				filter.FuncParam,
				filter.CaseSensitive,
				line,
				data,
			)
			if err != nil {
				return out_list.ListStrings(), err
			}

			if filter.NotFunc {
				matched = !matched
			}

			if matched {
				if filter.Add {
					out_list.Add(line)
				} else {
					out_list.Remove(line)
				}
			}
		}

	}

	return out_list.ListStrings(), nil
}
