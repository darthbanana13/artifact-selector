package os

import (
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type IOS interface {
	SetTargetOS(string) error
	FilterArtifact(github.Artifact) (github.Artifact, bool)
}
