package fetcher

import (
	"io"
	"net/http"
)

type (
	IFetcher interface {
		FetchArtifacts(userRepo string) (ReleaseInfo, []Artifact, error)
	}

	IURLBuilder interface {
		MakeURL(string, string) (string, error)
	}

	IRequestPreparer interface {
		PrepareRequest(url string) *http.Request
	}

	IHTTPClient interface {
		Do(*http.Request) (*http.Response, error)
	}

	IJsonParser interface {
		ParseJson(io.Reader) (ReleaseInfo, []Artifact, error)
	}

	IFetcherTemplate interface {
		IURLBuilder
		IRequestPreparer
		IHTTPClient
		IJsonParser
	}
)
