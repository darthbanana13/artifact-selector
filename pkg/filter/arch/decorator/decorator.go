package decorator

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter/arch"
)

type Constructor func(targetArch string) (arch.IArch, error)
