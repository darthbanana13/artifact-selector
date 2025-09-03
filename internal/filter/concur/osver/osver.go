package osver

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/filter"
)

const (
	Missing = "missing"
)

// NOTE: This doesn't work for named versions like "Ubuntu Noble Numbat" or "OpenSuse Tumbleweed". Without resorting to a
//
//	maintenance heavy list of named versions, using a user-defined regex or automatically parsing /etc/os-release for the
//	version name would be the best possible alternative
func FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	osAlias := filter.GetStringMetadata(artifact.Metadata, "os")
	r := regexp.MustCompile(fmt.Sprintf("[_\\-\\. ]%s[_\\-\\. ]?([\\d\\.].*\\d)", osAlias))
	matches := r.FindStringSubmatch(strings.ToLower(artifact.FileName))
	if len(matches) == 2 {
		artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "os-ver", matches[1])
	} else {
		artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "os-ver", Missing)
	}
	return artifact, true
}
