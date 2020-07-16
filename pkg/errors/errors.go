package errors

import (
	"errors"
	"fmt"
)

func Wrap(err error, format string, a ...interface{}) error {
	if err == nil {
		return New(format, a...)
	}
	text := fmt.Sprintf(format, a...)
	return fmt.Errorf(text+": %w", err)
}

func New(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}

func Cause(err error) error {
	for true {
		parent := errors.Unwrap(err)
		if parent == nil {
			break
		}
		err = parent
	}
	return err
}
