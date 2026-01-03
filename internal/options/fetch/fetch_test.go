package fetch_test

import (
	"context"
	"errors"
	"log/slog"
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
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	fetcher := fetch.NewFetcher(mockExecutor, mockReader, mockMessenger, logger)
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
		nixosExpr := `(builtins.getFlake (toString ./nix)).packages.${builtins.currentSystem}.nixos-options`
		hmExpr := `(builtins.getFlake (toString ./nix)).packages.${builtins.currentSystem}.home-manager-options`
		darwinExpr := `(builtins.getFlake (toString ./nix)).packages.${builtins.currentSystem}.darwin-options`

		mockExecCall1 := mockExecutor.On("Execute", ctx, nixosExpr).Return("../../../nix/nixos-options.nix", nil).Once()
		mockExecCall2 := mockExecutor.On("Execute", ctx, hmExpr).Return("../../../nix/nixos-options.nix", nil).Once()
		mockExecCall3 := mockExecutor.On("Execute", ctx, darwinExpr).Return("../../../nix/nixos-options.nix", nil).Once()
		mockReaderCall := mockReader.On("Read", "../../../nix/nixos-options.nix").Return(defaultOptionsData, nil).Times(3)
		mockMessengerCall := mockMessenger.On("Send", mock.Anything).Return()

		options, err := fetcher.Fetch(ctx, defaultSources)
		assert.NoError(t, err)
		// All three sources use the same test data, so we get 291 unique options total
		assert.Len(t, options, 291)

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockMessenger.AssertExpectations(t)

		mockExecCall1.Unset()
		mockExecCall2.Unset()
		mockExecCall3.Unset()
		mockReaderCall.Unset()
		mockMessengerCall.Unset()
	})

	t.Run("Should fail to read file", func(t *testing.T) {
		nixosExpr := `(builtins.getFlake (toString ./nix)).packages.${builtins.currentSystem}.nixos-options`

		mockExecCall := mockExecutor.On("Execute", ctx, nixosExpr).Return("../../../nix/nixos-options.nix", nil)
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
		nixosExpr := `(builtins.getFlake (toString ./nix)).packages.${builtins.currentSystem}.nixos-options`

		mockExecCall := mockExecutor.On("Execute", ctx, nixosExpr).Return("", errors.New("failed to execute cmd"))
		mockMessengerCall := mockMessenger.On("Send", mock.Anything).Return()

		_, err := fetcher.Fetch(ctx, defaultSources)
		assert.ErrorContains(t, err, "failed to execute cmd")

		mockExecutor.AssertExpectations(t)
		mockReader.AssertExpectations(t)
		mockMessenger.AssertExpectations(t)

		mockExecCall.Unset()
		mockMessengerCall.Unset()
	})
}
