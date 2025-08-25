package decorator

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/musl"
)

type Constructor func(bool) (musl.IMusl, error)

type NilMuslDecoratorErr error
