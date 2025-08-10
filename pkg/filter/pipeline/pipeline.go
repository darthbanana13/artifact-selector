package pipeline

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

func Process(artifacts <-chan filter.Artifact, filters ...filter.IFilter) <-chan filter.Artifact {
	output := artifacts

	for _, f := range filters {
		output = f.Filter(output)
	}

	return output
}
