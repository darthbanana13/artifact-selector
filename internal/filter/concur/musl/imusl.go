package musl

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type IMusl interface {
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
