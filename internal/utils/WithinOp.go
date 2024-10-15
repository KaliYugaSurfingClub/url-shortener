package utils

import "fmt"

func WithinOp(op string, err *error) {
	if *err != nil {
		*err = fmt.Errorf("%s: %w", op, *err)
	}
}
