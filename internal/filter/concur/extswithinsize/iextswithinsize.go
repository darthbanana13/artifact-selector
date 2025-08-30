package extswithinsize

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type IWithinSize interface {
	SetTargetExts([]string) error
	SetMaxSize(uint64) error
	SetPercentage(float64) error
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
}
