package metafilter

import (
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/AnimusPEXUS/utils/set2"
)

/*
create functions for FilterList().
parameter - parameter defined in filter and passed to function
value_to_match - value which function have to check
*/

type FilterFunction func(
	parameter string,
	value_to_match interface{},
) (bool, error)

type FilterFunctions map[string]FilterFunction

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

		ret = append(ret, new_item)

	}

	return ret, nil
}

func FilterListItem(
	item interface{},
	filters Filters,
	functions FilterFunctions,
	prefilled bool,
	check_filters bool,
) (
	remove bool,
	err error,
) {
	if check_filters {
		for _, i := range filters {
			if _, ok := functions[i.Func]; !ok {
				err = errors.New("requested function not found: " + i.Func)
				return
			}
		}
	}

	remove = !prefilled

	for _, filter := range filters {
		funct, ok := functions[filter.Func]

		if !ok {
			err = errors.New("function with this name not found")
			return
		}

		matched, err := funct(filter.FuncParam, item)
		if err != nil {
			return false, err
		}

		if filter.NotFunc {
			matched = !matched
		}

		if matched {
			remove = !filter.Add
		}

	}

	return false, nil
}

// Filters subject date passed by in_objects, with filter set passed by filters.
// functions should contain functions asked by filters.
func FilterList(
	in_objects interface{},
	filters Filters,
	functions FilterFunctions,
	prefilled bool,
	eqcheckfunc set2.EQCheckFunc,
) (
	interface{},
	error,
) {

	out_objects := set2.NewSet()
	out_objects.SetEQCheckFunc(eqcheckfunc)

	in_objects_v := reflect.ValueOf(in_objects)

	if prefilled {
		if in_objects_v.Kind() != reflect.Slice {
			return nil, errors.New("in_objects_v must be slice")
		}

		for i := 0; i != in_objects_v.Len(); i++ {
			out_objects.Add(in_objects_v.Index(i).Interface())
		}
	}

	for _, i := range filters {
		if _, ok := functions[i.Func]; !ok {
			return nil, errors.New("requested function not found: " + i.Func)
		}
	}

	for i := 0; i != in_objects_v.Len(); i++ {

		item := in_objects_v.Index(i).Interface()

		remove, err := FilterListItem(item, filters, functions, prefilled, false)
		if err != nil {
			return nil, err
		}

		if remove {
			out_objects.Remove(item)
		} else {
			out_objects.Add(item)
		}

	}

	return out_objects.List(), nil
}
