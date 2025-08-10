package metadata

import (
	"slices"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext"
)

type Ext struct {
	TargetExts []string
}

func NewExt(targetExts []string) (ext.IExt, error) {
	e := &Ext{}
	err := e.SetTargetExts(targetExts)
	return e, err
}

func (e *Ext) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	if slices.Contains(e.TargetExts, artifact.Metadata["ext"].(string)) {
		return artifact, true
	}
	return artifact, false
}

func (e *Ext) SetTargetExts(targetExts []string) error {
	e.TargetExts = targetExts
	return nil
}
