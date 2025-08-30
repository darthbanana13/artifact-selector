package handleerr

import (
	"errors"

	"github.com/darthbanana13/artifact-selector/internal/filter/concur/regex"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/regex/decorator"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
)

type HandleErrDecorator struct {
	regex.IRegex
}

func HandleErrConstructorDecorator() funcdecorator.FunctionDecorator[decorator.Constructor] {
	return func(rfc decorator.Constructor) decorator.Constructor {
		return func(expr, metadataKey string, toLower, filter, exclude bool) (regex.IRegex, error) {
			if expr == "" {
				return nil, EmptyRegexErr{}
			}
			rf, err := rfc(expr, metadataKey, toLower, filter, exclude)
			if err != nil {
				return rf, err
			}
			return NewHandleErrDecorator(rf)
		}
	}
}

func NewHandleErrDecorator(rf regex.IRegex) (regex.IRegex, error) {
	if rf == nil {
		return nil, decorator.NilRegexDecoratorErr(errors.New("Regex/IRegex cannot be nil"))
	}
	return &HandleErrDecorator{
		IRegex: rf,
	}, nil
}
