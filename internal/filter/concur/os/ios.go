package os

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type IOS interface {
	TargetOS() string
	TargetAliases() []string
	ExcludedAliases() []string
	SetTargetOS(string) error
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
