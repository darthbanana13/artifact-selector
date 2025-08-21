package handleerr

import "fmt"

type GithubJSONStruct struct {
	Err error
}

func (e *GithubJSONStruct) Error() string {
	return fmt.Sprintf("Could not parse the Github API response. Has the API changed? Error: %v", e.Err)
}

func (e *GithubJSONStruct) Unwrap() error {
	return e.Err
}

type InvalidGithubRepoFormat struct {
	UserRepo string
}

func (e *InvalidGithubRepoFormat) Error() string {
	return fmt.Sprintf("Invalid Github repo format: %s, expected format is 'user/repo-name'", e.UserRepo)
}

type GithubHTTPResp struct {
	Expected int
	Received int
}

func (e *GithubHTTPResp) Error() string {
	return fmt.Sprintf("Expected HTTP Response Code %v, got code %v instead", e.Expected, e.Received)
}
