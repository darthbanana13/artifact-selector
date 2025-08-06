package log

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/github"
	"github.com/darthbanana13/artifact-selector/pkg/log"
)

type LogDecorator struct {
	Ft github.FetcherTemplate
	L  log.ILogger
	github.IFetcher
}

func NewLogFetcherDecorator(fetcher github.IFetcher, logger log.ILogger) *LogDecorator {
	ld := &LogDecorator{
		L:        logger,
		IFetcher: fetcher,
	}
	ld.Ft = github.FetcherTemplate{IFetcher: ld}

	return ld
}

func (ld *LogDecorator) FetchArtifacts(userRepo string) (github.ReleasesInfo, error) {
	ld.L.Info("fetching artifacts for Github repository", "user repo", userRepo)
	info, err := ld.Ft.FetchArtifacts(userRepo)
	if err != nil {
		ld.L.Info("error when fetching the Github artifacts")
	}
	return info, err
}

func (ld *LogDecorator) MakeUrl(userRepo string, endpoint string) (string, error) {
	url, err := ld.IFetcher.MakeUrl(userRepo, endpoint)
	if err != nil {
		ld.L.Info(err.Error())
	}
	ld.L.Debug("artifact location", "URL", url)
	return url, err
}

func PanicHandler(l log.ILogger) {
	if r := recover(); r != nil {
		l.Panic("panicked", r)
	}
}

func (ld *LogDecorator) PrepareRequest(url string) *http.Request {
	firstValue := true
	firstHeader := true
	req := ld.IFetcher.PrepareRequest(url)
	defer PanicHandler(ld.L)

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

	ld.L.Debug("added request headers", "headers", builder.String())
	return req
}

func (ld *LogDecorator) DoGetRequest(req *http.Request) (*http.Response, error) {
	resp, err := ld.IFetcher.DoGetRequest(req)
	if err != nil {
		ld.L.Info(err.Error(), "HTTP Code", resp.StatusCode)
	}
	return resp, err
}

func (ld *LogDecorator) ReadResponseBody(r io.Reader) ([]byte, error) {
	body, err := ld.IFetcher.ReadResponseBody(r)
	if err != nil {
		ld.L.Info(err.Error(), "body", body)
	} else {
		ld.L.Debug("got response", "body", body)
	}
	return body, err
}

func (ld *LogDecorator) ParseJson(body []byte) (github.ReleasesInfo, error) {
	ld.L.Debug("checking if the Github json response is valid")
	info, err := ld.IFetcher.ParseJson(body)
	if err != nil {
		ld.L.Info(err.Error())
	}
	ld.L.Debug("parsed json", "json body", info)
	return info, err
}
