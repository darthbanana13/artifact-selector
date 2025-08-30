package log

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/darthbanana13/artifact-selector/internal/github"

	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// A mock for the gitlab.IFetcher
type MockFetcher struct {
	mock.Mock
}

func (mf *MockFetcher) FetchArtifacts(repo string) (github.ReleasesInfo, error) {
	args := mf.Called()
	return args.Get(0).(github.ReleasesInfo), args.Error(1)
}

func (mf *MockFetcher) TestValidRepoName(repo string) error {
	args := mf.Called()
	return args.Error(0)
}

func (mf *MockFetcher) MakeUrl(repo string, urlTemplate string) string {
	args := mf.Called()
	return args.String(0)
}

func (mf *MockFetcher) PrepareRequest(url string) *http.Request {
	args := mf.Called()
	return args.Get(0).(*http.Request)
}

func (mf *MockFetcher) GetUrlBody(req *http.Request, reader github.BodyReader) ([]byte, error) {
	args := mf.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (mf *MockFetcher) ReadBody(r io.Reader) ([]byte, error) {
	args := mf.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (mf *MockFetcher) ParseJson(body []byte) (github.ReleasesInfo, error) {
	args := mf.Called()
	return args.Get(0).(github.ReleasesInfo), args.Error(1)
}

// A mock for the log.ILogger interface
type MockLogger struct {
	mock.Mock
}

func (ml *MockLogger) Debug(msg string) {
	ml.Called(msg)
}

func (ml *MockLogger) Info(msg string) {
	ml.Called(msg)
}

func (ml *MockLogger) Warn(msg string) {
	ml.Called(msg)
}

func (ml *MockLogger) Error(msg string) {
	ml.Called(msg)
}

func (ml *MockLogger) Fatal(msg string) {
	ml.Called(msg)
}

func TestFetchArtifacts(t *testing.T) {
	ml, mf := &MockLogger{}, &MockFetcher{}
	mf.On("TestValidRepoName", mock.Anything).Return(nil).Once()
	mf.On("MakeUrl", mock.Anything, mock.Anything).Return("").Once()
	mf.On("PrepareRequest", mock.Anything).Return(&http.Request{}).Once()
	mf.On("GetUrlBody", mock.Anything).Return([]byte{}, nil).Once()
	mf.On("ReadBody", mock.Anything).Return([]byte{}, nil).Once()
	mf.On("ParseJson", mock.Anything).Return(github.ReleasesInfo{}, nil).Once()

	ml.On("Info", mock.Anything).Twice()
	ml.On("Debug", mock.Anything).Times(5)

	d := NewLogDecorator(ml, mf)
	_, err := d.FetchArtifacts("")
	assert.NoError(t, err)
	_, err = d.ReadBody(bytes.NewReader(nil))
	assert.NoError(t, err)

	mf.AssertExpectations(t)
	ml.AssertExpectations(t)
}

func TestErrRepoName(t *testing.T) {
	ml, mf := &MockLogger{}, &MockFetcher{}
	mf.On("TestValidRepoName", mock.Anything).Return(fmt.Errorf("error")).Once()

	ml.On("Info", mock.Anything).Once()
	ml.On("Debug", mock.Anything).Once()
	ml.On("Error", mock.Anything).Twice()

	d := NewLogDecorator(ml, mf)
	_, err := d.FetchArtifacts("")
	assert.Error(t, err)

	mf.AssertExpectations(t)
	ml.AssertExpectations(t)
}

func TestErrRequest(t *testing.T) {
	ml, mf := &MockLogger{}, &MockFetcher{}
	mf.On("TestValidRepoName", mock.Anything).Return(nil).Once()
	mf.On("MakeUrl", mock.Anything, mock.Anything).Return("").Once()
	mf.On("PrepareRequest", mock.Anything).Return(&http.Request{}).Once()
	mf.On("GetUrlBody", mock.Anything).Return([]byte{}, fmt.Errorf("error")).Once()

	ml.On("Info", mock.Anything).Twice()
	ml.On("Debug", mock.Anything).Times(3)
	ml.On("Warn", mock.Anything).Once()
	ml.On("Error", mock.Anything).Once()

	d := NewLogDecorator(ml, mf)
	_, err := d.FetchArtifacts("")
	assert.Error(t, err)

	mf.AssertExpectations(t)
	ml.AssertExpectations(t)
}

func TestErrParseJson(t *testing.T) {
	ml, mf := &MockLogger{}, &MockFetcher{}
	mf.On("TestValidRepoName", mock.Anything).Return(nil).Once()
	mf.On("MakeUrl", mock.Anything, mock.Anything).Return("").Once()
	mf.On("PrepareRequest", mock.Anything).Return(&http.Request{}).Once()
	mf.On("GetUrlBody", mock.Anything).Return([]byte{}, nil).Once()
	mf.On("ParseJson", mock.Anything).Return(github.ReleasesInfo{}, fmt.Errorf("error")).Once()

	ml.On("Info", mock.Anything).Twice()
	ml.On("Debug", mock.Anything).Times(5)
	ml.On("Error", mock.Anything).Twice()

	d := NewLogDecorator(ml, mf)
	_, err := d.FetchArtifacts("")
	assert.Error(t, err)

	mf.AssertExpectations(t)
	ml.AssertExpectations(t)
}

func TestErrReadBody(t *testing.T) {
	ml, mf := &MockLogger{}, &MockFetcher{}
	mf.On("ReadBody", mock.Anything).Return([]byte{}, fmt.Errorf("error")).Once()

	ml.On("Debug", mock.Anything).Times(2)
	ml.On("Error", mock.Anything).Once()

	d := NewLogDecorator(ml, mf)
	_, err := d.ReadBody(bytes.NewReader(nil))
	assert.Error(t, err)

	mf.AssertExpectations(t)
	ml.AssertExpectations(t)
}
