package ext

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type IExt interface {
	SetTargetExts([]string) error
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
