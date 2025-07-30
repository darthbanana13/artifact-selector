package handleerr

import (
	"fmt"
)

type UnsupportedOSErr struct {
	OS  string
	Err error
}

func (e UnsupportedOSErr) Error() string {
	return fmt.Sprintf("Unsupported OS: %s", e.OS)
}

func (e UnsupportedOSErr) Unwrap() error {
	return e.Err
}
