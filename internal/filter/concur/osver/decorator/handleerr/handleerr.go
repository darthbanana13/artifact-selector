package handleerr

import (
	"slices"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/os"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/osver"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
)

func HandleErrDecorator() funcdecorator.FunctionDecorator[concur.FilterFunc] {
	return func(ovf concur.FilterFunc) concur.FilterFunc {
		return func(artifact filter.Artifact) (filter.Artifact, bool) {
			osAlias := filter.GetStringMetadata(artifact.Metadata, "os")
			if slices.Contains([]string{os.Missing, filter.None}, osAlias) {
				artifact.Metadata, _ = filter.AddMetadata(artifact.Metadata, "os-ver", osver.Missing)
				return artifact, true
			}
			return ovf(artifact)
		}
	}
}

