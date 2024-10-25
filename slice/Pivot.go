package slice

import (
	"errors"
	"reflect"
)

func Pivot(slice interface{}, pice_size int) (interface{}, error) {
	r := reflect.ValueOf(slice)

	if r.Kind() != reflect.Slice {
		return nil, errors.New("slice must be slice")
	}

	ret := reflect.MakeSlice(reflect.SliceOf(reflect.SliceOf(r.Type().Elem())), 0, 0) // V

	i := 0

	for {
		i0 := i * pice_size

		len_a := r.Len()

		if i0 >= len_a {
			break
		}

		i1 := (i * pice_size) + pice_size

		if i1 > len_a {
			i1 = len_a
		}

		t_s := r.Slice(i0, i1)

		ret = reflect.Append(
			ret,
			t_s,
		)
		i++
	}

	return ret.Interface(), nil
}
