package os

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type IOS interface {
	SetTargetOS(string) error
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
