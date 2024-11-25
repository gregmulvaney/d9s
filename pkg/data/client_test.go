package data

import (
	"testing"
)

func TestNewDockerClient(t *testing.T) {
	docker := NewDockerClient()

	if docker.Client == nil {
		t.Fatalf("Expected Docker API Client to be initialized and got nil")
	}
}
