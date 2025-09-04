package sort

import (
	"sort"

	"github.com/darthbanana13/artifact-selector/internal/filter"
)

func SortChan(artifacts <-chan filter.Artifact) []filter.Artifact {
	return Sort(ChanToSlice(artifacts))
}

func Sort(artifacts []filter.Artifact) []filter.Artifact {
	sort.Slice(
		artifacts,
		func(i, j int) bool {
			return artifacts[i].Metadata["rank"].(float64) > artifacts[j].Metadata["rank"].(float64)
		},
	)
	return artifacts
}

func ChanToSlice(artifacts <-chan filter.Artifact) []filter.Artifact {
	var out []filter.Artifact
	for artifact := range artifacts {
		out = append(out, artifact)
	}
	return out
}
