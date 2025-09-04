package builder

import (
	"github.com/darthbanana13/artifact-selector/internal/filter/concur"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/metadata/osver"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/metadata/osver/decorator/handleerr"
	logger "github.com/darthbanana13/artifact-selector/internal/filter/concur/os/metadata/osver/decorator/log"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
	"github.com/darthbanana13/artifact-selector/internal/log"
)

type OSVerBuilder struct {
	L log.ILogger
}

func NewOSVerBuilder() *OSVerBuilder {
	return &OSVerBuilder{}
}

func (ovb *OSVerBuilder) WithLogger(l log.ILogger) *OSVerBuilder {
	ovb.L = l
	return ovb
}

func (ovb *OSVerBuilder) makeDecorators() []funcdecorator.FunctionDecorator[concur.FilterFunc] {
	decorators := []funcdecorator.FunctionDecorator[concur.FilterFunc]{
		handleerr.HandleErrDecorator(),
	}
	if ovb.L != nil {
		decorators = append(decorators, logger.LogDecorator(ovb.L))
	}
	return decorators
}

func (ovb *OSVerBuilder) Build() (concur.FilterFunc, error) {
	decorators := ovb.makeDecorators()

	return funcdecorator.DecorateFunction(concur.FilterFunc(osver.FilterArtifact), decorators...), nil
}
