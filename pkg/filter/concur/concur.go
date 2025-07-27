package concur

import (
	"sync"

	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type FilterFunc func(github.Artifact) (github.Artifact, bool)

func (f FilterFunc) Filter(artifacts <-chan github.Artifact) <-chan github.Artifact {
	return FilterChannel(artifacts, f)
}

func FilterChannel(artifacts <-chan github.Artifact, filterFunc FilterFunc) <-chan github.Artifact {
	output := make(chan github.Artifact)

	go func() {
		defer close(output)

		var wg sync.WaitGroup
		defer wg.Wait()

		for artifact := range artifacts {
			wg.Add(1)
			go func(artifact github.Artifact) {
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

