package glogdecorate

import (
  "fmt"
	"io"
  "net/http"
  "strings"

	"github.com/darthbanana13/artifact-selector/pkg/github"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogFetcherDecorator struct {
  ft github.FetcherTemplate
  l log.ILogger
  github.IFetcher
}

func NewLogFetcherDecorator(logger log.ILogger, fetcher github.IFetcher) *LogFetcherDecorator {
  d := &LogFetcherDecorator{
    l: logger,
    IFetcher: fetcher,
  }
  d.ft = github.FetcherTemplate{IFetcher: d}

  return d
}

func (d *LogFetcherDecorator) FetchArtifacts(userRepo string) (github.ReleasesInfo, error) {
  d.l.Info(
    fmt.Sprintf("Fetching artifacts for Github Repository %v", userRepo),
  )
  info, err := d.ft.FetchArtifacts(userRepo)
  if err != nil {
    d.l.Error("Error when fetching the Github Artifacts")
  }
  return info, err
}

func (d *LogFetcherDecorator) TestValidRepoName(userRepo string) error {
  d.l.Debug(
    fmt.Sprintf("Checking if repo name \"%v\" is valid", userRepo),
  )
  err := d.IFetcher.TestValidRepoName(userRepo)
  if err != nil {
    d.l.Error(err.Error())
  }
  return err
}

func (d *LogFetcherDecorator) PrepareRequest(url string) *http.Request {
  d.l.Info(
    fmt.Sprintf("Target URL is %v", url),
  )

  req := d.IFetcher.PrepareRequest(url)

  var builder strings.Builder
  fmt.Fprintf(&builder, "Headers:")
  for name, values := range req.Header {
    for _, value := range values {
      fmt.Fprintf(&builder, "\n%v: %v", name, value)
    }
  }

  d.l.Debug(builder.String())

  return req
}

func (d *LogFetcherDecorator) GetUrlBody(req *http.Request, reader github.BodyReader) ([]byte, error) {
  d.l.Debug("Making GET request")
  resp, err := d.IFetcher.GetUrlBody(req, d.ReadBody)
  if err != nil {
    d.l.Warn(err.Error())
  }
  return resp, err
}

func (d *LogFetcherDecorator) ReadBody(r io.Reader) ([]byte, error) {
  d.l.Debug("Reading response body")
  body, err := d.IFetcher.ReadBody(r)
  if err != nil {
    d.l.Error(err.Error())
    d.l.Debug(
      fmt.Sprintf("Retrieved body:\n%v", body),
    )
  }
  return body, err
}

func (d *LogFetcherDecorator) ParseJson(body []byte) (github.ReleasesInfo, error) {
  d.l.Debug("Checking if Github json response is valid")
  info, err := d.IFetcher.ParseJson(body)
  if err != nil {
    d.l.Error(err.Error())
    d.l.Debug(string(body))
  }
  return info, err
}
