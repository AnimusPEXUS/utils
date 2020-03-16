package copy

// TODO: under construction

// NOTE: probably, the idea of such functions (Copy and DeepCopy) are pointless:
//       Copy is pointless, becouse Shallow copy faster to write by hand,
//       DeepCopy is pointless, becouse it is virtually impossible to pass some
//       list of fields to be excluded from virtually unlimitd structure nesting.
//       Also, usage of reflect it self is messy

// import (
// 	"reflect"
// )

// func CopyStructSimple(dst, src interface{}) error {

// 	src_r := reflect.ValueOf(src)
// 	dst_r := reflect.ValueOf(dst)

// 	t := src_r.Type()

// 	new_dst := reflect.NewAt(t, dst)

// 	field_count := t.NumField()

// 	for (i:=0; i != field_count ; i++ ){
// 		field:= t.Field(i)

// 		switch field.Type.Kind() {
// 			case reflect
// 		}

// 	}

// 	return nil
// }
