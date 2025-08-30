package decorator

import "github.com/darthbanana13/artifact-selector/internal/filter/concur/ext"

type Constructor func(targetExts []string) (ext.IExt, error)

type NilExtDecoratorErr error
