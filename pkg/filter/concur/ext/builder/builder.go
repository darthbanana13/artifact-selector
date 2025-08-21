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

type ExtBuilder struct {
	Exts        []string
	L           log.ILogger
	Lname       string
	Constructor decorator.Constructor //strategy
}

func NewExtFilterBuilder() *ExtBuilder {
	return &ExtBuilder{}
}

func (eb *ExtBuilder) WithLogger(l log.ILogger) *ExtBuilder {
	eb.L = l
	return eb
}

func (eb *ExtBuilder) WithLoggerName(name string) *ExtBuilder {
	eb.Lname = name
	return eb
}

func (eb *ExtBuilder) WithExts(exts []string) *ExtBuilder {
	eb.Exts = exts
	return eb
}

func (eb *ExtBuilder) WithConstructor(constructor decorator.Constructor) *ExtBuilder {
	eb.Constructor = constructor
	return eb
}

func (eb *ExtBuilder) makeDecorators() []funcdecorator.FunctionDecorator[decorator.Constructor] {
	decorators := []funcdecorator.FunctionDecorator[decorator.Constructor]{
		handleerr.HandleErrConstructorDecorator(),
	}
	if eb.L == nil {
		return decorators
	}
	return append(decorators, logger.LogConstructorDecorator(eb.L, eb.Lname))
}

func (eb *ExtBuilder) Build() (concur.FilterFunc, error) {
	if eb.Constructor == nil {
		return nil, EmptyConstructorErr(errors.New("constructor cannot be nil for ExtFilterBuilder"))
	}
	decorators := eb.makeDecorators()

	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		ext.NewExt,
		decorators...,
	)

	extF, err := constructor(eb.Exts)
	if err != nil {
		return nil, err
	}

	return extF.FilterArtifact, nil
}
