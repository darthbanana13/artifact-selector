package rank

import (
	"math"
	"sync"

	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type IRank interface {
	Rank(<-chan filter.Artifact, uint) <-chan filter.Artifact
}

type RankFunc func(filter.Artifact) uint

func (f RankFunc) Rank(artifacts <-chan filter.Artifact, magnitude uint) <-chan filter.Artifact {
	return RankChannel(artifacts, f, magnitude)
}

func RankChannel(artifacts <-chan filter.Artifact, rankFunk RankFunc, magnitude uint) <-chan filter.Artifact {
	output := make(chan filter.Artifact)

	go func() {
		defer close(output)

		var wg sync.WaitGroup
		defer wg.Wait()

		for artifact := range artifacts {
			wg.Add(1)
			go func(artifact filter.Artifact) {
				defer wg.Done()
				rank := rankFunk(artifact)
				userRank := math.Pow10(int(magnitude))*float64(rank) + artifact.Metadata["rank"].(float64)
				artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "rank", userRank)
				output <- artifact
			}(artifact)
		}
	}()

	return output
}

func InitRank(artifacts <-chan filter.Artifact) <-chan filter.Artifact {
	output := make(chan filter.Artifact)
	go func() {
		defer close(output)
		for artifact := range artifacts {
			artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "rank", float64(0))
			output <- artifact
		}
	}()
	return output
}
