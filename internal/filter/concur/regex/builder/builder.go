package builder

import (
	"github.com/darthbanana13/artifact-selector/internal/filter/concur"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/regex"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/regex/decorator"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/regex/decorator/handleerr"
	logger "github.com/darthbanana13/artifact-selector/internal/filter/concur/regex/decorator/log"
	"github.com/darthbanana13/artifact-selector/internal/funcdecorator"
	"github.com/darthbanana13/artifact-selector/internal/log"
)

type RegexBuilder struct {
	L       log.ILogger
	Lname   string
	Expr    string
	MetaKey string
	ToLower bool
	Filter  bool
	Exclude bool
}

func NewRegexBuilder() *RegexBuilder {
	return &RegexBuilder{}
}

func (rb *RegexBuilder) WithLogger(l log.ILogger) *RegexBuilder {
	rb.L = l
	return rb
}

func (rb *RegexBuilder) WithLoggerName(name string) *RegexBuilder {
	rb.Lname = name
	return rb
}

func (rb *RegexBuilder) WithExpr(expr string) *RegexBuilder {
	rb.Expr = expr
	return rb
}

func (rb *RegexBuilder) WithMetaKey(metaKey string) *RegexBuilder {
	rb.MetaKey = metaKey
	return rb
}

func (rb *RegexBuilder) WithToLower(toLower bool) *RegexBuilder {
	rb.ToLower = toLower
	return rb
}

func (rb *RegexBuilder) WithFilter(filter bool) *RegexBuilder {
	rb.Filter = filter
	return rb
}

func (rb *RegexBuilder) WithExclude(exclude bool) *RegexBuilder {
	rb.Exclude = exclude
	return rb
}

func (rb *RegexBuilder) makeDecorators() []funcdecorator.FunctionDecorator[decorator.Constructor] {
	decorators := []funcdecorator.FunctionDecorator[decorator.Constructor]{
		handleerr.HandleErrConstructorDecorator(),
	}
	if rb.L != nil {
		decorators = append(decorators, logger.LogConstructorDecorator(rb.L, rb.Lname))
	}
	return decorators
}

func (rb *RegexBuilder) constructorWithDecorators() (decorator.Constructor, error) {
	constructor := funcdecorator.DecorateFunction[decorator.Constructor](
		regex.NewRegex,
		rb.makeDecorators()...,
	)
	return constructor, nil
}

func (rb *RegexBuilder) Build() (concur.FilterFunc, error) {
	constructor, err := rb.constructorWithDecorators()
	if err != nil {
		return nil, err
	}

	regexFilter, err := constructor(rb.Expr, rb.MetaKey, rb.ToLower, rb.Filter, rb.Exclude)
	if err != nil {
		return nil, err
	}

	return regexFilter.FilterArtifact, nil
}
