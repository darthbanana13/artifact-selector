package extswithinsize

import (
	"math"
	"slices"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

type WithinSize struct {
	MaxSize    uint64
	Percentage float64
	Exts       []string
}

func NewWithinSize(maxSize uint64, percentage float64, exts []string) (IWithinSize, error) {
	return &WithinSize{
		MaxSize:    maxSize,
		Percentage: percentage,
		Exts:       exts,
	}, nil
}

func (ws *WithinSize) SetTargetExts(exts []string) error {
	ws.Exts = exts
	return nil
}

func (ws *WithinSize) SetMaxSize(maxSize uint64) error {
	ws.MaxSize = maxSize
	return nil
}

func (ws *WithinSize) SetPercentage(percentage float64) error {
	ws.Percentage = percentage
	return nil
}

func (ws *WithinSize) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	if slices.Contains(ws.Exts, artifact.Metadata["ext"].(string)) {
		if PercentDiff(artifact.Size, ws.MaxSize) <= ws.Percentage {
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
