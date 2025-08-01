package filter

import (
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type ReleasesInfo struct {
	Version    string
	PreRelease bool
	Draft      bool
	Artifacts  []Artifact
}

type Artifact struct {
	Source   github.Artifact
	Rank     uint
	Metadata map[string]any
}

type IFilter interface {
	Filter(<-chan Artifact) <-chan Artifact
}
