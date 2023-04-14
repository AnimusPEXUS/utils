package anyutils

import (
	"errors"
	"fmt"
	"reflect"
)

// NOTE: it is unadvised to use this function until
// golang gets conditional compilation. othervise this function generates excessive unused code
// https://github.com/golang/go/issues/45380
func TraverseObjectTree001[T string | float32 | float64](object_tree any, names ...string) (T, error) {

	var zero_T T

	// value of object_tree
	vo_ot := reflect.ValueOf(object_tree)

	// current object tree leaf
	c_ot_l := vo_ot

	for ii, i := range names {
		fmt.Println("ii:", ii, " i:", i)

		if c_ot_l.Kind() == reflect.Map {
			c_ot_l = c_ot_l.MapIndex(reflect.ValueOf(i))
		} else {
			return zero_T, errors.New("invalid object tree structure")
		}
	}

	if c_ot_l.Kind() == reflect.Interface {
		c_ot_l = c_ot_l.Elem()
	}

	switch any(zero_T).(type) {
	default:
		panic("programming error")
	case string:
		if c_ot_l.Kind() != reflect.String {
			return zero_T, errors.New("invalid object tree structure")
		}
		return any(c_ot_l.String()).(T), nil
	case float32:
	case float64:
		switch c_ot_l.Kind() {
		default:
			return zero_T, errors.New("invalid object tree structure")
		case reflect.Float32:
		case reflect.Float64:
			return any(c_ot_l.Float()).(T), nil
		}
		return zero_T, errors.New("invalid object tree structure")
	}
	return zero_T, errors.New("programming error")
}

func TraverseObjectTree002(object_tree any, unwrap_last_any bool, names ...string) (any, error) {
	// value of object_tree
	vo_ot := reflect.ValueOf(object_tree)

	// current object tree leaf
	c_ot_l := vo_ot

	for ii, i := range names {
		fmt.Println("ii:", ii, " i:", i)

		if c_ot_l.Kind() == reflect.Map {
			c_ot_l = c_ot_l.MapIndex(reflect.ValueOf(i))
		} else {
			return nil, errors.New("invalid object tree structure")
		}
	}

	if unwrap_last_any {
		if c_ot_l.Kind() == reflect.Interface {
			c_ot_l = c_ot_l.Elem()
		}
	}

	switch c_ot_l.Kind() {
	default:
		return nil, errors.New("invalid object tree structure")
	case reflect.Interface:
		return c_ot_l.Interface(), nil
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		return c_ot_l.Uint(), nil
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		return c_ot_l.Int(), nil
	case reflect.String:
		return c_ot_l.String(), nil
	case reflect.Float32:
	case reflect.Float64:
		return c_ot_l.Float(), nil
	}
	return nil, errors.New("invalid object tree structure")
}

func TraverseObjectTree002_float64(object_tree any, unwrap_last_any bool, names ...string) (float64, bool, error) {
	res, err := TraverseObjectTree002(object_tree, unwrap_last_any, names...)
	if err != nil {
		return 0, false, err
	}

	var ret float64
	var ok bool

	switch res.(type) {
	default:
		return ret, false, errors.New("no type match")
	case float32:
		ret_x, ok := res.(float32)
		if !ok {
			return 0, false, errors.New("can't obtain float32")
		}
		ret = float64(ret_x)
		return ret, true, nil
	case float64:
		ret, ok = res.(float64)
		if !ok {
			return 0, false, errors.New("can't obtain float64")
		}
		return ret, true, nil
	}
}

func TraverseObjectTree002_string(object_tree any, unwrap_last_any bool, names ...string) (string, bool, error) {
	res, err := TraverseObjectTree002(object_tree, unwrap_last_any, names...)
	if err != nil {
		return "", false, err
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
