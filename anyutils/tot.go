package anyutils

import (
	"errors"
	"reflect"
)

// NOTE: it is unadvised to use this function until
// golang gets conditional compilation. othervise this function generates excessive unused code
// https://github.com/golang/go/issues/45380
// func TraverseObjectTree001[T string | float32 | float64](object_tree any, names ...string) (T, error) {

// 	var zero_T T

// 	// value of object_tree
// 	vo_ot := reflect.ValueOf(object_tree)

// 	// current object tree leaf
// 	c_ot_l := vo_ot

// 	for _, i := range names {
// 		// fmt.Println("ii:", ii, " i:", i)

// 		if c_ot_l.Kind() == reflect.Map {
// 			c_ot_l = c_ot_l.MapIndex(reflect.ValueOf(i))
// 		} else {
// 			return zero_T, errors.New("invalid object tree structure")
// 		}
// 	}

// 	if c_ot_l.Kind() == reflect.Interface {
// 		c_ot_l = c_ot_l.Elem()
// 	}

// 	switch any(zero_T).(type) {
// 	default:
// 		panic("programming error")
// 	case string:
// 		if c_ot_l.Kind() != reflect.String {
// 			return zero_T, errors.New("invalid object tree structure")
// 		}
// 		return any(c_ot_l.String()).(T), nil
// 	case float32:
// fallthrough
// 	case float64:
// 		switch c_ot_l.Kind() {
// 		default:
// 			return zero_T, errors.New("invalid object tree structure")
// 		case reflect.Float32, reflect.Float64:
// 			return any(c_ot_l.Float()).(T), nil
// 		}
// 		return zero_T, errors.New("invalid object tree structure")
// 	}
// 	return zero_T, errors.New("programming error")
// }

// results: 0 - value, 1 - found, 2 - error
func TraverseObjectTree002(
	object_tree any,
	unwrap_last_any bool,
	not_found_not_error bool,
	names ...string,
) (any, bool, error) {
	// value of object_tree
	vo_ot := reflect.ValueOf(object_tree)

	// current object tree leaf
	c_ot_l := vo_ot

	// names_l := len(names)

	for _, i := range names {

		if c_ot_l.Kind() == reflect.Map {
			c_ot_l = c_ot_l.MapIndex(reflect.ValueOf(i))
			isvalid := c_ot_l.IsValid()
			if !isvalid {
				if not_found_not_error {
					return nil, false, nil
				} else {
					return nil, false, errors.New("item not found")
				}
			}
		} else {
			return nil, false, errors.New("invalid object tree structure")
		}
	}

	if unwrap_last_any {
		if c_ot_l.Kind() == reflect.Interface {
			c_ot_l = c_ot_l.Elem()
		}
	}

	switch c_ot_l.Kind() {
	default:
		return nil, false, errors.New("invalid object tree structure")
	case reflect.Interface:
		return c_ot_l.Interface(), true, nil
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return c_ot_l.Uint(), true, nil
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return c_ot_l.Int(), true, nil
	case reflect.String:
		return c_ot_l.String(), true, nil
	case reflect.Float32, reflect.Float64:
		return c_ot_l.Float(), true, nil
	}
}

// results: 0 - value, 1 - found, 2 - error
func TraverseObjectTree002_float64(
	object_tree any,
	unwrap_last_any bool,
	not_found_not_error bool,
	names ...string,
) (float64, bool, error) {
	res, found, err := TraverseObjectTree002(
		object_tree,
		unwrap_last_any,
		not_found_not_error,
		names...,
	)

	if !found {
		if not_found_not_error {
			return 0, false, nil
		} else {
			return 0, false, errors.New("not found")
		}
	}

	if err != nil {
		return 0, found, err
	}

	var ret float64
	var ok bool

	switch res.(type) {
	default:
		return ret, found, errors.New("no type match")
	case float32:
		ret_x, ok := res.(float32)
		if !ok {
			return 0, found, errors.New("can't obtain float32")
		}
		ret = float64(ret_x)
		return ret, found, nil
	case float64:
		ret, ok = res.(float64)
		if !ok {
			return 0, false, errors.New("can't obtain float64")
		}
		return ret, found, nil
	}
}

func TraverseObjectTree002_string(
	object_tree any,
	unwrap_last_any bool,
	not_found_not_error bool,
	names ...string,
) (string, bool, error) {
	res, found, err := TraverseObjectTree002(
		object_tree,
		unwrap_last_any,
		not_found_not_error,
		names...,
	)

	if !found {
		if not_found_not_error {
			return "", false, nil
		} else {
			return "", false, errors.New("not found")
		}
	}

	if err != nil {
		return "", found, err
	}

	var ret string
	var ok bool

	switch res.(type) {
	default:
		return ret, false, errors.New("no type match")
	case string:
		ret, ok = res.(string)
		if !ok {
			return "", false, errors.New("can't obtain float64")
		}
		return ret, true, nil
	}
}

func TraverseObjectTree002_int64(
	object_tree any,
	unwrap_last_any bool,
	not_found_not_error bool,
	names ...string,
) (int64, bool, error) {
	res, found, err := TraverseObjectTree002(
		object_tree,
		unwrap_last_any,
		not_found_not_error,
		names...,
	)

	if !found {
		if not_found_not_error {
			return 0, false, nil
		} else {
			return 0, false, errors.New("not found")
		}
	}

	if err != nil {
		return 0, found, err
	}

	var ret int64
	var ok bool

	switch res.(type) {
	default:
		return ret, false, errors.New("no type match")
	case int, int8, int16, int32, int64:
		ret, ok = res.(int64)
		if !ok {
			return 0, false, errors.New("can't obtain float64")
		}
		return ret, true, nil
	}
}

func TraverseObjectTree002_uint64(
	object_tree any,
	unwrap_last_any bool,
	not_found_not_error bool,
	names ...string,
) (uint64, bool, error) {
	res, found, err := TraverseObjectTree002(
		object_tree,
		unwrap_last_any,
		not_found_not_error,
		names...,
	)

	if !found {
		if not_found_not_error {
			return 0, false, nil
		} else {
			return 0, false, errors.New("not found")
		}
	}

	if err != nil {
		return 0, found, err
	}

	var ret uint64
	var ok bool

	switch res.(type) {
	default:
		return ret, false, errors.New("no type match")
	case uint, uint8, uint16, uint32, uint64:
		ret, ok = res.(uint64)
		if !ok {
			return 0, false, errors.New("can't obtain float64")
		}
		return ret, true, nil
	}
}
