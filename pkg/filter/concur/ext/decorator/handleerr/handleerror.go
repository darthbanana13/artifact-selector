package handleerr

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/pkg/filter"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext"
	"github.com/darthbanana13/artifact-selector/pkg/filter/concur/ext/decorator"
	"github.com/darthbanana13/artifact-selector/pkg/funcdecorator"
)

type HandleErrDecorator struct {
	Ext ext.IExt
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
		Ext: ext,
	}, nil
}

func (hed *HandleErrDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	return hed.Ext.FilterArtifact(artifact)
}

func (hed *HandleErrDecorator) SetTargetExts(targetExts []string) error {
	if err := CheckValidExts(targetExts); err != nil {
		return err
	}
	return hed.Ext.SetTargetExts(targetExts)
}

func CheckValidExts(exts []string) error {
	if len(exts) == 0 {
		return EmptyExtsErr(errors.New("At least one extension is required"))
	}
	return nil
}
