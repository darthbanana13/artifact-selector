package regex

import (
	"regexp"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/filter"
)

type Regex struct {
	R       *regexp.Regexp
	MetaKey string
	ToLower bool
	Filter  bool //if the artifact contains the regex, should it also filter it from the values, or just add metadata
	Exclude bool //should we exclude the value if it matches the regex, or include only the matches?
}

func NewRegex(expr, metadataKey string, toLower, filter, exclude bool) (IRegex, error) {
	r, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &Regex{
		R:       r,
		ToLower: toLower,
		Filter:  filter,
		Exclude: exclude,
		MetaKey: metadataKey,
	}, nil
}

func (r *Regex) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	filename := r.applyToLower(artifact.FileName)
	match := r.R.MatchString(filename)
	if r.MetaKey != "" {
		artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, r.MetaKey, match)
	}
	match = r.applyExclude(match)
	return r.applyFilter(artifact, match)
}

func (r *Regex) applyFilter(artifact filter.Artifact, match bool) (filter.Artifact, bool) {
	if r.Filter {
		return artifact, match
	}
	return artifact, true
}

func (r *Regex) applyExclude(match bool) bool {
	if r.Exclude {
		match = !match
	}
	return match
}

func (r *Regex) applyToLower(filename string) string {
	if r.ToLower {
		return strings.ToLower(filename)
	}
	return filename
}

func (r *Regex) SetMetadataKey(metadataKey string) error {
	r.MetaKey = metadataKey
	return nil
}

func (r *Regex) MetadataKey() string {
	return r.MetaKey
}
