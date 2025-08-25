package log

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/fetcher"
	"github.com/darthbanana13/artifact-selector/pkg/fetcher/github/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/fetcher/github/decorator/handleerr"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	logging "github.com/darthbanana13/artifact-selector/pkg/log"
)

type Log struct {
	fetcher.IFetcherTemplate
	// Ft github.FetcherTemplate
	L logging.ILogger
}

func LogConstructorDecorator(logger logging.ILogger) funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(fc decorator.Constructor) decorator.Constructor {
		return func(client fetcher.IHTTPClient) (fetcher.IFetcherTemplate, error) {
			f, err := fc(client)
			if err != nil {
				return nil, err
			}
			return NewLogDecorator(f, logger)
		}
	}
}

func NewLogDecorator(fetcher fetcher.IFetcherTemplate, logger logging.ILogger) (fetcher.IFetcherTemplate, error) {
	if logger == nil {
		return nil, decorator.NilGithubDecoratorErr(errors.New("Logger can not be nil!"))
	}
	l := &Log{
		L:                logger,
		IFetcherTemplate: fetcher,
	}
	return l, nil
}

// TODO: This should be in a fetcher decorator
// func (ld *LogDecorator) FetchArtifacts(userRepo, version string) (fetcher.ReleaseInfo, []fetcher.Artifact, error) {
// 	ld.L.Info("fetching artifacts for Github repository", "user repo", userRepo)
// 	info, artifacts, err := ld.Ft.FetchArtifacts(userRepo, version)
// 	if err != nil {
// 		ld.L.Info("error when fetching the Github artifacts")
// 	}
// 	return info, artifacts, err
// }

const MakeURLErr = "Could not make URL for request"

func (l *Log) MakeURL(userRepo string, version string) (string, error) {
	url, err := l.IFetcherTemplate.MakeURL(userRepo, version)
	if err != nil {
		invalidFormat := &handleerr.InvalidGithubRepoFormat{}
		if errors.As(err, &invalidFormat) {
			l.L.Info(MakeURLErr, "Error", invalidFormat.Error(), "Repo format", invalidFormat.UserRepo)
			return url, invalidFormat
		}
		l.L.Info(MakeURLErr, "Error", err.Error())
	}
	l.L.Debug("Preparing to make request", "URL", url)
	return url, err
}

func PanicHandler(l logging.ILogger) {
	if r := recover(); r != nil {
		l.Panic("Panicked", r)
	}
}

// TODO: Break up/refactor into smaller functions
func (l *Log) PrepareRequest(url string) *http.Request {
	firstValue := true
	firstHeader := true
	req := l.IFetcherTemplate.PrepareRequest(url)
	defer PanicHandler(l.L)

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

	l.L.Debug("Added request headers", "Headers", builder.String())
	return req
}

const DoErr = "Error when fetching the Github artifacts"

func (l *Log) Do(req *http.Request) (*http.Response, error) {
	resp, err := l.IFetcherTemplate.Do(req)
	if err != nil {
		httpRespErr := &handleerr.GithubHTTPResp{}
		if errors.As(err, &httpRespErr) {
			l.L.Info(DoErr, "Error", httpRespErr.Error, "Expected", httpRespErr.Expected, "Received", httpRespErr.Received)
			return resp, httpRespErr
		}
		l.L.Info(DoErr, "Error", err.Error())
	}
	return resp, err
}

func (l *Log) ParseJson(r io.Reader) (fetcher.ReleaseInfo, []fetcher.Artifact, error) {
	l.L.Debug("Checking if the Github JSON response is valid")
	info, artifacts, err := l.IFetcherTemplate.ParseJson(r)
	if err != nil {
		l.L.Info("Error trying to parse the Github JSON response", "Error", err.Error())
	}
	l.L.Debug("Parsed json", "Release Info", info, "Artifacts", artifacts)
	return info, artifacts, err
}
