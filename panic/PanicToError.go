package panic

import (
	"errors"
	"fmt"
)

// call this inside defer, to convert panic to error
func PanicToError() error {
	e := recover()
	if e != nil {
		switch e.(type) {
		case error:
			return e.(error)
		case string:
			return errors.New(e.(string))
		default:
			return errors.New(fmt.Sprint(e))
		}
	}
	return nil
}
