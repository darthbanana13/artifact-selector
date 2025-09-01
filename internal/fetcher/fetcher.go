package fetcher

import (
	"io"
	"net/http"
)

const Latest = "latest"

type (
	URLBuilder func(string, string) (string, error)

	RequestPreparer func(url string) *http.Request

	HTTPClient func(*http.Request) (*http.Response, error)

	JsonParser func(io.Reader) (ReleaseInfo, []Artifact, error)

	ReleaseInfo struct {
		Version    string
		PreRelease bool
		Draft      bool
	}

	Artifact struct {
		InfoUrl      string
		FileName     string
		ContentType  string
		Size         uint64
		DownloadLink string
		Checksum     string
	}

	Fetcher struct {
		IFetcherTemplate
	}
)

func NewFetcher(template IFetcherTemplate) *Fetcher {
	return &Fetcher{IFetcherTemplate: template}
}

func (f *Fetcher) FetchArtifacts(userRepo, version string) (ReleaseInfo, []Artifact, error) {
	repoURL, err := f.IFetcherTemplate.MakeURL(userRepo, version)
	if err != nil {
		return ReleaseInfo{}, nil, err
	}

	req := f.IFetcherTemplate.PrepareRequest(repoURL)

	resp, err := f.IFetcherTemplate.Do(req)
	if err != nil {
		return ReleaseInfo{}, nil, err
	}
	defer resp.Body.Close()

	return f.IFetcherTemplate.ParseJson(resp.Body)
}
