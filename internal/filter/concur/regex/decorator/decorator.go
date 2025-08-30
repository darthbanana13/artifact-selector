package decorator

import "github.com/darthbanana13/artifact-selector/internal/filter/concur/regex"

type Constructor func(string, string, bool, bool, bool) (regex.IRegex, error)

type NilRegexDecoratorErr error
