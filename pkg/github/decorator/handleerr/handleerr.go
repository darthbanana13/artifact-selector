package handleerr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type HandleErrDecorator struct {
	Ft github.FetcherTemplate
	github.IFetcher
}

func NewHandleErrorDecorator(fetcher github.IFetcher) *HandleErrDecorator {
	hed := &HandleErrDecorator{
		IFetcher: fetcher,
	}
	hed.Ft = github.FetcherTemplate{IFetcher: hed}
	return hed
}

func (hed *HandleErrDecorator) FetchArtifacts(userRepo string) (github.ReleasesInfo, error) {
	return hed.Ft.FetchArtifacts(userRepo)
}

func (hed *HandleErrDecorator) MakeUrl(userRepo string, endpoint string) (string, error) {
	i := strings.Index(userRepo, "/")
	if i > 0 && i < len(userRepo)-1 && strings.Count(userRepo, "/") == 1 {
		return hed.IFetcher.MakeUrl(userRepo, endpoint)
	}
	return "", fmt.Errorf("Given repo '%v', does not have the format 'user/repo-name'", userRepo)
}

func (hed *HandleErrDecorator) ParseJson(data []byte) (github.ReleasesInfo, error) {
	info, err := hed.IFetcher.ParseJson(data)
	if err != nil {
		return info, fmt.Errorf("Could not parse Github API reponse. Has the API changed?")
	}
	return info, err
}

func (hed *HandleErrDecorator) DoGetRequest(req *http.Request) (*http.Response, error) {
	resp, err := hed.IFetcher.DoGetRequest(req)

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("Expected HTTP Response Code %v, got code %v instead.", http.StatusOK, resp.StatusCode)
	}

	return resp, err
}
