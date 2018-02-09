package sort

import (
	"errors"
	"reflect"
)

// func Sort(
// 	slices [][]interface{},
// 	control_slice_index int,
// 	compare func(i, j interface{}) (int, error),
// 	reverse bool,
// ) error {
//
// 	if !(control_slice_index >= 0 && control_slice_index < len(slices)) {
// 		return errors.New("control_slice_index out of range")
// 	}
//
// 	css_len := len(slices[control_slice_index])
//
// 	for _, i := range slices {
// 		if len(i) != css_len {
// 			return errors.New("invalid lengths of slices in 'slices'")
// 		}
// 	}
//
// 	if css_len < 2 {
// 		return nil
// 	}
//
// 	swap := func(i, j int) {
// 		for _, k := range slices {
// 			k[i], k[j] = k[j], k[i]
// 		}
// 	}
//
// 	for i := 0; i < css_len-1; i++ {
// 		for j := i + 1; j < css_len; j++ {
//
// 			ii := slices[control_slice_index][i]
// 			ij := slices[control_slice_index][j]
//
// 			res, err := compare(ii, ij)
// 			if err != nil {
// 				return err
// 			}
//
// 			if !reverse {
// 				if res < 0 {
// 					swap(i, j)
// 				}
// 			} else {
// 				if res > 0 {
// 					swap(i, j)
// 				}
// 			}
//
// 		}
// 	}
// 	return nil
// }

func Sort(
	slices interface{},
	control_slice_index int,
	compare func(i, j interface{}) (int, error),
	reverse bool,
) error {

	value_of_slices := reflect.ValueOf(slices)

	if value_of_slices.Kind() != reflect.Slice {
		// NOTE: this is type checking error, so it must be panic, not error result
		panic("'slices' must be slice")
	}

	for i := 0; i != value_of_slices.Len(); i++ {
		if reflect.ValueOf(value_of_slices.Index(i).Interface()).Kind() !=
			reflect.Slice {
			// NOTE: this is type checking error, so it must be panic, not error
			//       result
			panic("all elements of 'slices' must be slices")
		}
	}

	if !(control_slice_index >= 0 && control_slice_index < value_of_slices.Len()) {
		return errors.New("control_slice_index out of range")
	}

	controll_slice :=
		reflect.ValueOf(value_of_slices.Index(control_slice_index).Interface())
	controll_slice_len := controll_slice.Len()

	// NOTE: this must be above controll_slice_len check, since this will help
	//       user to find data structure error. also this is input value
	//       validation check, as user may assume presence of this check inside
	//       this function
	for i := 0; i != value_of_slices.Len(); i++ {
		if reflect.ValueOf(value_of_slices.Index(i).Interface()).Len() !=
			controll_slice_len {
			return errors.New("invalid lengths of slices in 'slices'")
		}
	}

	if controll_slice_len < 2 {
		return nil
	}

	swap := func(i, j int) {
		for k := 0; k != value_of_slices.Len(); k++ {
			reflect.Swapper(value_of_slices.Index(k).Interface())(i, j)
		}
	}

	for i := 0; i < controll_slice_len-1; i++ {
		for j := i + 1; j < controll_slice_len; j++ {

			ii := controll_slice.Index(i).Interface()
			ij := controll_slice.Index(j).Interface()

			res, err := compare(ii, ij)
			if err != nil {
				return err
			}

			if reverse {
				if res < 0 {
					swap(i, j)
				}
			} else {
				if res > 0 {
					swap(i, j)
				}
			}

		}
	}
	return nil
}
