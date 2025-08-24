package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/os"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/os/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogDecorator struct {
	OS os.IOS
	L  logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(ofc decorator.Constructor) decorator.Constructor {
		return func(targetOS string) (os.IOS, error) {
			logger.Debug("Setting", "Target OS", targetOS)
			of, err := ofc(targetOS)
			if err != nil {
				return of, err
			}
			return NewLogDecorator(of, logger)
		}
	}
}

func NewLogDecorator(osf os.IOS, logger logging.ILogger) (os.IOS, error) {
	if logger == nil {
		return nil, decorator.NilOSDecoratorErr(errors.New("Logger can not be nil!"))
	}
	if osf == nil {
		err := decorator.NilOSDecoratorErr(errors.New("OSFilter/IOS cannot be nil"))
		logger.Info(err.Error())
		return nil, err
	}
	return &LogDecorator{
		OS: osf,
		L:  logger,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := ld.OS.FilterArtifact(artifact)
	ld.L.Debug("OS filtered",
		"Artifact", filteredArtifact,
		"Keep", keep,
		"Match Type", filter.GetStringMetadata(filteredArtifact.Metadata, "os"),
	)
	return filteredArtifact, keep
}

func (ld *LogDecorator) SetTargetOS(targetOS string) error {
	ld.L.Debug("Setting", "Target OS", targetOS)
	return ld.OS.SetTargetOS(targetOS)
}
