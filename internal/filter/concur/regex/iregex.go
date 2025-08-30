package regex

import "github.com/darthbanana13/artifact-selector/internal/filter"

type IRegex interface {
	FilterArtifact(filter.Artifact) (filter.Artifact, bool)
	SetMetadataKey(string) error
	MetadataKey() string
}
