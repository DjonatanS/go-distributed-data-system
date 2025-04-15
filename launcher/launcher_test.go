package launcher

import (
	"testing"
	"time"
)

func TestClusterConfigDefaults(t *testing.T) {
	config := DefaultClusterConfig()

	if config.MasterPort != 50051 {
		t.Errorf("Expected default master port 50051, got %d", config.MasterPort)
	}

	if config.ApiPort != 9092 {
		t.Errorf("Expected default API port 9092, got %d", config.ApiPort)
	}

	if config.WorkerCount != 3 {
		t.Errorf("Expected default worker count 3, got %d", config.WorkerCount)
	}

	if config.MasterHost != "localhost" {
		t.Errorf("Expected default master host 'localhost', got '%s'", config.MasterHost)
	}

	if config.BaseWorkerID != "worker" {
		t.Errorf("Expected default base worker ID 'worker', got '%s'", config.BaseWorkerID)
	}

	if config.ReconnectDelay != 5*time.Second {
		t.Errorf("Expected default reconnect delay of 5s, got %v", config.ReconnectDelay)
	}
}

func TestClusterCreation(t *testing.T) {
	config := ClusterConfig{
		MasterPort:     60051,
		ApiPort:        8080,
		WorkerCount:    5,
		MasterHost:     "test-host",
		BaseWorkerID:   "test-worker",
		ReconnectDelay: 10 * time.Second,
	}

	cluster := NewCluster(config)

	if cluster.config.MasterPort != 60051 {
		t.Errorf("Cluster has wrong master port, expected 60051, got %d", cluster.config.MasterPort)
	}

	if cluster.config.ApiPort != 8080 {
		t.Errorf("Cluster has wrong API port, expected 8080, got %d", cluster.config.ApiPort)
	}

	if cluster.config.WorkerCount != 5 {
		t.Errorf("Cluster has wrong worker count, expected 5, got %d", cluster.config.WorkerCount)
	}

	if cluster.config.MasterHost != "test-host" {
		t.Errorf("Cluster has wrong master host, expected 'test-host', got '%s'", cluster.config.MasterHost)
	}

	if cluster.config.BaseWorkerID != "test-worker" {
		t.Errorf("Cluster has wrong base worker ID, expected 'test-worker', got '%s'", cluster.config.BaseWorkerID)
	}

	if cluster.config.ReconnectDelay != 10*time.Second {
		t.Errorf("Cluster has wrong reconnect delay, expected 10s, got %v", cluster.config.ReconnectDelay)
	}

	if cluster.ctx == nil {
		t.Errorf("Cluster context should not be nil")
	}

	if cluster.cancel == nil {
		t.Errorf("Cluster cancel function should not be nil")
	}

	if cap(cluster.workers) != 5 {
		t.Errorf("Cluster workers slice has wrong capacity, expected 5, got %d", cap(cluster.workers))
	}
}

// Skip actual cluster start/stop tests as they require running services
// In a real-world scenario, you'd use test containers or mocks for these
