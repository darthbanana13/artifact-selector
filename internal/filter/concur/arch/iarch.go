package arch

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type IArch interface {
	SetTargetArch(string) error
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
