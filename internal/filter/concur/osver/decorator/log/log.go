package log

import (
	"github.com/darthbanana13/artifact-selector/internal/filter"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/internal/log"
)

func LogDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[concur.FilterFunc] {
	return func(ovf concur.FilterFunc) concur.FilterFunc {
		return func(artifact filter.Artifact) (filter.Artifact, bool) {
			filteredArtifact, keep := ovf(artifact)
			logger.Debug("OS Version extracted",
				"Artifact", filteredArtifact,
				"Keep", keep,
				"OS Version", filteredArtifact.Metadata["os-ver"],
			)
			return filteredArtifact, keep
		}
	}
}
