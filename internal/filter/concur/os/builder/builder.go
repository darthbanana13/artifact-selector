package builder

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/internal/filter/concur"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator/handleerr"
	logger "github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator/log"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
	"github.com/darthbanana13/artifact-selector/internal/log"
)

type OSFilterBuilder struct {
	Decorators []funcdecorator.FunctionDecorator[decorator.Constructor] //TODO: This does not make the builder resusable
	OS         string
}

func NewOSBuilder() *OSFilterBuilder {
	ofb := &OSFilterBuilder{
		Decorators: []funcdecorator.FunctionDecorator[decorator.Constructor]{
			handleerr.HandleErrConstructorDecorator(),
		},
	}
	return ofb
}

func (ofb *OSFilterBuilder) WithLogger(l log.ILogger) *OSFilterBuilder {
	ofb.Decorators = append(ofb.Decorators, logger.LogConstructorDecorator(l))
	return ofb
}

func (ofb *OSFilterBuilder) WithOS(os string) *OSFilterBuilder {
	ofb.OS = os
	return ofb
}

func (ofb *OSFilterBuilder) Build() (concur.FilterFunc, error) {
	if ofb.OS == "" {
		return nil, errors.New("OS is required for OSFilterBuilder")
	}

	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		os.NewOS,
		ofb.Decorators...,
	)

	os, err := constructor(ofb.OS)
	if err != nil {
		return nil, err
	}

	return os.FilterArtifact, nil
}
