package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter/arch"
	"github.com/darthbanana13/artifact-selector/pkg/filter/arch/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogDecorator struct {
	arch arch.IArch
	l    logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(afc decorator.Constructor) decorator.Constructor {
		return func(targetArch string) (arch.IArch, error) {
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
	if arch == nil {
		err := decorator.NilArchDecoratorErr(errors.New("ArchFilter/IArch cannot be nil"))
		logger.Info(err.Error())
		return nil, err
	}
	return &LogDecorator{
		arch: arch,
		l:    logger,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	filteredArtifact, ok := ld.arch.FilterArtifact(artifact)
	ld.l.Debug("Architecture filtered", "Artifact", filteredArtifact, "Matched", ok)
	return filteredArtifact, ok
}

func (ld *LogDecorator) SetTargetArch(targetArch string) error {
	ld.l.Debug("Setting", "Target Architecture", targetArch)
	return ld.arch.SetTargetArch(targetArch)
}
