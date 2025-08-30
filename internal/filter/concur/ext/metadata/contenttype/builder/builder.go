package builder

import (
	"github.com/darthbanana13/artifact-selector/internal/filter/concur"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/ext/metadata/contenttype"
	logger "github.com/darthbanana13/artifact-selector/internal/filter/concur/ext/metadata/contenttype/decorator/log"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
	"github.com/darthbanana13/artifact-selector/internal/log"
)

type ContentTypeBuilder struct {
	L log.ILogger
}

func NewContentTypeFilterBuilder() *ContentTypeBuilder {
	return &ContentTypeBuilder{}
}

func (ctb *ContentTypeBuilder) WithLogger(l log.ILogger) *ContentTypeBuilder {
	ctb.L = l
	return ctb
}

func (ctb *ContentTypeBuilder) Build() (concur.FilterFunc, error) {
	if ctb.L == nil {
		return contenttype.FilterArtifact, nil
	}
	return funcdecorator.DecorateFunction(contenttype.FilterArtifact, logger.LogDecorator(ctb.L)), nil
}
