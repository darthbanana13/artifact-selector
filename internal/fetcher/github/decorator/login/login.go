package login

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/darthbanana13/artifact-selector/internal/fetcher"
	"github.com/darthbanana13/artifact-selector/internal/fetcher/github/decorator"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
)

var (
	classicTokenRegex     = regexp.MustCompile(`^(ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{36}$`)
	fineGrainedTokenRegex = regexp.MustCompile(`^github_pat_[A-Za-z0-9_]{75,82}$`)
)

type LoginDecorator struct {
	fetcher.IFetcherTemplate
	Token string
}

func LoginConstructorDecorator(token string) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(fc decorator.Constructor) decorator.Constructor {
		return func(client fetcher.IHTTPClient) (fetcher.IFetcherTemplate, error) {
			f, err := fc(client)
			if err != nil {
				return nil, err
			}
			return NewLoginDecorator(f, token)
		}
	}
}

func NewLoginDecorator(fetcher fetcher.IFetcherTemplate, token string) (fetcher.IFetcherTemplate, error) {
	if !IsValidGithubToken(token) {
		return nil, decorator.NilGithubDecoratorErr(fmt.Errorf("Invalid GitHub token format"))
	}
	return &LoginDecorator{
		IFetcherTemplate: fetcher,
		Token:            token,
	}, nil
}

func (ld *LoginDecorator) PrepareRequest(url string) *http.Request {
	req := ld.IFetcherTemplate.PrepareRequest(url)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ld.Token))
	return req
}

func IsValidGithubToken(token string) bool {
	return classicTokenRegex.MatchString(token) || fineGrainedTokenRegex.MatchString(token)
}
