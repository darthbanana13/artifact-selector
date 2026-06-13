package filter

import (
	"errors"
	"maps"

	"github.com/darthbanana13/artifact-selector/internal/fetcher"
)

const (
	None = "none"
)

type ReleasesInfo struct {
	Version   string
	Artifacts []Artifact
}

type Artifact struct {
	fetcher.Artifact
	Metadata map[string]any
}

type IFilter interface {
	Filter(<-chan Artifact) <-chan Artifact
}

func AddMetadata(metadata map[string]any, valKeys ...any) (map[string]any, error) {
	if len(valKeys)%2 != 0 {
		return metadata, errors.New("metadata key-value pairs must be even")
	}
	var newMetadata = maps.Clone(metadata)
	for i := 0; i < len(valKeys)-1; i += 2 {
		key, ok := valKeys[i].(string)
		if ok == false {
			return newMetadata, errors.New("metadata key must be a string")
		}
		newMetadata[key] = valKeys[i+1]
	}
	return newMetadata, nil
}

func GetStringMetadata(metadata map[string]any, key string) string {
	val, ok := metadata[key]
	if !ok {
		return None
	}
	return val.(string)
}

func (a *Artifact) AddMetadata(valKeys ...any) error {
	metadata, ok := AddMetadata(a.Metadata, valKeys...)
	a.Metadata = metadata
	return ok
}

func (a *Artifact) GetStringMetadata(key string) string {
	return GetStringMetadata(a.Metadata, key)
}
