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
	Name       string
}

func LogConstructorDecorator(logger logging.ILogger, name string) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(wsc decorator.Constructor) decorator.Constructor {
		return func(maxSize uint64, percentage float64, exts []string) (extswithinsize.IWithinSize, error) {
			ws, err := wsc(maxSize, percentage, exts)
			if err != nil {
				return ws, err
			}
			filter, err := NewLogDecorator(ws, logger, name)
			if err != nil {
				return ws, err
			}
			logger.Debug("WithinSize filter created",
				"Decorator", name,
				"Percentage", percentage,
				"Max Size", maxSize,
				"Exts", exts,
			)
			return filter, nil
		}
	}
}

func NewLogDecorator(ws extswithinsize.IWithinSize, logger logging.ILogger, name string) (extswithinsize.IWithinSize, error) {
	if logger == nil {
		return nil, decorator.NilWithinSizeDecoratorErr(errors.New("Logger can not be nil!"))
	}
	return &LogDecorator{
		WithinSize: ws,
		L:          logger,
		Name:       name,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := ld.WithinSize.FilterArtifact(artifact)
	ld.L.Debug("WithinSize filtered",
		"Decorator", ld.Name,
		"Artifact", filteredArtifact,
		"Keep", keep,
	)
	return filteredArtifact, keep
}

func (ld *LogDecorator) SetTargetExts(targetExts []string) error {
	ld.L.Debug("Setting", "Decorator", ld.Name, "Target Exts", targetExts)
	return ld.WithinSize.SetTargetExts(targetExts)
}

func (ld *LogDecorator) SetMaxSize(maxSize uint64) error {
	ld.L.Debug("Setting", "Decorator", ld.Name, "Max Size", maxSize)
	return ld.WithinSize.SetMaxSize(maxSize)
}

func (ld *LogDecorator) SetPercentage(percentage float64) error {
	ld.L.Debug("Setting", "Decorator", ld.Name, "Percentage", percentage)
	return ld.WithinSize.SetPercentage(percentage)
}
