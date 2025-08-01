package ext

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

type IExt interface {
	SetTargetExts([]string) error
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
