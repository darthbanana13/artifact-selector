package handleerror

import (
	"errors"
	"strings"

	"github.com/darthbanana13/artifact-selector/pkg/filter/os"
	"github.com/darthbanana13/artifact-selector/pkg/filter/os/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type OSHandleDecorator struct {
	os os.IOS
}

func HandleErrorConstructorDecorator() funcdecorator.FunctionDecorator[decorator.Constructor] {
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
			return NewOSHandleDecorator(of)
		}
	}
}

func NewOSHandleDecorator(os os.IOS) (os.IOS, error) {
	if os == nil {
		return nil, decorator.NilOSDecoratorErr(errors.New("OSFilter/IOS cannot be nil"))
	}
	return &OSHandleDecorator{
		os: os,
	}, nil
}

func (ohd *OSHandleDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	return ohd.os.FilterArtifact(artifact)
}

func (ohd *OSHandleDecorator) SetTargetOS(targetOS string) error {
	if err := CheckValidOS(targetOS); err != nil {
		return err
	}
	return ohd.os.SetTargetOS(targetOS)
}

func CheckValidOS(osName string) error {
	_, osOk := os.OSMap[osName]
	_, distroOk := os.DistroMap[osName]
	if !osOk && !distroOk {
		return UnsupportedOSErr{OS: osName}
	}
	return nil
}
