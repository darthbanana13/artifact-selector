package builder

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter/concur"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/decorator/handleerr"
	logger "github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/decorator/log"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type ArchFilterBuilder struct {
	Decorators []funcdecorator.FunctionDecorator[decorator.Constructor]
	L          log.ILogger
	Arch       string
}

func NewArchBuilder() *ArchFilterBuilder {
	afb := &ArchFilterBuilder{
		Decorators: []funcdecorator.FunctionDecorator[decorator.Constructor]{
			handleerr.HandleErrConstructorDecorator(),
		},
	}
	return afb
}

func (afb *ArchFilterBuilder) WithLogger(l log.ILogger) *ArchFilterBuilder {
	afb.Decorators = append(afb.Decorators, logger.LogConstructorDecorator(l))
	return afb
}

func (afb *ArchFilterBuilder) WithArch(arch string) *ArchFilterBuilder {
	afb.Arch = arch
	return afb
}

func (afb *ArchFilterBuilder) Build() (concur.FilterFunc, error) {
	if afb.Arch == "" {
		return nil, errors.New("Architecture is required for ArchFilterBuilder")
	}

	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		arch.NewArch,
		afb.Decorators...,
	)

	arch, err := constructor(afb.Arch)
	if err != nil {
		return nil, err
	}

	return arch.FilterArtifact, nil
}
