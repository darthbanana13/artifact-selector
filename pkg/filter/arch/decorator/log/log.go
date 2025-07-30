package log

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter/arch"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/filter/arch/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type ArchLogDecorator struct {
	arch arch.IArch
	l    logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(afc decorator.Constructor) decorator.Constructor {
		return func(targetArch string) (arch.IArch, error) {
			af, err := afc(targetArch)
			if err != nil {
				return af, err
			}
			return NewArchLogDecorator(af, logger), nil
		}
	}
}

func NewArchLogDecorator(arch arch.IArch, logger logging.ILogger) arch.IArch {
	return &ArchLogDecorator{
		arch: arch,
		l:    logger,
	}
}

func (ald *ArchLogDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	filteredArtifact, ok := ald.arch.FilterArtifact(artifact)
	ald.l.Debug("Filtered artifact based on architecture", "Artifact", filteredArtifact, "Matched", ok)
	return filteredArtifact, ok
}

func (ald *ArchLogDecorator) SetTargetArch(targetArch string) error {
	ald.l.Debug("Searching for architecture" + targetArch)
	return ald.arch.SetTargetArch(targetArch)
}
