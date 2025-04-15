package core

import (
	"testing"
	"time"
)

func TestWorkerConfigDefaults(t *testing.T) {
	config := DefaultWorkerConfig()

	if config.ID != "worker-default" {
		t.Errorf("Expected default ID 'worker-default', got '%s'", config.ID)
	}

	if config.MasterAddr != "localhost:50051" {
		t.Errorf("Expected default master address 'localhost:50051', got '%s'", config.MasterAddr)
	}

	if config.ReconnectDelay != 5*time.Second {
		t.Errorf("Expected default reconnect delay of 5s, got %v", config.ReconnectDelay)
	}
}

func TestWorkerNodeCreation(t *testing.T) {
	config := WorkerConfig{
		ID:             "test-worker",
		MasterAddr:     "test-host:12345",
		ReconnectDelay: 10 * time.Second,
	}

	worker := NewWorkerNode(config)

	if worker.config.ID != "test-worker" {
		t.Errorf("Worker has wrong ID, expected 'test-worker', got '%s'", worker.config.ID)
	}

	if worker.config.MasterAddr != "test-host:12345" {
		t.Errorf("Worker has wrong master address, expected 'test-host:12345', got '%s'", worker.config.MasterAddr)
	}

	if worker.config.ReconnectDelay != 10*time.Second {
		t.Errorf("Worker has wrong reconnect delay, expected 10s, got %v", worker.config.ReconnectDelay)
	}

	if worker.ctx == nil {
		t.Errorf("Worker context should not be nil")
	}

	if worker.cancel == nil {
		t.Errorf("Worker cancel function should not be nil")
	}
}

func TestWorkerStop(t *testing.T) {
	worker := NewWorkerNode(DefaultWorkerConfig())

	// Ensure Stop doesn't panic when no connection exists
	worker.Stop()

	// Verify context is canceled
	if worker.ctx.Err() == nil {
		t.Errorf("Worker context should be canceled after Stop()")
	}
}
