package arch

import (
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type IArch interface {
	SetTargetArch(string) error
	FilterArtifact(github.Artifact) (github.Artifact, bool)
}
