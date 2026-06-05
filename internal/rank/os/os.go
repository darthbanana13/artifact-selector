package os

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os"
)

type OS struct {
	TotalTargets int
}

func NewOS() *OS {
	return &OS{}
}

func (o *OS) RankArtifact(artifact filter.Artifact) uint {
	osIndex := artifact.Metadata[os.MetadataOSIndexKey].(int)
	if osIndex < 0 {
		return uint(0)
	}
	return uint(osIndex)
}
