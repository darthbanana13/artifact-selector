package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os/decorator"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/internal/log"
)

type LogDecorator struct {
	OS os.IOS
	L  logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(ofc decorator.Constructor) decorator.Constructor {
		return func(targetOS string) (os.IOS, error) {
			of, err := ofc(targetOS)
			if err != nil {
				return of, err
			}
			of, err = NewLogDecorator(of, logger)
			if err != nil {
				return of, err
			}
			logger.Debug("Setting", "Target OS", targetOS)
			return of, nil
		}
	}
}

func NewLogDecorator(osf os.IOS, logger logging.ILogger) (os.IOS, error) {
	if logger == nil {
		return nil, decorator.NilOSDecoratorErr(errors.New("Logger can not be nil!"))
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
		"Match Alias", filter.GetStringMetadata(filteredArtifact.Metadata, "os"),
	)
	return filteredArtifact, keep
}

func (ld *LogDecorator) SetTargetOS(targetOS string) error {
	ld.L.Debug("Setting", "Target OS", targetOS)
	return ld.OS.SetTargetOS(targetOS)
}
