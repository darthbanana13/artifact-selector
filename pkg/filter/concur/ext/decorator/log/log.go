package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogDecorator struct {
	Ext ext.IExt
	L   logging.ILogger
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
		Ext: ef,
		L:   logger,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := ld.Ext.FilterArtifact(artifact)
	ld.L.Debug("Ext filtered", "Artifact", filteredArtifact, "Keep", keep)
	return filteredArtifact, keep
}

func (ld *LogDecorator) SetTargetExts(targetExts []string) error {
	ld.L.Debug("Setting", "Target Exts", targetExts)
	return ld.Ext.SetTargetExts(targetExts)
}
