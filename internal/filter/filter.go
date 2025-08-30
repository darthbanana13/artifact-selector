package filter

import (
	"maps"

	"github.com/darthbanana13/artifact-selector/internal/fetcher"
)

const (
	None = "none"
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

func AddMetadata(metadata map[string]any, key string, val any) map[string]any {
	var newMetadata = maps.Clone(metadata)
	newMetadata[key] = val
	return newMetadata
}

func GetStringMetadata(metadata map[string]any, key string) string {
	val, ok := metadata[key]
	if !ok {
		return None
	}
	return val.(string)
}
