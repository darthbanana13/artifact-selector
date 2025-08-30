package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// A mock for the http.Client
type HttpMockClient struct {
	mock.Mock
}

func (c *HttpMockClient) Do(req *http.Request) (*http.Response, error) {
	args := c.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// A mock for the io.ReadCloser interface
type ReadCloserMock struct {
	mock.Mock
	io.Reader
}

func (rc *ReadCloserMock) Read(p []byte) (n int, err error) {
	return rc.Reader.Read(p)
}

func (rc *ReadCloserMock) Close() error {
	args := rc.Called()
	return args.Error(0)
}

func githubReadCloserResponse() *ReadCloserMock {
	data, err := os.ReadFile("../../testdata/github/neovim.json")
	if err != nil {
		fmt.Println("Error opening neovim.json")
		panic(err)
	}

	mockBody := &ReadCloserMock{
		Reader: bytes.NewReader(data),
	}

	return mockBody
}

func expectedNeovimReleasesInfo() ReleasesInfo {
	data, err := os.ReadFile("../../testdata/github/releaseInfoNeovim.json")
	if err != nil {
		fmt.Println("Error opening releaseInfoNeovim.json")
		panic(err)
	}
	var info ReleasesInfo
	json.Unmarshal(data, &info)
	return info
}

func TestRepoName(t *testing.T) {
	var err error
	mockClient := &HttpMockClient{}
	f := NewGithubFetcher(mockClient)

	_, err = f.FetchArtifacts("neovim")
	assert.Error(t, err)

	_, err = f.FetchArtifacts("/neovim")
	assert.Error(t, err)

	_, err = f.FetchArtifacts("neovim/")
	assert.Error(t, err)

	_, err = f.FetchArtifacts("neovim/neovim/neovim")
	assert.Error(t, err)
}

func TestBadRespCode(t *testing.T) {
	mockClient, mockBody := &HttpMockClient{}, &ReadCloserMock{}

	mockBody.On("Close").Return(nil).Once()

	resp := &http.Response{StatusCode: http.StatusCreated, Body: mockBody}
	mockClient.On("Do", mock.Anything).Return(resp, nil).Once()

	f := NewGithubFetcher(mockClient)
	_, gotErr := f.FetchArtifacts("neovim/neovim")

	assert.Error(t, gotErr)
	mockClient.AssertExpectations(t)
	mockBody.AssertExpectations(t)
}

func TestBadReq(t *testing.T) {
	mockClient, mockBody := &HttpMockClient{}, &ReadCloserMock{}

	resp := &http.Response{StatusCode: http.StatusOK, Body: mockBody}
	err := fmt.Errorf("The request issued had a problem")
	mockClient.On("Do", mock.Anything).Return(resp, err).Once()

	f := NewGithubFetcher(mockClient)
	_, gotErr := f.FetchArtifacts("neovim/neovim")

	assert.Error(t, gotErr)
	mockClient.AssertExpectations(t)
	mockBody.AssertNotCalled(t, "Close")
}

func TestJsonParsing(t *testing.T) {
	unexpectedBody := "Artifact list"
	mockBody := &ReadCloserMock{
		Reader: bytes.NewReader([]byte(unexpectedBody)),
	}
	mockBody.On("Close").Return(nil).Once()
	resp := &http.Response{StatusCode: http.StatusOK, Body: mockBody}

	mockClient := &HttpMockClient{}
	mockClient.On("Do", mock.Anything).Return(resp, nil).Once()

	f := NewGithubFetcher(mockClient)
	_, gotErr := f.FetchArtifacts("neovim/neovim")

	assert.Error(t, gotErr)
	mockBody.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestFetchArtifacts(t *testing.T) {
	mockBody := githubReadCloserResponse()
	mockBody.On("Close").Return(nil).Once()
	resp := &http.Response{StatusCode: http.StatusOK, Body: mockBody}

	var gotReq *http.Request
	mockClient := &HttpMockClient{}
	mockClient.
		On("Do", mock.Anything).
		Return(resp, nil).
		Run(func(args mock.Arguments) {
			gotReq = args.Get(0).(*http.Request)
		}).
		Once()

	f := NewGithubFetcher(mockClient)
	gotInfo, gotErr := f.FetchArtifacts("neovim/neovim")

	expectedInfo := expectedNeovimReleasesInfo()

	assert.Equal(t, expectedInfo, gotInfo)
	assert.NoError(t, gotErr)

	assert.Equal(t, "https://api.github.com/repos/neovim/neovim/releases/latest", gotReq.URL.String())
	assert.Equal(t, "application/vnd.github+json", strings.Join(gotReq.Header.Values("Accept"), "; "))
	assert.Equal(t, "2022-11-28", strings.Join(gotReq.Header.Values("X-GitHub-Api-Version"), "; "))

	mockBody.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}
