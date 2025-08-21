package transmute

import (
	"github.com/darthbanana13/artifact-selector/pkg/fetcher"
	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

func ToFilter(artifacts <-chan fetcher.Artifact) <-chan filter.Artifact {
	output := make(chan filter.Artifact)
	go func() {
		defer close(output)
		for artifact := range artifacts {
			output <- filter.Artifact{
				Source:   artifact,
				Metadata: make(map[string]any),
			}
		}
	}()
	return output
}
