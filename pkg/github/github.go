package github

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "strings"
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
  FileType      string  `json:"content_type"`
  Size          uint64  `json:"size"`
  DownloadLink  string  `json:"browser_download_url"`
}

type IFetcher interface {
  FetchArtifacts(string) (ReleasesInfo, error)
  TestValidRepoName(string) error
  MakeUrl(string, string) string
  PrepareRequest(string) *http.Request
  GetUrlBody(*http.Request, BodyReader) ([]byte, error)
  ReadBody(io.Reader) ([]byte, error)
  ParseJson([]byte) (ReleasesInfo, error)
}

type httpClient interface {
  Do(*http.Request) (*http.Response, error)
}

type FetcherTemplate struct {
  IFetcher
}

type HttpFetcher struct {
  *FetcherTemplate
  c httpClient 
}

func NewHttpFetcher(client httpClient) *HttpFetcher {
  f := &HttpFetcher{ 
    c:  client,
  }
  f.FetcherTemplate = &FetcherTemplate{f}

  return f
}

type BodyReader func(io.Reader) ([]byte, error)

func (f *FetcherTemplate) FetchArtifacts(userRepo string) (ReleasesInfo, error) {
  if err := f.TestValidRepoName(userRepo); err != nil {
    return ReleasesInfo{}, err
  }

  body, err := f.GetUrlBody(
    f.PrepareRequest(f.MakeUrl(userRepo, LatestReleases)),
    f.ReadBody,
  )
  if err != nil {
    return ReleasesInfo{}, err
  }

  return f.ParseJson(body)
}

func (HttpFetcher) TestValidRepoName(userRepo string) error {
  i := strings.Index(userRepo, "/")
  if i > 0 && i < len(userRepo) - 1 && strings.Count(userRepo, "/") == 1 {
    return nil 
  }
  return fmt.Errorf("Given repo '%v', does not have the format 'user/repo-name'", userRepo)
}

func (HttpFetcher) MakeUrl(userRepo string, releasesUrlTemplate string) string {
  return fmt.Sprintf(releasesUrlTemplate, userRepo)
}

func (f *HttpFetcher) PrepareRequest(url string) *http.Request {
  // https://docs.github.com/en/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release
  req, err := http.NewRequest(http.MethodGet, url, nil)
  if err != nil {
    panic("The standard library for http has changed")
  }
  req.Header.Set("Accept", "application/vnd.github+json")
  req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
  // req.Header.Set("Authorization", "Bearer <TOKEN>")

  return req
}

func (f *HttpFetcher) GetUrlBody(req *http.Request, reader BodyReader) ([]byte, error) {
  resp, err := f.c.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("Expected HTTP Response Code %v, got code %v instead.", http.StatusOK, resp.StatusCode)
  }

  return reader(resp.Body)
}

func (HttpFetcher) ReadBody(r io.Reader) ([]byte, error) {
  body, err := io.ReadAll(r)
  if err != nil {
    return body, err
  }
  return body, nil
}

func (HttpFetcher) ParseJson(body []byte) (ReleasesInfo, error) {
  var info ReleasesInfo
  if err := json.Unmarshal(body, &info); err != nil {
    return info, fmt.Errorf("Could not parse Github API reponse. Has the API changed?")
  }
  return info, nil
}
