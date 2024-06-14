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

type MockMessenger struct {
	mock.Mock
}

func (m *MockMessenger) Send(msg string) {
	m.Called(msg)
}

type MockReader struct {
	mock.Mock
}

func (m *MockReader) Read(r string) ([]byte, error) {
	args := m.Called(r)
	return args.Get(0).([]byte), args.Error(1)
}

func TestFetch(t *testing.T) {
	mockExecutor := new(MockCmdExecutor)
	mockReader := new(MockReader)
	mockMessenger := new(MockMessenger)

	fetcher := fetch.NewFetcher(mockExecutor, mockReader, mockMessenger)
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
		mockMessengerCall := mockMessenger.On("Send", mock.Anything).Return()

		options, err := fetcher.Fetch(ctx, defaultSources)
		assert.NoError(t, err)
		assert.Len(t, options, 291)

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockMessenger.AssertExpectations(t)

		mockExecCall.Unset()
		mockReaderCall.Unset()
		mockMessengerCall.Unset()
	})

	t.Run("Should fail to read file", func(t *testing.T) {
		mockExecCall := mockExecutor.On("Execute", ctx, nixFile).Return("../../../nix/nixos-options.nix", nil)
		mockReaderCall := mockReader.On("Read", "../../../nix/nixos-options.nix").Return(
			[]byte{}, errors.New("failed to read file"),
		)
		mockMessengerCall := mockMessenger.On("Send", mock.Anything).Return()

		_, err := fetcher.Fetch(ctx, defaultSources)
		assert.ErrorContains(t, err, "failed to read file")

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockMessenger.AssertExpectations(t)

		mockExecCall.Unset()
		mockReaderCall.Unset()
		mockMessengerCall.Unset()
	})

	t.Run("Should fail to execute cmd", func(t *testing.T) {
		mockExecCall := mockExecutor.On("Execute", ctx, nixFile).Return("", errors.New("failed to execute cmd"))
		mockMessengerCall := mockMessenger.On("Send", mock.Anything).Return()

		_, err := fetcher.Fetch(ctx, defaultSources)
		assert.ErrorContains(t, err, "failed to execute cmd")

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockMessenger.AssertExpectations(t)

		mockExecCall.Unset()
		mockMessengerCall.Unset()
	})

	t.Run("Should fail to execute home-manager cmd", func(t *testing.T) {
		mockExecCall := mockExecutor.On("Execute", ctx, nixFile).Return("", errors.New("failed to execute cmd"))
		mockMessengerCall := mockMessenger.On("Send", mock.Anything).Return().Once()
		mockMessenger.On("Send",
			`failed to get home-manager options, try to run:\n`+
				`nix-channel --add https://github.com/nix-community/home-manager/archive/master.tar.gz home-manager\n`+
				`nix-channel --update\n\n`).Return().Once()

		hmSource := entities.Sources{
			HomeManager: nixFile,
		}

		_, err := fetcher.Fetch(ctx, hmSource)
		assert.ErrorContains(t, err, "failed to execute cmd")

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockMessenger.AssertExpectations(t)

		mockExecCall.Unset()
		mockMessengerCall.Unset()
	})
}
