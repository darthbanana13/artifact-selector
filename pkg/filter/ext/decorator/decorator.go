package decorator

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter/ext"
)

type Constructor func(targetExts []string) (ext.IExt, error)

type NilExtDecoratorErr error
