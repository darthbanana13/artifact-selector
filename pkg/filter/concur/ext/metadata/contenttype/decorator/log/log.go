package log

import (
	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

func LogDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[concur.FilterFunc] {
	return func(ctf concur.FilterFunc) concur.FilterFunc {
		return func(artifact filter.Artifact) (filter.Artifact, bool) {
			filteredArtifact, keep := ctf(artifact)
			logger.Debug("Content Type filtered",
				"Artifact", filteredArtifact,
				"Keep", keep,
				"Content Type Match", filteredArtifact.Metadata["contentType"],
			)
			return filteredArtifact, keep
		}
	}
}
