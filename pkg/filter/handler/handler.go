package handler

import (
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type IFilterHandler interface {
	Filter(github.ReleasesInfo) github.ReleasesInfo
	SetNext(IFilterHandler)
}
