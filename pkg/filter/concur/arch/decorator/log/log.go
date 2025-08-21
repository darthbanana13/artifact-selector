package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/arch/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogDecorator struct {
	Arch arch.IArch
	L    logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(afc decorator.Constructor) decorator.Constructor {
		return func(targetArch string) (arch.IArch, error) {
			logger.Debug("Setting", "Target Architecture", targetArch)
			af, err := afc(targetArch)
			if err != nil {
				return af, err
			}
			return NewLogDecorator(af, logger)
		}
	}
}

func NewLogDecorator(arch arch.IArch, logger logging.ILogger) (arch.IArch, error) {
	if logger == nil {
		return nil, decorator.NilArchDecoratorErr(errors.New("Logger can not be nil!"))
	}
	return &LogDecorator{
		Arch: arch,
		L:    logger,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := ld.Arch.FilterArtifact(artifact)
	ld.L.Debug("Architecture filtered", "Artifact", filteredArtifact, "Keep", keep)
	return filteredArtifact, keep
}

func (ld *LogDecorator) SetTargetArch(targetArch string) error {
	ld.L.Debug("Setting", "Target Architecture", targetArch)
	return ld.Arch.SetTargetArch(targetArch)
}
