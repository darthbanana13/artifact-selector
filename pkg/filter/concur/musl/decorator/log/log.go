package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/musl"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/musl/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogDecorator struct {
	Musl musl.IMusl
	L    logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(mfc decorator.Constructor) decorator.Constructor {
		return func(filter bool) (musl.IMusl, error) {
			mf, err := mfc(filter)
			if err != nil {
				return mf, err
			}
			mf, err = NewLogDecorator(mf, logger)
			if err != nil {
				return mf, err
			}
			logger.Debug("Setting", "Musl Filter with filter status", filter)
			return mf, nil
		}
	}
}

func NewLogDecorator(mf musl.IMusl, logger logging.ILogger) (musl.IMusl, error) {
	if logger == nil {
		return nil, decorator.NilMuslDecoratorErr(errors.New("Logger can not be nil!"))
	}
	return &LogDecorator{
		Musl: mf,
		L:    logger,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := ld.Musl.FilterArtifact(artifact)
	ld.L.Debug("Musl filtered",
		"Artifact", filteredArtifact,
		"Keep", keep,
		"Match Type", filteredArtifact.Metadata["musl"],
	)
	return filteredArtifact, keep
}
