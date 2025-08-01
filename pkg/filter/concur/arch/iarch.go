package arch

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

type IArch interface {
	SetTargetArch(string) error
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
