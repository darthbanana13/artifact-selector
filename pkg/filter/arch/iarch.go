package arch

import (
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type IArch interface {
	SetTargetArch(string) error
	// Filter(<-chan github.Artifact) <-chan github.Artifact
	FilterArtifact(github.Artifact) (github.Artifact, bool)
}
