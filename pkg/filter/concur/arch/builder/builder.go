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
	decorators []funcdecorator.FunctionDecorator[decorator.Constructor]
	logger log.ILogger
	arch   string
}

func NewArchFilterBuilder() *ArchFilterBuilder {
	afb := &ArchFilterBuilder{
		decorators: []funcdecorator.FunctionDecorator[decorator.Constructor]{
			handleerr.HandleErrConstructorDecorator(),
		},
	}
	return afb
}

func (afb *ArchFilterBuilder) WithLogger(l log.ILogger) *ArchFilterBuilder {
	afb.decorators = append(afb.decorators, logger.LogConstructorDecorator(l))
	return afb
}

func (afb *ArchFilterBuilder) WithArch(arch string) *ArchFilterBuilder {
	afb.arch = arch
	return afb
}

func (afb *ArchFilterBuilder) Build() (concur.FilterFunc, error) {
	if afb.arch == "" {
		return nil, errors.New("Architecture is required for ArchFilterBuilder")
	}

	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		arch.NewArchFilter,
		afb.decorators...,
	)

	arch, err := constructor(afb.arch)
	if err != nil {
		return nil, err
	}

	return arch.FilterArtifact, nil
}

