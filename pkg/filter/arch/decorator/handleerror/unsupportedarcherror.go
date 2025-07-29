package handleerror

import (
	"fmt"
)

type UnsupportedArchErr struct {
	Arch string
	Err  error
}

func (e UnsupportedArchErr) Error() string {
	return fmt.Sprintf("Unsupported architecture: %s. %s", e.Arch, e.Err.Error())
}

func (e UnsupportedArchErr) Unwrap() error {
	return e.Err
}

type NilArchDecoratorErr error
