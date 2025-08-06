package transmute

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

func ToFilter(artifacts <-chan github.Artifact) <-chan filter.Artifact {
	output := make(chan filter.Artifact)
	go func() {
		defer close(output)
		for artifact := range artifacts {
			output <- filter.Artifact{
				Source: artifact,
				Metadata: make(map[string]any),
			}
		}
	}()
	return output
}
