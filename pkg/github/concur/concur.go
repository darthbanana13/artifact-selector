package concur

import (
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

func FetchArtifacts(fetcher github.IFetcher, userRepo string) (<-chan github.Artifact, github.ReleasesInfo, error) {
	info, err := fetcher.FetchArtifacts(userRepo)

	if err != nil {
		return nil, info, err
	}

	output := make(chan github.Artifact)
	go func() {
		defer close(output)
		for _, artifact := range info.Artifacts {
			output <- artifact
		}
	}()
	return output, info, nil
}
