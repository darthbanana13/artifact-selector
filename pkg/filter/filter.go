package filter

import (
	"github.com/darthbanana13/artifact-selector/pkg/fetcher"
)

type ReleasesInfo struct {
	Version    string
	PreRelease bool
	Draft      bool
	Artifacts  []Artifact
}

type Artifact struct {
	fetcher.Artifact
	Metadata map[string]any
}

type IFilter interface {
	Filter(<-chan Artifact) <-chan Artifact
}
