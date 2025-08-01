package handleerr

import (
	"errors"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
)

type HandleErrDecorator struct {
	Arch arch.IArch
}

func HandleErrConstructorDecorator() funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(afc decorator.Constructor) decorator.Constructor {
		return func(targetArch string) (arch.IArch, error) {
			targetArch = strings.ToLower(targetArch)
			if err := CheckValidArch(targetArch); err != nil {
				return nil, err
			}
			af, err := afc(targetArch)
			if err != nil {
				return af, err
			}
			return NewHandleErrDecorator(af)
		}
	}
}

func NewHandleErrDecorator(arch arch.IArch) (arch.IArch, error) {
	if arch == nil {
		return nil, decorator.NilArchDecoratorErr(errors.New("ArchFilter/IArch cannot be nil"))
	}
	return &HandleErrDecorator{
		Arch: arch,
	}, nil
}

func (hed *HandleErrDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	return hed.Arch.FilterArtifact(artifact)
}

func (hed *HandleErrDecorator) SetTargetArch(targetArch string) error {
	if err := CheckValidArch(targetArch); err != nil {
		return err
	}
	return hed.Arch.SetTargetArch(targetArch)
}

func CheckValidArch(archName string) error {
	if _, ok := arch.ArchMap[archName]; !ok {
		err := UnsupportedArchErr{Arch: archName}
		return err
	}
	return nil
}
