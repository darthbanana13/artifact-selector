package handleerr

import (
	"errors"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/filter/arch"
	"github.com/darthbanana13/artifact-selector/pkg/filter/arch/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type HandleErrDecorator struct {
	arch arch.IArch
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
		arch: arch,
	}, nil
}

func (hed *HandleErrDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	return hed.arch.FilterArtifact(artifact)
}

func (hed *HandleErrDecorator) SetTargetArch(targetArch string) error {
	if err := CheckValidArch(targetArch); err != nil {
		return err
	}
	return hed.arch.SetTargetArch(targetArch)
}

func CheckValidArch(archName string) error {
	if _, ok := arch.ArchMap[archName]; !ok {
		err := UnsupportedArchErr{Arch: archName}
		return err
	}
	return nil
}
