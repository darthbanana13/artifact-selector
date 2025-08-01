package decorator

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/os"
)

type Constructor func(targetOS string) (os.IOS, error)

type NilOSDecoratorErr error
