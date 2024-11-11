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
  Filename      string  `json:"name"`
  FileType      string  `json:"content_type"`
  Size          uint64  `json:"size"`
  DownloadLink  string  `json:"browser_download_url"`
}

func FetchArtifacts(userRepo string) (ReleasesInfo, error) {
  if err := TestValidRepoName(userRepo); err != nil {
    return ReleasesInfo{}, err
  }

  body, err := GetUrlBody(fmt.Sprintf(LatestReleases, userRepo))
  if err != nil {
    return ReleasesInfo{}, err
  }

  return PraseJson(body)
}

func PraseJson(body []byte) (ReleasesInfo, error) {
  var info ReleasesInfo
  if err := json.Unmarshal(body, &info); err != nil {
    return info, fmt.Errorf("Could not parse Github API reponse. Has the API changed?")
  }
  return info, nil
}

func GetUrlBody(url string) ([]byte, error) {
  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("Expected HTTP Response Code %v when fetching %v, got code %v instead.", http.StatusOK, url, resp.StatusCode)
  }

  return FetchBody(resp.Body)
}

func FetchBody(r io.Reader) ([]byte, error) {
  body, err := io.ReadAll(r)
  if err != nil {
    return nil, err
  }
  return body, nil
}

func TestValidRepoName(userRepo string) error {
  i := strings.Index(userRepo, "/")
  if i > 0 && i < len(userRepo) - 1 {
    return nil 
  }
  return fmt.Errorf("Given repo '%v', does not have the format 'user/repo-name'", userRepo)
}
