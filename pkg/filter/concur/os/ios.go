package os

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

type IOS interface {
	SetTargetOS(string) error
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
