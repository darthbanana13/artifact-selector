package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/extswithinsize"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/extswithinsize/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogDecorator struct {
	WithinSize extswithinsize.IWithinSize
	L          logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(wsc decorator.Constructor) decorator.Constructor {
		return func(maxSize uint64, percentage float64, exts []string) (extswithinsize.IWithinSize, error) {
			ws, err := wsc(maxSize, percentage, exts)
			if err != nil {
				return ws, err
			}
			return NewLogDecorator(ws, logger)
		}
	}
}

func NewLogDecorator(ws extswithinsize.IWithinSize, logger logging.ILogger) (extswithinsize.IWithinSize, error) {
	if logger == nil {
		return nil, decorator.NilWithinSizeDecoratorErr(errors.New("Logger can not be nil!"))
	}
	if ws == nil {
		err := decorator.NilWithinSizeDecoratorErr(errors.New("WithinSizeFilter/IWithinSize cannot be nil"))
		logger.Info(err.Error())
		return nil, err
	}
	return &LogDecorator{
		WithinSize: ws,
		L:          logger,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := ld.WithinSize.FilterArtifact(artifact)
	ld.L.Debug("WithinSize filtered", "Artifact", filteredArtifact, "Keep", keep)
	return filteredArtifact, keep
}

func (ld *LogDecorator) SetTargetExts(targetExts []string) error {
	ld.L.Debug("Setting", "Target Exts", targetExts)
	return ld.WithinSize.SetTargetExts(targetExts)
}

func (ld *LogDecorator) SetMaxSize(maxSize uint64) error {
	ld.L.Debug("Setting", "Max Size", maxSize)
	return ld.WithinSize.SetMaxSize(maxSize)
}

func (ld *LogDecorator) SetPercentage(percentage float64) error {
	ld.L.Debug("Setting", "Percentage", percentage)
	return ld.WithinSize.SetPercentage(percentage)
}
