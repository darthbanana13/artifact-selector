package handleerror

import (
	"fmt"
)

type UnsupportedArchErr struct {
	Arch string
	Err  error
}

func (e UnsupportedArchErr) Error() string {
	return fmt.Sprintf("Unsupported architecture: %s", e.Arch)
}

func (e UnsupportedArchErr) Unwrap() error {
	return e.Err
}
