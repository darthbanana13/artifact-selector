package builder

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/internal/filter/concur"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator/extos"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator/handleerr"
	logger "github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator/log"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
	"github.com/darthbanana13/artifact-selector/internal/log"
)

type OSFilterBuilder struct {
	L			log.ILogger
	OS 		string
	ExtOS	bool
}

func NewOSBuilder() *OSFilterBuilder {
	ofb := &OSFilterBuilder{}
	return ofb
}

func (ofb *OSFilterBuilder) WithLogger(l log.ILogger) *OSFilterBuilder {
	ofb.L = l
	return ofb
}

func (ofb *OSFilterBuilder) WithOS(os string) *OSFilterBuilder {
	ofb.OS = os
	return ofb
}

func (ofb *OSFilterBuilder) WithExtOS() *OSFilterBuilder {
	ofb.ExtOS = true
	return ofb
}

func (ofb *OSFilterBuilder) makeDecorators() []funcdecorator.FunctionDecorator[decorator.Constructor] {
	decorators := []funcdecorator.FunctionDecorator[decorator.Constructor]{
		handleerr.HandleErrConstructorDecorator(),
	}
	if ofb.ExtOS == true {
		decorators = append(decorators, extos.ExtOsConstructorDecorator())
	}
	if ofb.L != nil {
		decorators = append(decorators, logger.LogConstructorDecorator(ofb.L))
	}
	return decorators
}

func (ofb *OSFilterBuilder) Build() (concur.FilterFunc, error) {
	if ofb.OS == "" {
		return nil, errors.New("OS is required for OSFilterBuilder")
	}
	decorators := ofb.makeDecorators()

	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		os.NewOS,
		decorators...,
	)

	os, err := constructor(ofb.OS)
	if err != nil {
		return nil, err
	}

	return os.FilterArtifact, nil
}
