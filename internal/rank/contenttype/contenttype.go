package contenttype

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/ext/metadata/contenttype"
)

type ContentType struct {
}

func NewContentType() *ContentType {
	return &ContentType{}
}

func (c *ContentType) RankArtifact(artifact filter.Artifact) uint {
	switch artifact.Metadata[contenttype.MetadataIndex].(string) {
	case contenttype.Match:
		return 3
	case contenttype.Unknown, contenttype.Missing:
		return 2
	}
	return 1 //contenttype.Mismatch
}
