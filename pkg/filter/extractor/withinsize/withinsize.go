package withinsize

import (
	"math"
	"slices"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

type WithinSize struct {
	MaxSize uint64
	Percentage float64
	Exts []string
}

func NewWithinSize(maxSize uint64, percentage float64, exts []string) (*WithinSize, error) {
	return &WithinSize{
		MaxSize:    maxSize,
		Percentage: percentage,
		Exts:       exts,
	}, nil
}

func (fe *WithinSize) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	if slices.Contains(fe.Exts, artifact.Metadata["ext"].(string)) {
		if PercentDiff(artifact.Source.Size, fe.MaxSize) <= fe.Percentage {
			return artifact, true
		}
		return artifact, false
	}
	return artifact, true
}

func PercentDiff(a, b uint64) float64 {
	if a == 0 && b == 0 {
		return 0
	}

	fa, fb := float64(a), float64(b)

	num := math.Abs(fa - fb)
	den := (fa + fb) / 2

	return (num / den) * 100
}
