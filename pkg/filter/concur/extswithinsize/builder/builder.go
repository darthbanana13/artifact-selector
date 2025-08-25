package builder

import (
	"errors"
	"sync"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/extswithinsize"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/extswithinsize/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/extswithinsize/decorator/handleerr"
	logger "github.com/darthbanana13/artifact-selector/pkg/filter/concur/extswithinsize/decorator/log"
	"github.com/darthbanana13/artifact-selector/pkg/filter/extractor/max"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type WithinSizeBuilder struct {
	L          log.ILogger
	Lname      string
	maxChann   <-chan filter.Artifact
	exts       []string
	maxSize    uint64
	percentage float64
}

func NewWithinSizeFilterBuilder() *WithinSizeBuilder {
	return &WithinSizeBuilder{}
}

func (wsb *WithinSizeBuilder) WithLogger(l log.ILogger) *WithinSizeBuilder {
	wsb.L = l
	return wsb
}

func (wsb *WithinSizeBuilder) WithLoggerName(name string) *WithinSizeBuilder {
	wsb.Lname = name
	return wsb
}

func (wsb *WithinSizeBuilder) WithMaxSize(maxSize uint64) *WithinSizeBuilder {
	wsb.maxSize = maxSize
	return wsb
}

func (wsb *WithinSizeBuilder) WithChannelMax(input <-chan filter.Artifact) *WithinSizeBuilder {
	wsb.maxChann = input
	return wsb
}

func (wsb *WithinSizeBuilder) WithExts(exts []string) *WithinSizeBuilder {
	wsb.exts = exts
	return wsb
}

func (wsb *WithinSizeBuilder) WithPercentage(percentage float64) *WithinSizeBuilder {
	wsb.percentage = percentage
	return wsb
}

func (wsb *WithinSizeBuilder) makeDecorators() []funcdecorator.FunctionDecorator[decorator.Constructor] {
	decorators := []funcdecorator.FunctionDecorator[decorator.Constructor]{
		handleerr.HandleErrConstructorDecorator(),
	}
	if wsb.L == nil {
		return decorators
	}
	return append(decorators, logger.LogConstructorDecorator(wsb.L, wsb.Lname))
}

func (wsb *WithinSizeBuilder) constructorWithDecorators() (decorator.Constructor, error) {
	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		extswithinsize.NewWithinSize,
		wsb.makeDecorators()...,
	)
	return constructor, nil
}

func BuildStrategyDeferred(
	constructor decorator.Constructor,
	maxChann <-chan filter.Artifact,
	percentage float64,
	exts []string,
) (concur.FilterFunc, error) {
	var (
		withinSizeOnce     sync.Once
		withinSizeFilter   extswithinsize.IWithinSize
		withinSizeStrategy concur.FilterFunc
	)
	withinSizeStrategy = func(artifact filter.Artifact) (filter.Artifact, bool) {
		withinSizeOnce.Do(func() { // We try to construct it before this, to assure no error is returned
			withinSizeFilter, _ = constructor(max.Find(maxChann), percentage, exts)
		})
		return withinSizeFilter.FilterArtifact(artifact)
	}
	return withinSizeStrategy, nil
}

func (wsb *WithinSizeBuilder) buildStrategy(constructor decorator.Constructor) (concur.FilterFunc, error) {
	withinSizeFilter, err := constructor(wsb.maxSize, wsb.percentage, wsb.exts)
	if err != nil {
		return nil, err
	}
	if wsb.maxSize > 0 {
		return withinSizeFilter.FilterArtifact, nil
	}
	return BuildStrategyDeferred(constructor, wsb.maxChann, wsb.percentage, wsb.exts)
}

func (wsb *WithinSizeBuilder) Build() (concur.FilterFunc, error) {
	if wsb.maxSize == 0 && wsb.maxChann == nil {
		return nil, errors.New("maxSize or maxChann must be set for WithinSizeFilterBuilder")
	}
	if wsb.maxSize > 0 && wsb.maxChann != nil {
		return nil, errors.New("maxSize and maxChann cannot both be set for WithinSizeFilterBuilder")
	}
	constructor, err := wsb.constructorWithDecorators()
	if err != nil {
		return nil, err
	}
	return wsb.buildStrategy(constructor)
}
