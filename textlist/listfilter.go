package textlist

import (
	"errors"
	"regexp"
	"strings"

	"github.com/AnimusPEXUS/utils/set"
)

type FilterFunctions map[string]func(
	parameter string,
	case_sensitive bool,
	value_to_match string,
) (bool, error)

type Filters []*FilterItem

type FilterItem struct {
	Add           bool
	NotFunc       bool
	Func          string
	CaseSensitive bool
	FuncParam     string
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

var StdFunctions = FilterFunctions{}

func FilterListStd(in_list []string, filters Filters) ([]string, error) {
	return FilterList(in_list, filters, StdFunctions)
}

func FilterList(in_list []string, filters Filters, functions FilterFunctions) (
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
			)
			if err != nil {
				return out_list.ListStrings(), err
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
