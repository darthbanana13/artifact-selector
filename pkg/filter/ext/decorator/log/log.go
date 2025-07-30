package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter/ext"
	"github.com/darthbanana13/artifact-selector/pkg/filter/ext/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogDecorator struct {
	ext ext.IExt
	l   logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(efc decorator.Constructor) decorator.Constructor {
		return func(targetExts []string) (ext.IExt, error) {
			ef, err := efc(targetExts)
			if err != nil {
				return ef, err
			}
			return NewLogDecorator(ef, logger)
		}
	}
}

func NewLogDecorator(ef ext.IExt, logger logging.ILogger) (ext.IExt, error) {
	if logger == nil {
		return nil, decorator.NilExtDecoratorErr(errors.New("Logger can not be nil!"))
	}
	if ef == nil {
		err := decorator.NilExtDecoratorErr(errors.New("ExtFilter/IExt cannot be nil"))
		logger.Info(err.Error())
		return nil, err
	}
	return &LogDecorator{
		ext: ef,
		l:   logger,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	filteredArtifact, ok := ld.ext.FilterArtifact(artifact)
	ld.l.Debug("Ext filtered", "Artifact", filteredArtifact, "Matched", ok)
	return filteredArtifact, ok
}

func (ld *LogDecorator) SetTargetExts(targetExts []string) error {
	ld.l.Debug("Setting", "Target Exts", targetExts)
	return ld.ext.SetTargetExts(targetExts)
}
