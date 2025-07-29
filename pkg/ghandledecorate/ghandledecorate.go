package ghandledecorate

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type HandleErrorDecorator struct {
	ft github.FetcherTemplate
	github.IFetcher
}

func NewHandleErrorDecorator(fetcher github.IFetcher) *HandleErrorDecorator {
	d := &HandleErrorDecorator{
		IFetcher: fetcher,
	}
	d.ft = github.FetcherTemplate{IFetcher: d}
	return d
}

func (d *HandleErrorDecorator) FetchArtifacts(userRepo string) (github.ReleasesInfo, error) {
	return d.ft.FetchArtifacts(userRepo)
}

func (d *HandleErrorDecorator) MakeUrl(userRepo string, endpoint string) (string, error) {
	i := strings.Index(userRepo, "/")
	if i > 0 && i < len(userRepo)-1 && strings.Count(userRepo, "/") == 1 {
		return d.IFetcher.MakeUrl(userRepo, endpoint)
	}
	return "", fmt.Errorf("Given repo '%v', does not have the format 'user/repo-name'", userRepo)
}

func (d *HandleErrorDecorator) ParseJson(data []byte) (github.ReleasesInfo, error) {
	info, err := d.IFetcher.ParseJson(data)
	if err != nil {
		return info, fmt.Errorf("Could not parse Github API reponse. Has the API changed?")
	}
	return info, err
}

func (d *HandleErrorDecorator) DoGetRequest(req *http.Request) (*http.Response, error) {
	resp, err := d.IFetcher.DoGetRequest(req)

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("Expected HTTP Response Code %v, got code %v instead.", http.StatusOK, resp.StatusCode)
	}

	return resp, err
}
