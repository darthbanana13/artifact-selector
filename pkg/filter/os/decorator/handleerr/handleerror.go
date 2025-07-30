package handleerr

import (
	"errors"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/filter/os"
	"github.com/darthbanana13/artifact-selector/pkg/filter/os/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type HandleErrDecorator struct {
	os os.IOS
}

func HandleErrConstructorDecorator() funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(ofc decorator.Constructor) decorator.Constructor {
		return func(targetOS string) (os.IOS, error) {
			targetOS = strings.ToLower(targetOS)
			if err := CheckValidOS(targetOS); err != nil {
				return nil, err
			}
			of, err := ofc(targetOS)
			if err != nil {
				return of, err
			}
			return NewHandleErrDecorator(of)
		}
	}
}

func NewHandleErrDecorator(os os.IOS) (os.IOS, error) {
	if os == nil {
		return nil, decorator.NilOSDecoratorErr(errors.New("OSFilter/IOS cannot be nil"))
	}
	return &HandleErrDecorator{
		os: os,
	}, nil
}

func (hed *HandleErrDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	return hed.os.FilterArtifact(artifact)
}

func (hed *HandleErrDecorator) SetTargetOS(targetOS string) error {
	if err := CheckValidOS(targetOS); err != nil {
		return err
	}
	return hed.os.SetTargetOS(targetOS)
}

func CheckValidOS(osName string) error {
	_, osOk := os.OSMap[osName]
	_, distroOk := os.DistroMap[osName]
	if !osOk && !distroOk {
		return UnsupportedOSErr{OS: osName}
	}
	return nil
}
