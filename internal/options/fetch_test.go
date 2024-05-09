package options_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"gitlab.com/hmajid2301/optinix/internal/options"
)

type MockCmdExecutor struct {
	mock.Mock
}

func (m *MockCmdExecutor) Executor(expression string) (string, error) {
	args := m.Called(expression)
	return args.String(0), args.Error(1)
}

type MockReader struct {
	mock.Mock
}

func (m *MockReader) Read(r string) ([]byte, error) {
	args := m.Called(r)
	return args.Get(0).([]byte), args.Error(1)
}

type FetchTestSuite struct {
	suite.Suite
	mockExecutor *MockCmdExecutor
	mockReader   *MockReader
	mockFetcher  options.Fetcher
}

func (s *FetchTestSuite) SetupTest() {
	s.mockExecutor = new(MockCmdExecutor)
	s.mockReader = new(MockReader)
	s.mockFetcher = options.NewFetcher(s.mockExecutor, s.mockReader)
}

func TestFetchTestSuite(t *testing.T) {
	suite.Run(t, new(FetchTestSuite))
}

func (s *FetchTestSuite) TestFetch() {
	defaultOptionsData, err := os.ReadFile("../../testdata/nixos-options.json")
	s.NoError(err)

	defaultExpression, err := os.ReadFile("./nix/nixos-options.nix")
	s.NoError(err)

	expression := string(defaultExpression)
	defaultSources := options.Sources{
		NixOS:       expression,
		HomeManager: expression,
		Darwin:      expression,
	}

	s.T().Cleanup(func() {})

	ctx := context.Background()
	s.Run("Should successfully fetch options", func() {
		mockExecCall := s.mockExecutor.On("Executor", string(defaultExpression)).Return("./nix/nixos-options.nix", nil)
		mockReaderCall := s.mockReader.On("Read", "./nix/nixos-options.nix").Return(defaultOptionsData, nil)

		options, err := s.mockFetcher.Fetch(ctx, defaultSources)
		s.NoError(err)
		s.Len(options, 291)

		s.mockExecutor.AssertExpectations(s.T())
		s.mockReader.AssertExpectations(s.T())

		// TODO: refactor this
		mockExecCall.Unset()
		mockReaderCall.Unset()
	})

	s.Run("Should fail to read file", func() {
		mockExecCall := s.mockExecutor.On("Executor", string(defaultExpression)).Return("./nix/nixos-options.nix", nil)
		mockReaderCall := s.mockReader.On("Read", "./nix/nixos-options.nix").Return(
			[]byte{}, errors.New("failed to read file"),
		)

		_, err := s.mockFetcher.Fetch(ctx, defaultSources)
		s.ErrorContains(err, "failed to read file")

		s.mockExecutor.AssertExpectations(s.T())
		s.mockReader.AssertExpectations(s.T())
		mockExecCall.Unset()
		mockReaderCall.Unset()
	})

	s.Run("Should fail to execute cmd", func() {
		mockExecCall := s.mockExecutor.On("Executor", mock.Anything).Return("", errors.New("failed to execute cmd"))
		mockReaderCall := s.mockReader.On("Read", mock.Anything).Return(defaultOptionsData, nil)

		_, err := s.mockFetcher.Fetch(ctx, defaultSources)
		s.ErrorContains(err, "failed to execute cmd")

		s.mockExecutor.AssertExpectations(s.T())
		s.mockReader.AssertExpectations(s.T())
		mockExecCall.Unset()
		mockReaderCall.Unset()
	})

	// TODO: add more test cases
}
