package concur

import (
	"sync"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

type FilterFunc func(filter.Artifact) (filter.Artifact, bool)

func (f FilterFunc) Filter(artifacts <-chan filter.Artifact) <-chan filter.Artifact {
	return FilterChannel(artifacts, f)
}

func FilterChannel(artifacts <-chan filter.Artifact, filterFunc FilterFunc) <-chan filter.Artifact {
	output := make(chan filter.Artifact)

	go func() {
		defer close(output)

		var wg sync.WaitGroup
		defer wg.Wait()

		for artifact := range artifacts {
			wg.Add(1)
			go func(artifact filter.Artifact) {
				defer wg.Done()
				filteredArtifact, ok := filterFunc(artifact)
				if ok {
					output <- filteredArtifact
				}
			}(artifact)
		}
	}()

	return output
}
