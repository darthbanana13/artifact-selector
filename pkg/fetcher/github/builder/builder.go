package builder

import (
	"errors"
	"net/http"

	"github.com/darthbanana13/artifact-selector/pkg/fetcher"
	"github.com/darthbanana13/artifact-selector/pkg/fetcher/github"
	"github.com/darthbanana13/artifact-selector/pkg/fetcher/github/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/fetcher/github/decorator/handleerr"
	glogdecorator "github.com/darthbanana13/artifact-selector/pkg/fetcher/github/decorator/log"
	"github.com/darthbanana13/artifact-selector/pkg/fetcher/github/decorator/login"
	"github.com/darthbanana13/artifact-selector/pkg/fetcher/github/retryclient"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type GithubFetcher struct {
	Decorators []funcdecorator.FunctionDecorator[decorator.Constructor] //TODO: This does not make the builder reusable
	Logger     log.ILogger
	MaxRetries int
}

func NewGihubFetcher() *GithubFetcher {
	return &GithubFetcher{
		Decorators: []funcdecorator.FunctionDecorator[decorator.Constructor]{
			handleerr.HandleErrConstructorDecorator(),
		},
	}
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
	gf.Decorators = append(gf.Decorators, login.LoginConstructorDecorator(token))
	return gf
}

func (gf *GithubFetcher) createConstructor() decorator.Constructor {
	if gf.Logger != nil {
		gf.Decorators = append(gf.Decorators, glogdecorator.LogConstructorDecorator(gf.Logger))
	}
	return funcdecorator.DecorateFunction[decorator.Constructor](
		github.NewGithubFetcher,
		gf.Decorators...,
	)
}

func (gf *GithubFetcher) Build() (fetcher.IFetcherTemplate, error) {
	constructor := gf.createConstructor()
	if gf.MaxRetries == 0 {
		return constructor(http.DefaultClient)
	}
	if gf.Logger == nil {
		return nil, NotImplemendted(errors.New("Fetcher retry without logger not implemented"))
	}

	return constructor(retryclient.NewRetryClient(gf.MaxRetries, gf.Logger))
}
