package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/regex"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/regex/decorator"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/internal/log"
)

type LogDecorator struct {
	regex.IRegex
	L     logging.ILogger
	LName string
}

func LogConstructorDecorator(
	logger logging.ILogger,
	name string,
) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(rfc decorator.Constructor) decorator.Constructor {
		return func(expr, metadataKey string, toLower, filter, exclude bool) (regex.IRegex, error) {
			rf, err := rfc(expr, metadataKey, toLower, filter, exclude)
			if err != nil {
				return rf, err
			}
			rf, err = NewLogDecorator(rf, logger, name)
			if err != nil {
				return rf, err
			}
			logger.Debug("Setting",
				"Regex", expr,
				"Decorator", name,
				"Metadata Key", metadataKey,
				"To Lower", toLower,
				"Filter", filter,
				"Exclude", exclude,
			)
			return rf, err
		}
	}
}

func NewLogDecorator(rf regex.IRegex, logger logging.ILogger, name string) (regex.IRegex, error) {
	if logger == nil {
		return nil, decorator.NilRegexDecoratorErr(errors.New("Logger can not be nil!"))
	}
	return &LogDecorator{
		IRegex: rf,
		L:      logger,
		LName:  name,
	}, nil
}

func (ld *LogDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filteredArtifact, keep := ld.IRegex.FilterArtifact(artifact)
	key := ld.IRegex.MetadataKey()
	if key == "" {
		ld.L.Debug("Regex filtered",
			"Decorator", ld.LName,
			"Artifact", filteredArtifact,
			"Keep", keep,
		)
	} else {
		ld.L.Debug("Regex filtered",
			"Decorator", ld.LName,
			"Artifact", filteredArtifact,
			"Keep", keep,
			"Match Type", filteredArtifact.Metadata[ld.IRegex.MetadataKey()],
		)
	}
	return filteredArtifact, keep
}
