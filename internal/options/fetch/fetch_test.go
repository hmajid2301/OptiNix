package fetch_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
	"gitlab.com/hmajid2301/optinix/internal/options/fetch"
)

type MockCmdExecutor struct {
	mock.Mock
}

func (m *MockCmdExecutor) Execute(ctx context.Context, path string) (string, error) {
	args := m.Called(ctx, path)
	return args.String(0), args.Error(1)
}

type MockReader struct {
	mock.Mock
}

type MockMessenger struct {
	mock.Mock
}

func (m *MockMessenger) Send(msg string) {
	m.Called(msg)
}

func (m *MockReader) Read(r string) ([]byte, error) {
	args := m.Called(r)
	return args.Get(0).([]byte), args.Error(1)
}

func TestFetch(t *testing.T) {
	mockExecutor := new(MockCmdExecutor)
	mockReader := new(MockReader)
	mockUpdater := new(MockMessenger)

	fetcher := fetch.NewFetcher(mockExecutor, mockReader, mockUpdater)
	defaultOptionsData, err := os.ReadFile("../../../testdata/nixos-options.json")
	assert.NoError(t, err)

	nixFile := "../../../nix/nixos-options.nix"
	defaultSources := entities.Sources{
		NixOS:       nixFile,
		HomeManager: nixFile,
		Darwin:      nixFile,
	}

	ctx := context.Background()
	t.Run("Should successfully fetch options", func(t *testing.T) {
		mockExecCall := mockExecutor.On("Execute", ctx, nixFile).Return("../../../nix/nixos-options.nix", nil)
		mockReaderCall := mockReader.On("Read", "../../../nix/nixos-options.nix").Return(defaultOptionsData, nil)
		mockUpdaterCall := mockUpdater.On("Send", mock.Anything).Return()

		options, err := fetcher.Fetch(ctx, defaultSources)
		assert.NoError(t, err)
		assert.Len(t, options, 291)

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockUpdater.AssertExpectations(t)

		// TODO: refactor this
		mockExecCall.Unset()
		mockReaderCall.Unset()
		mockUpdaterCall.Unset()
	})

	t.Run("Should fail to read file", func(t *testing.T) {
		mockExecCall := mockExecutor.On("Execute", ctx, nixFile).Return("../../../nix/nixos-options.nix", nil)
		mockReaderCall := mockReader.On("Read", "../../../nix/nixos-options.nix").Return(
			[]byte{}, errors.New("failed to read file"),
		)
		mockUpdaterCall := mockUpdater.On("Send", mock.Anything).Return()

		_, err := fetcher.Fetch(ctx, defaultSources)
		assert.ErrorContains(t, err, "failed to read file")

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockUpdater.AssertExpectations(t)

		mockExecCall.Unset()
		mockReaderCall.Unset()
		mockUpdaterCall.Unset()
	})

	t.Run("Should fail to execute cmd", func(t *testing.T) {
		mockExecCall := mockExecutor.On("Execute", ctx, nixFile).Return("", errors.New("failed to execute cmd"))
		mockUpdaterCall := mockUpdater.On("Send", mock.Anything).Return()

		_, err := fetcher.Fetch(ctx, defaultSources)
		assert.ErrorContains(t, err, "failed to execute cmd")

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockUpdater.AssertExpectations(t)

		mockExecCall.Unset()
		mockUpdaterCall.Unset()
	})

	// TODO: add more test cases
}
