package handler

import (
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type IFilter interface {
	Filter(<-chan github.Artifact) <-chan github.Artifact
}
