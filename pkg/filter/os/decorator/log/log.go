package log

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter/os"
	"github.com/darthbanana13/artifact-selector/pkg/filter/os/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type OSLogDecorator struct {
	os os.IOS
	l  logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(ofc decorator.Constructor) decorator.Constructor {
		return func(targetOS string) (os.IOS, error) {
			of, err := ofc(targetOS)
			if err != nil {
				return of, err
			}
			return NewOSLogDecorator(of, logger)
		}
	}
}

func NewOSLogDecorator(osf os.IOS, logger logging.ILogger) (os.IOS, error) {
	if logger == nil {
		return nil, decorator.NilOSDecoratorErr(errors.New("Logger can not be nil!"))
	}
	if osf == nil {
		err := decorator.NilOSDecoratorErr(errors.New("OSFilter/IOS cannot be nil"))
		logger.Info(err.Error())
		return nil, err
	}
	return &OSLogDecorator{
		os: osf,
		l:  logger,
	}, nil
}

func (old *OSLogDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	filteredArtifact, ok := old.os.FilterArtifact(artifact)
	old.l.Debug("OS filtered", "Artifact", filteredArtifact, "Matched", ok)
	return filteredArtifact, ok
}

func (old *OSLogDecorator) SetTargetOS(targetOS string) error {
	old.l.Debug("Setting", "Target OS", targetOS)
	return old.os.SetTargetOS(targetOS)
}

