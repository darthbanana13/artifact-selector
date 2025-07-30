package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter/arch"
	"github.com/darthbanana13/artifact-selector/pkg/filter/arch/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type ArchLogDecorator struct {
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
			return NewArchLogDecorator(af, logger)
		}
	}
}

func NewArchLogDecorator(arch arch.IArch, logger logging.ILogger) (arch.IArch, error) {
	if logger == nil {
		return nil, decorator.NilArchDecoratorErr(errors.New("Logger can not be nil!"))
	}
	if arch == nil {
		err := decorator.NilArchDecoratorErr(errors.New("ArchFilter/IArch cannot be nil"))
		logger.Info(err.Error())
		return nil, err
	}
	return &ArchLogDecorator{
		arch: arch,
		l:    logger,
	}, nil
}

func (ald *ArchLogDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	filteredArtifact, ok := ald.arch.FilterArtifact(artifact)
	ald.l.Debug("Architecture filtering", "Artifact", filteredArtifact, "Matched", ok)
	return filteredArtifact, ok
}

func (ald *ArchLogDecorator) SetTargetArch(targetArch string) error {
	ald.l.Debug("Setting", "Target Architecture", targetArch)
	return ald.arch.SetTargetArch(targetArch)
}
