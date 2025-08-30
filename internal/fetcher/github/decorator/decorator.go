package decorator

import "github.com/darthbanana13/artifact-selector/internal/fetcher"

type Constructor func(client fetcher.IHTTPClient) (fetcher.IFetcherTemplate, error)

type NilGithubDecoratorErr error
