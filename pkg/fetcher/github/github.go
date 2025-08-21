package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/darthbanana13/artifact-selector/pkg/fetcher"
)

const (
	LatestRelease = "https://api.github.com/repos/%s/releases/latest"
	TagRelease    = "https://api.github.com/repos/%s/releases/tags/%s"
)

type (
	ReleaseInfo struct {
		Version    string     `json:"tag_name"`
		PreRelease bool       `json:"prerelease"`
		Draft      bool       `json:"draft"`
		Artifacts  []Artifact `json:"assets"`
	}

	Artifact struct {
		InfoUrl      string `json:"url"`
		FileName     string `json:"name"`
		ContentType  string `json:"content_type"`
		Size         uint64 `json:"size"`
		DownloadLink string `json:"browser_download_url"`
	}

	Github struct {
		fetcher.IHTTPClient
	}
)

func NewGithubFetcher(client fetcher.IHTTPClient) (fetcher.IFetcherTemplate, error) {
	f := &Github{
		IHTTPClient: client,
	}

	return f, nil
}

func (Github) MakeURL(userRepo, version string) (string, error) {
	if version == fetcher.Latest {
		return fmt.Sprintf(LatestRelease, userRepo), nil
	}
	return fmt.Sprintf(TagRelease, userRepo, version), nil
}

func (Github) PrepareRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic("The standard library for net/http has changed")
	}
	// https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	return req
}

func (Github) ParseJson(r io.Reader) (fetcher.ReleaseInfo, []fetcher.Artifact, error) {
	var info ReleaseInfo
	dec := json.NewDecoder(r)
	err := dec.Decode(&info)
	if err != nil {
		return fetcher.ReleaseInfo{}, nil, err
	}
	releases, artifacts := ConvertReleaseInfo(info)
	return releases, artifacts, nil
}

func ConvertReleaseInfo(info ReleaseInfo) (fetcher.ReleaseInfo, []fetcher.Artifact) {
	return fetcher.ReleaseInfo{
			Version:    info.Version,
			PreRelease: info.PreRelease,
			Draft:      info.Draft,
		},
		ConvertArtifacts(info.Artifacts)
}

func ConvertArtifacts(artifacts []Artifact) []fetcher.Artifact {
	result := make([]fetcher.Artifact, 0, len(artifacts))
	for _, artifact := range artifacts {
		result = append(result, fetcher.Artifact{
			InfoUrl:      artifact.InfoUrl,
			FileName:     artifact.FileName,
			ContentType:  artifact.ContentType,
			Size:         artifact.Size,
			DownloadLink: artifact.DownloadLink,
		})
	}
	return result
}
