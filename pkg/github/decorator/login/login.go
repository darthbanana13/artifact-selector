package login

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/darthbanana13/artifact-selector/pkg/github"
)

var (
	// Regular expressions for known GitHub token formats
	classicTokenRegex     = regexp.MustCompile(`^(ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{36}$`)
	fineGrainedTokenRegex = regexp.MustCompile(`^github_pat_[A-Za-z0-9_]{75,82}$`)
)

type LoginDecorator struct {
	Ft github.FetcherTemplate
	Token  string
	github.IFetcher
}

func NewLoginDecorator(fetcher github.IFetcher, token string) (*LoginDecorator, error) {
	if !IsValidGithubToken(token) {
		return nil, fmt.Errorf("Invalid GitHub token format")
	}
	ld := &LoginDecorator{
		IFetcher: fetcher,
		Token:        token,
	}
	ld.Ft = github.FetcherTemplate{IFetcher: ld}
	return ld, nil
}

func (ld *LoginDecorator) PrepareRequest(url string) *http.Request {
	req := ld.IFetcher.PrepareRequest(url)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ld.Token))
	return req
}

func IsValidGithubToken(token string) bool {
	return classicTokenRegex.MatchString(token) || fineGrainedTokenRegex.MatchString(token)
}
