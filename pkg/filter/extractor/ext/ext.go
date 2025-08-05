package ext

import (
	"slices"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

type Ext struct {
	TargetExts	[]string
	Output			chan<- uint64	
}

func NewExt(targetExts []string, output chan<- uint64) (*Ext, error) {
	return &Ext{
		TargetExts: targetExts,
		Output: output,
	}, nil
}

func (e *Ext) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	if slices.Contains(e.TargetExts, artifact.Metadata["ext"].(string)) {
		e.Output <- artifact.Source.Size
	}
	return artifact, true
}

func (e *Ext) End() {
	close(e.Output)
}
