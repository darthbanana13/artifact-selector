package handleerr

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter/ext"
	"github.com/darthbanana13/artifact-selector/pkg/filter/ext/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
	"github.com/darthbanana13/artifact-selector/pkg/github"
)

type HandleErrDecorator struct {
	ext ext.IExt
}

func HandleErrConstructorDecorator() funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(efc decorator.Constructor) decorator.Constructor {
		return func(targetExts []string) (ext.IExt, error) {
			if err := CheckValidExts(targetExts); err != nil {
				return nil, err
			}
			ef, err := efc(targetExts)
			if err != nil {
				return ef, err
			}
			return NewHandleErrDecorator(ef)
		}
	}
}

func NewHandleErrDecorator(ext ext.IExt) (ext.IExt, error) {
	if ext == nil {
		return nil, decorator.NilExtDecoratorErr(errors.New("ExtFilter/IExt cannot be nil"))
	}
	return &HandleErrDecorator{
		ext: ext,
	}, nil
}

func (hed *HandleErrDecorator) FilterArtifact(artifact github.Artifact) (github.Artifact, bool) {
	return hed.ext.FilterArtifact(artifact)
}

func (hed *HandleErrDecorator) SetTargetExts(targetExts []string) error {
	if err := CheckValidExts(targetExts); err != nil {
		return err
	}
	return hed.ext.SetTargetExts(targetExts)
}

func CheckValidExts(exts []string) error {
	for _, e := range exts {
		if _, ok := ext.ExtensionContentType[e]; !ok {
			return UnsupportedExtErr{Ext: e}
		}
	}
	return nil
}
