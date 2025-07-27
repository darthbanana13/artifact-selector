package decorator

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter/arch"
)

type Constructor func(targetArch string) (arch.IArch, error)

type ConstructorDecorator func(Constructor) Constructor

func DecorateConstructor(afc Constructor, decorators ...ConstructorDecorator) Constructor {
	decorated := afc
	for _, decorator := range decorators {
		decorated = decorator(decorated)
	}
	return decorated
}

