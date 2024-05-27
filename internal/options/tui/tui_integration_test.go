package tui_test

import (
	"testing"
)

func TestIntegrationTUI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
