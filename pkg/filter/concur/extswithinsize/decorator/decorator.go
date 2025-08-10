package decorator

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/extswithinsize"
)

type Constructor func(maxSize uint64, percentage float64, exts []string) (extswithinsize.IWithinSize, error)

type NilWithinSizeDecoratorErr error
