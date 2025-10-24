package os

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type OS struct {
	TotalTargets int
}

func NewOS() *OS {
	return &OS{}
}

func (o *OS) RankArtifact(artifact filter.Artifact) uint {
	return uint(artifact.Metadata["os-index"].(int))
}
