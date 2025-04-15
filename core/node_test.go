package core

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNodeServiceGrpcServer_DistributeTask(t *testing.T) {
	server := &NodeServiceGrpcServer{
		CmdChannel:   make(chan string),
		workers:      make([]NodeService_AssignTaskServer, 0),
		workersMutex: sync.RWMutex{},
	}

	// No workers registered, should return false
	if server.DistributeTask("test command") {
		t.Errorf("DistributeTask should return false when no workers are available")
	}
}

func TestNodeServiceGrpcGetInstance(t *testing.T) {
	server1 := GetNodeServiceGrpcServer()
	server2 := GetNodeServiceGrpcServer()

	// Should return the same instance (singleton pattern)
	if server1 != server2 {
		t.Errorf("GetNodeServiceGrpcServer should return the same instance each time")
	}

	if server1.CmdChannel == nil {
		t.Errorf("NodeServiceGrpcServer should have a non-nil CmdChannel")
	}

	if server1.workers == nil {
		t.Errorf("NodeServiceGrpcServer should have a non-nil workers slice")
	}
}

func TestMasterNodeGetInstance(t *testing.T) {
	// This test just makes sure the GetMasterNode function doesn't panic
	// A real integration test would need a running system
	t.Run("GetInstance_NoPanic", func(t *testing.T) {
		// Use a timeout to prevent hanging if there's an issue
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		done := make(chan struct{})
		var err error

		go func() {
			defer func() {
				if r := recover(); r != nil {
					err = r.(error)
				}
				close(done)
			}()

			// This will panic if there's a port conflict or other initialization issue
			_ = GetMasterNode()
		}()

		select {
		case <-done:
			if err != nil {
				t.Errorf("GetMasterNode panicked: %v", err)
			}
		case <-ctx.Done():
			t.Errorf("GetMasterNode timed out or hung")
		}
	})
}
