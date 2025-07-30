package separator

import (
	"fmt"
	"regexp"
)

func MakeAliasRegex(alias string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf("[_\\-\\. ]%s([_\\-\\. ]|$)", alias))
}
