package retryclient

import (
	"net/http"

	"github.com/darthbanana13/artifact-selector/internal/log"

	"github.com/hashicorp/go-retryablehttp"
)

func NewRetryClient(maxRetries int, logger log.ILogger) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = maxRetries
	retryClient.Logger = NewLeveledLoggerAdapter(logger)

	return retryClient.StandardClient()
}
