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

func FetchArtifacts(userRepo string) error {
  if err := TestValidRepoName(userRepo); err != nil {
    return err
  }

  url := fmt.Sprintf(LatestReleases, userRepo)
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("Expected HTTP Response Code %v when fetching %v, got code %v instead.", http.StatusOK, url, resp.StatusCode)
  }

  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  var info ReleasesInfo
  if err := json.Unmarshal(body, &info); err != nil {
    return fmt.Errorf("Could not parse Github API reponse for %v. Has the API changed?", url)
  }
  minfo, _ := json.Marshal(info)
  fmt.Println(string(minfo))

  return nil
}

func TestValidRepoName(userRepo string) error {
  i := strings.Index(userRepo, "/")
  if i > 0 && i < len(userRepo) - 1 {
    return nil 
  }
  return fmt.Errorf("Given repo '%v', does not have the format 'user/repo-name'", userRepo)
}
