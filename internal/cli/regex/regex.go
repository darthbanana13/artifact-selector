package regex

import (
	"fmt"
	"strings"

	"github.com/darthbanana13/artifact-selector/internal/filter/concur"
	"github.com/darthbanana13/artifact-selector/internal/filter/concur/regex/builder"
	"github.com/darthbanana13/artifact-selector/internal/log"
)

func ProcessRegexParams(regexes []string, toLowers, filters, metaKeys []string, l log.ILogger) ([]concur.FilterFunc, error) {
	regexFilters := make([]concur.FilterFunc, 0, len(regexes))
	regexBuilder := builder.
		NewRegexBuilder().
		WithLogger(l)
	for i, expr := range regexes {
		regexFilter, err := regexBuilder.
			WithExpr(expr).
			WithToLower(GetToLower(toLowers, i)).
			WithLoggerName(GetMetaKey(metaKeys, i)).
			WithMetaKey(GetMetaKey(metaKeys, i)).
			WithFilter(GetFilter(filters, i)).
			WithExclude(GetExclude(filters, i)).
			Build()
		if err != nil {
			return nil, err
		}
		regexFilters = append(regexFilters, regexFilter)
	}
	return regexFilters, nil
}

func GetToLower(toLowers []string, index int) bool {
	if index < len(toLowers) {
		return strings.EqualFold(toLowers[index], "yes") || strings.EqualFold(toLowers[index], "y")
	}
	return false
}

func GetMetaKey(metaKeys []string, index int) string {
	if index < len(metaKeys) {
		return metaKeys[index]
	}
	return fmt.Sprintf("UserRegex%v", index)
}

func GetFilter(filters []string, index int) bool {
	if index < len(filters) {
		return strings.EqualFold(filters[index], "yes") ||
			strings.EqualFold(filters[index], "y") ||
			strings.EqualFold(filters[index], "exclude") ||
			strings.EqualFold(filters[index], "e")
	}
	return false
}

func GetExclude(filters []string, index int) bool {
	if index < len(filters) {
		return strings.EqualFold(filters[index], "exclude") ||
			strings.EqualFold(filters[index], "e")
	}
	return false
}
