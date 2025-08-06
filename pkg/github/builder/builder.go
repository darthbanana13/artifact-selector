package builder

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/github"
	"github.com/darthbanana13/artifact-selector/pkg/github/decorator/handleerr"
	glogdecorator "github.com/darthbanana13/artifact-selector/pkg/github/decorator/log"
	"github.com/darthbanana13/artifact-selector/pkg/github/decorator/login"
	"github.com/darthbanana13/artifact-selector/pkg/github/retryclient"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type GithubFetcher struct {
	Logger 			log.ILogger
	MaxRetries	int
	Token				string
}

func NewGihubFetcher() *GithubFetcher {
	return &GithubFetcher{}
}

func (gf *GithubFetcher) WithLogger(logger log.ILogger) *GithubFetcher {
	gf.Logger = logger
	return gf
}

func (gf *GithubFetcher) WithRetry(maxRetries int) *GithubFetcher {
	gf.MaxRetries = maxRetries
	return gf
}

func (gf *GithubFetcher) WithLogin(token string) *GithubFetcher {
	gf.Token = token
	return gf
}

func (gf *GithubFetcher) createWithRetry() (github.IFetcher, error) {
	var fetcher github.IFetcher
	if gf.MaxRetries == 0 {
		return github.NewDefaultHttpFetcher(), nil
	}
	if gf.Logger == nil {
		return nil, errors.New("Fetcher retry without logger not implemented")
	}
	logAdapter := retryclient.NewLeveledLoggerAdapter(gf.Logger)
	fetcher = github.NewHttpFetcher(retryclient.NewRetryClient(gf.MaxRetries, logAdapter))
	return fetcher, nil
}

func (gf *GithubFetcher) applyLogin(fetcher github.IFetcher) (github.IFetcher, error) {
	if gf.Token == "" {
		return fetcher, nil
	}
	return login.NewLoginDecorator(fetcher, gf.Token)
}

func (gf *GithubFetcher) applyLog(fetcher github.IFetcher) github.IFetcher {
	if gf.Logger == nil {
		return fetcher
	}
	return glogdecorator.NewLogFetcherDecorator(fetcher, gf.Logger)
}


func (gf *GithubFetcher) Build() (github.IFetcher, error) {
	fetcher, err := gf.createWithRetry()

	if err != nil {
		return fetcher, err
	}

	fetcher, err = gf.applyLogin(fetcher)
	fetcher =	handleerr.NewHandleErrorDecorator(fetcher)
	fetcher = gf.applyLog(fetcher)

	return fetcher, err
}
