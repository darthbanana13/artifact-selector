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
	Ext  ext.IExt
	L    logging.ILogger
	Name string
}

func LogConstructorDecorator(logger logging.ILogger, name string) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(efc decorator.Constructor) decorator.Constructor {
		return func(targetExts []string) (ext.IExt, error) {
			logger.Debug("Setting", "Decorator", name, "Target Exts", targetExts)
			ef, err := efc(targetExts)
			if err != nil {
				return ef, err
			}
			return NewLogDecorator(ef, logger, name)
		}
	}
}

func NewLogDecorator(ef ext.IExt, logger logging.ILogger, name string) (ext.IExt, error) {
	if logger == nil {
		return nil, decorator.NilExtDecoratorErr(errors.New("Logger can not be nil!"))
	}
	return &LogDecorator{
		Ext:  ef,
		L:    logger,
		Name: name,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := ld.Ext.FilterArtifact(artifact)

	ld.L.Debug("Ext filtered",
		"Decorator", ld.Name,
		"Artifact", filteredArtifact,
		"Keep", keep,
		"Matched extension", filter.GetStringMetadata(filteredArtifact.Metadata, "ext"),
	)
	return filteredArtifact, keep
}

func (ld *LogDecorator) SetTargetExts(targetExts []string) error {
	ld.L.Debug("Setting", "Decorator", ld.Name, "Target Exts", targetExts)
	return ld.Ext.SetTargetExts(targetExts)
}
