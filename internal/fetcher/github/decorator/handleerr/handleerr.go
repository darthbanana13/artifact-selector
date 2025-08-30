package handleerr

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/fetcher"
	"github.com/darthbanana13/artifact-selector/internal/fetcher/github/decorator"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
)

type HandleErr struct {
	fetcher.IFetcherTemplate
}

func HandleErrConstructorDecorator() funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(fc decorator.Constructor) decorator.Constructor {
		return func(client fetcher.IHTTPClient) (fetcher.IFetcherTemplate, error) {
			if client == nil {
				return nil, decorator.NilGithubDecoratorErr(errors.New("Client can not be nil"))
			}
			f, err := fc(client)
			if err != nil {
				return nil, err
			}
			return NewHandleErrorDecorator(f)
		}
	}
}

func NewHandleErrorDecorator(fetcher fetcher.IFetcherTemplate) (fetcher.IFetcherTemplate, error) {
	he := &HandleErr{
		IFetcherTemplate: fetcher,
	}
	return he, nil
}

func (he *HandleErr) MakeURL(userRepo, version string) (string, error) {
	i := strings.Index(userRepo, "/")
	if i > 0 && i < len(userRepo)-1 && strings.Count(userRepo, "/") == 1 {
		return he.IFetcherTemplate.MakeURL(userRepo, version)
	}
	return "", &InvalidGithubRepoFormat{UserRepo: userRepo}
}

func (he *HandleErr) Do(req *http.Request) (*http.Response, error) {
	resp, err := he.IFetcherTemplate.Do(req)

	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, &GithubHTTPResp{Expected: http.StatusOK, Received: resp.StatusCode}
	}

	return resp, err
}

func (he *HandleErr) ParseJson(r io.Reader) (fetcher.ReleaseInfo, []fetcher.Artifact, error) {
	info, artifacts, err := he.IFetcherTemplate.ParseJson(r)
	if err != nil {
		return info, artifacts, &GithubJSONStruct{Err: err}
	}
	return info, artifacts, err
}
