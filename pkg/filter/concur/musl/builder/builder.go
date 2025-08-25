package builder

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/musl"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/musl/decorator"
	logger "github.com/darthbanana13/artifact-selector/pkg/filter/concur/musl/decorator/log"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type MuslBuilder struct {
	L      log.ILogger
	Filter bool
}

func NewMuslFilterBuilder() *MuslBuilder {
	return &MuslBuilder{}
}

func (mb *MuslBuilder) WithLogger(l log.ILogger) *MuslBuilder {
	mb.L = l
	return mb
}

func (mb *MuslBuilder) WithFilter(filter bool) *MuslBuilder {
	mb.Filter = filter
	return mb
}

func (mb *MuslBuilder) Build() (concur.FilterFunc, error) {
	var constructor decorator.Constructor = musl.NewMusl
	if mb.L != nil {
		constructor = funcdecorator.DecorateFunction(constructor, logger.LogConstructorDecorator(mb.L))
	}
	muslFilter, err := constructor(mb.Filter)
	if err != nil {
		return nil, err
	}
	return muslFilter.FilterArtifact, nil
}
