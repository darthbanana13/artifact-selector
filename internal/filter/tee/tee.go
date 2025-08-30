package tee

import (
	"sync"

	"github.com/darthbanana13/artifact-selector/internal/filter"
)

func Tee(input <-chan filter.Artifact) (<-chan filter.Artifact, <-chan filter.Artifact) {
	out1 := make(chan filter.Artifact)
	out2 := make(chan filter.Artifact)

	go func() {
		defer close(out1)
		defer close(out2)

		var wg sync.WaitGroup
		defer wg.Wait()

		for artifact := range input {
			wg.Add(2)
			go func(a filter.Artifact) {
				defer wg.Done()
				out1 <- a
			}(artifact)
			go func(a filter.Artifact) {
				defer wg.Done()
				out2 <- a
			}(artifact)
		}
	}()
	return out1, out2
}
