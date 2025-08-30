package handleerr

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/internal/filter"
	exterrdecorator "github.com/darthbanana13/artifact-selector/internal/filter/concur/ext/decorator/handleerr"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/extswithinsize"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/extswithinsize/decorator"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
)

type HandleErrDecorator struct {
	WithinSize extswithinsize.IWithinSize
}

func HandleErrConstructorDecorator() funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(wsc decorator.Constructor) decorator.Constructor {
		return func(maxSize uint64, percentage float64, exts []string) (extswithinsize.IWithinSize, error) {
			if err := exterrdecorator.CheckValidExts(exts); err != nil {
				return nil, err
			}
			if err := CheckValidPercentage(percentage); err != nil {
				return nil, err
			}

			ws, err := wsc(maxSize, percentage, exts)
			if err != nil {
				return ws, err
			}
			return NewHandleErrDecorator(ws)
		}
	}
}

func NewHandleErrDecorator(ws extswithinsize.IWithinSize) (extswithinsize.IWithinSize, error) {
	if ws == nil {
		return nil, decorator.NilWithinSizeDecoratorErr(errors.New("WithinSize/IWithinSize cannot be nil"))
	}
	return &HandleErrDecorator{
		WithinSize: ws,
	}, nil
}

func (hed *HandleErrDecorator) SetTargetExts(exts []string) error {
	if err := exterrdecorator.CheckValidExts(exts); err != nil {
		return err
	}
	return hed.WithinSize.SetTargetExts(exts)
}

func (hed *HandleErrDecorator) SetMaxSize(maxSize uint64) error {
	return hed.WithinSize.SetMaxSize(maxSize)
}

func (hed *HandleErrDecorator) SetPercentage(percentage float64) error {
	if err := CheckValidPercentage(percentage); err != nil {
		return err
	}
	return hed.WithinSize.SetPercentage(percentage)
}

func (hed *HandleErrDecorator) FilterArtifact(artifact filter.Artifact) (filter.Artifact, bool) {
	return hed.WithinSize.FilterArtifact(artifact)
}

func CheckValidPercentage(percentage float64) error {
	if percentage < 0 {
		return InvalidPercentageErr{Percentage: percentage}
	}
	return nil
}
