package builder

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter/concur"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/decorator/handleerr"
	logger "github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/decorator/log"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type ExtFilterBuilder struct {
	decorators []funcdecorator.FunctionDecorator[decorator.Constructor]
	exts       []string
}

func NewExtFilterBuilder() *ExtFilterBuilder {
	efb := &ExtFilterBuilder{
		decorators: []funcdecorator.FunctionDecorator[decorator.Constructor]{
			handleerr.HandleErrConstructorDecorator(),
		},
	}
	return efb
}

func (efb *ExtFilterBuilder) WithLogger(l log.ILogger) *ExtFilterBuilder {
	efb.decorators = append(efb.decorators, logger.LogConstructorDecorator(l))
	return efb
}

func (efb *ExtFilterBuilder) WithExts(exts []string) *ExtFilterBuilder {
	efb.exts = exts
	return efb
}

func (efb *ExtFilterBuilder) Build() (concur.FilterFunc, error) {
	if len(efb.exts) == 0 {
		return nil, errors.New("at least one extension is required for ExtFilterBuilder")
	}

	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		ext.NewExtFilter,
		efb.decorators...,
	)

	extF, err := constructor(efb.exts)
	if err != nil {
		return nil, err
	}

	return extF.FilterArtifact, nil
}

