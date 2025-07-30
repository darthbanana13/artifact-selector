package ext

import (
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type IExt interface {
	SetTargetExts([]string) error
	FilterArtifact(github.Artifact) (github.Artifact, bool)
}
