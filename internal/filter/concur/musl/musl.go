package musl

import (
	"regexp"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/separator"
)

type Musl struct {
	Filter bool //if the artifact that contains musl should be filtered out or not
}

var reg *regexp.Regexp

func init() {
	reg = separator.MakeAliasRegex("musl")
}

func NewMusl(filter bool) (IMusl, error) {
	return &Musl{
		Filter: filter,
	}, nil
}

func (m *Musl) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	match := reg.MatchString(strings.ToLower(artifact.FileName))
	artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "musl", match)
	if m.Filter {
		return artifact, !match
	}
	return artifact, true
}
