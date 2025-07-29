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
	l  log.ILogger
	github.IFetcher
}

func NewLogFetcherDecorator(logger log.ILogger, fetcher github.IFetcher) *LogFetcherDecorator {
	d := &LogFetcherDecorator{
		l:        logger,
		IFetcher: fetcher,
	}
	d.ft = github.FetcherTemplate{IFetcher: d}

	return d
}

func (d *LogFetcherDecorator) FetchArtifacts(userRepo string) (github.ReleasesInfo, error) {
	d.l.Info("fetching artifacts for Github repository", "user repo", userRepo)
	info, err := d.ft.FetchArtifacts(userRepo)
	if err != nil {
		d.l.Info("error when fetching the Github artifacts")
	}
	return info, err
}

func (d *LogFetcherDecorator) MakeUrl(userRepo string, endpoint string) (string, error) {
	url, err := d.IFetcher.MakeUrl(userRepo, endpoint)
	if err != nil {
		d.l.Info(err.Error())
	}
	d.l.Debug("artifact location", "URL", url)
	return url, err
}

func PanicHandler(l log.ILogger) {
	if r := recover(); r != nil {
		l.Panic("panicked", r)
	}
}

func (d *LogFetcherDecorator) PrepareRequest(url string) *http.Request {
	firstValue := true
	firstHeader := true
	req := d.IFetcher.PrepareRequest(url)
	defer PanicHandler(d.l)

	var builder strings.Builder
	for name, values := range req.Header {
		if firstHeader {
			fmt.Fprintf(&builder, "%v: ", name)
			firstHeader = false
		} else {
			fmt.Fprintf(&builder, ", %v: ", name)
		}
		firstValue = true
		for _, value := range values {
			if firstValue {
				fmt.Fprintf(&builder, "%v", value)
				firstValue = false
			} else {
				fmt.Fprintf(&builder, ", %v", value)
			}
		}
	}

	d.l.Debug("added request headers", "headers", builder.String())
	return req
}

func (d *LogFetcherDecorator) DoGetRequest(req *http.Request) (*http.Response, error) {
	resp, err := d.IFetcher.DoGetRequest(req)
	if err != nil {
		d.l.Info(err.Error(), "HTTP Code", resp.StatusCode)
	}
	return resp, err
}

func (d *LogFetcherDecorator) ReadResponseBody(r io.Reader) ([]byte, error) {
	body, err := d.IFetcher.ReadResponseBody(r)
	if err != nil {
		d.l.Info(err.Error(), "body", body)
	} else {
		d.l.Debug("got response", "body", body)
	}
	return body, err
}

func (d *LogFetcherDecorator) ParseJson(body []byte) (github.ReleasesInfo, error) {
	d.l.Debug("checking if the Github json response is valid")
	info, err := d.IFetcher.ParseJson(body)
	if err != nil {
		d.l.Info(err.Error())
	}
	d.l.Debug("parsed json", "json body", info)
	return info, err
}
