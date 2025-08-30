package handleerr

type EmptyRegexErr struct{}

func (e EmptyRegexErr) Error() string {
	return "Regex pattern cannot be empty"
}
