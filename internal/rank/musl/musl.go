package musl

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/musl"
)

type Musl struct {
	PreferMusl bool
}

func NewPreferMusl() *Musl {
	return NewMusl(true)
}

func NewDislikeMusl() *Musl {
	return NewMusl(false)
}

func NewMusl(preferMusl bool) *Musl {
	return &Musl{
		PreferMusl: preferMusl,
	}
}

func (m *Musl) RankArtifact(artifact filter.Artifact) uint {
	if artifact.Metadata[musl.MetadataKey].(bool) == m.PreferMusl {
		return 2
	}
	return 1
}
