package concur

import (
	"github.com/darthbanana13/artifact-selector/pkg/fetcher"
)

func FetchArtifacts(fetch *fetcher.Fetcher, userRepo, version string) (<-chan fetcher.Artifact, fetcher.ReleaseInfo, error) {
	info, artifacts, err := fetch.FetchArtifacts(userRepo, version)

	if err != nil {
		return nil, info, err
	}

	output := make(chan fetcher.Artifact)
	go func() {
		defer close(output)
		for _, artifact := range artifacts {
			output <- artifact
		}
	}()
	return output, info, nil
}
