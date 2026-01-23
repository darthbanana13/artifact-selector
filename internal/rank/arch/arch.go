package arch

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/arch"
)

type Arch struct {
}

func NewArch() *Arch {
	return &Arch{}
}

func (a *Arch) RankArtifact(artifact filter.Artifact) uint {
	if artifact.Metadata["arch"].(string) == arch.Exact {
		return 2
	}
	return 1
}
