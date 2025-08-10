package handleerr

import "fmt"

type UnsupportedExtErr struct {
	Ext string
}

func (e UnsupportedExtErr) Error() string {
	return fmt.Sprintf("unsupported extension: %s", e.Ext)
}

type EmptyExtsErr error
