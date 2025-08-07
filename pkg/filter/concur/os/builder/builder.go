package builder

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter/concur"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/os"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/os/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/os/decorator/handleerr"
	logger "github.com/darthbanana13/artifact-selector/pkg/filter/concur/os/decorator/log"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type OSFilterBuilder struct {
	decorators []funcdecorator.FunctionDecorator[decorator.Constructor]
	logger     log.ILogger
	os         string
}

func NewOSFilterBuilder() *OSFilterBuilder {
	ofb := &OSFilterBuilder{
		decorators: []funcdecorator.FunctionDecorator[decorator.Constructor]{
			handleerr.HandleErrConstructorDecorator(),
		},
	}
	return ofb
}

func (ofb *OSFilterBuilder) WithLogger(l log.ILogger) *OSFilterBuilder {
	ofb.decorators = append(ofb.decorators, logger.LogConstructorDecorator(l))
	return ofb
}

func (ofb *OSFilterBuilder) WithOS(os string) *OSFilterBuilder {
	ofb.os = os
	return ofb
}

func (ofb *OSFilterBuilder) Build() (concur.FilterFunc, error) {
	if ofb.os == "" {
		return nil, errors.New("OS is required for OSFilterBuilder")
	}

	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		os.NewOSFilter,
		ofb.decorators...,
	)

	os, err := constructor(ofb.os)
	if err != nil {
		return nil, err
	}

	return os.FilterArtifact, nil
}
