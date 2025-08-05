package extractor

import (
	"sync"
	
	"github.com/darthbanana13/artifact-selector/pkg/filter"
)

type IExtractor interface {
	FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool)
	End()
}

type Extractor struct {
	IExtractor
}

func NewExtractor(extractor IExtractor) (*Extractor, error) {
	return &Extractor{
		IExtractor: extractor,
	}, nil
}

func (e Extractor) Extract(artifacts <-chan filter.Artifact) <-chan filter.Artifact {
	return ExtractChannel(artifacts, e)
}

//TODO: This function is too much like concur.FilterChannel. Find a way to make the code DRY
func ExtractChannel(artifacts <-chan filter.Artifact, extractor IExtractor) <-chan filter.Artifact {
	output := make(chan filter.Artifact)

	go func() {
		defer close(output)
		defer extractor.End()

		var wg sync.WaitGroup
		defer wg.Wait()

		for artifact := range artifacts {
			wg.Add(1)
			go func(artifact filter.Artifact) {
				defer wg.Done()
				filteredArtifact, ok := extractor.FilterArtifact(artifact)
				if ok {
					output <- filteredArtifact
				}
			}(artifact)
		}
	}()

	return output
}
