package github

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
)

const LatestReleases = "https://api.github.com/repos/%s/releases/latest"

type ReleasesInfo struct {
  Version     string      `json:"tag_name"`
  PreRelease  bool        `json:"prerelease"`
  Draft       bool        `json:"draft"`
  Artifacts   []Artifact  `json:"assets"`
}

type Artifact struct {
  InfoUrl       string  `json:"url"`
  FileName      string  `json:"name"`
  ContentType   string  `json:"content_type"`
  Size          uint64  `json:"size"`
  DownloadLink  string  `json:"browser_download_url"`
}

type IFetcher interface {
  FetchArtifacts(string) (ReleasesInfo, error)
  MakeUrl(string, string) (string, error)
  PrepareRequest(string) *http.Request
  DoGetRequest(*http.Request) (*http.Response, error)
  ReadResponseBody(io.Reader) ([]byte, error)
  ParseJson([]byte) (ReleasesInfo, error)
}

type HttpClient interface {
  Do(*http.Request) (*http.Response, error)
}

type FetcherTemplate struct {
  IFetcher
}

type HttpFetcher struct {
  *FetcherTemplate
  c HttpClient 
  t string
}

func NewHttpFetcher(client HttpClient) *HttpFetcher {
  f := &HttpFetcher{ 
    c:  client,
  }
  f.FetcherTemplate = &FetcherTemplate{f}

  return f
}

func NewDefaultHttpFetcher() *HttpFetcher {
  return NewHttpFetcher(http.DefaultClient)
}

func (f *FetcherTemplate) FetchArtifacts(userRepo string) (ReleasesInfo, error) {
  repoUrl, err := f.MakeUrl(userRepo, LatestReleases)
  if err != nil {
    return ReleasesInfo{}, err
  }

  resp, err := f.DoGetRequest(f.PrepareRequest(repoUrl))
  defer resp.Body.Close()
  if err != nil {
    return ReleasesInfo{}, err
  }

  body, err := f.ReadResponseBody(resp.Body)
  if err != nil {
    return ReleasesInfo{}, err
  }

  return f.ParseJson(body)
}

func (HttpFetcher) MakeUrl(userRepo string, releasesUrlTemplate string) (string, error) {
  return fmt.Sprintf(releasesUrlTemplate, userRepo), nil
}

func (f *HttpFetcher) PrepareRequest(url string) *http.Request {
  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
    panic("The standard library for net/http has changed")
  }
  // https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
  req.Header.Set("Accept", "application/vnd.github+json")
  req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
  // if len(f.t) > 0 {
  //   req.Header.Set("Authorization", "Bearer " + f.t)
  // }

  return req
}

func (f *HttpFetcher) DoGetRequest(req *http.Request) (*http.Response, error) {
  return f.c.Do(req)
}

func (HttpFetcher) ReadResponseBody(r io.Reader) ([]byte, error) {
  return io.ReadAll(r)
}

func (HttpFetcher) ParseJson(body []byte) (ReleasesInfo, error) {
  var info ReleasesInfo
  err := json.Unmarshal(body, &info)
  return info, err
}
