package ext

import (
	"slices"

	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type Ext struct {
	TargetExts []string
}

func NewExt(targetExts []string) *Ext {
	return &Ext{
		TargetExts: targetExts,
	}
}

func (e *Ext) RankArtifact(artifact filter.Artifact) uint {
	index := slices.Index(e.TargetExts, artifact.Metadata["ext"].(string))
	return uint(len(e.TargetExts) - index)
}
