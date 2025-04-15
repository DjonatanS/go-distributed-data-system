package core

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerConfig struct {
	ID             string
	MasterAddr     string
	ReconnectDelay time.Duration
}

func DefaultWorkerConfig() WorkerConfig {
	return WorkerConfig{
		ID:             "worker-default",
		MasterAddr:     "localhost:50051",
		ReconnectDelay: 5 * time.Second,
	}
}

type WorkerNode struct {
	conn   *grpc.ClientConn
	c      NodeServiceClient
	config WorkerConfig
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWorkerNode(config WorkerConfig) *WorkerNode {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerNode{
		config: config,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (n *WorkerNode) Init() (err error) {
	ctx, cancel := context.WithTimeout(n.ctx, 10*time.Second)
	defer cancel()

	n.conn, err = grpc.DialContext(
		ctx,
		n.config.MasterAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to master at %s: %w", n.config.MasterAddr, err)
	}

	n.c = NewNodeServiceClient(n.conn)
	return nil
}

func (n *WorkerNode) Start() error {
	fmt.Printf("Worker %s starting, connecting to master at %s\n", n.config.ID, n.config.MasterAddr)

	for {
		if n.ctx.Err() != nil {
			return n.ctx.Err()
		}

		if err := n.Init(); err != nil {
			fmt.Printf("Worker %s: Connection error: %v. Retrying in %v...\n",
				n.config.ID, err, n.config.ReconnectDelay)
			select {
			case <-time.After(n.config.ReconnectDelay):
				continue
			case <-n.ctx.Done():
				return n.ctx.Err()
			}
		}

		fmt.Printf("Worker %s connected to master\n", n.config.ID)

		_, err := n.c.ReportStatus(n.ctx, &Request{
			Action: fmt.Sprintf("worker:%s:connected", n.config.ID),
		})
		if err != nil {
			n.conn.Close()
			continue
		}

		if err := n.processTasks(); err != nil {
			fmt.Printf("Worker %s: Task processing error: %v\n", n.config.ID, err)
		}

		n.conn.Close()
	}
}

func (n *WorkerNode) Stop() {
	n.cancel()
	if n.conn != nil {
		n.conn.Close()
	}
	fmt.Printf("Worker %s stopped\n", n.config.ID)
}

func (n *WorkerNode) processTasks() error {
	stream, err := n.c.AssignTask(n.ctx, &Request{
		Action: fmt.Sprintf("worker:%s:ready", n.config.ID),
	})
	if err != nil {
		return err
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			return err
		}

		fmt.Printf("Worker %s received command: %s\n", n.config.ID, res.Data)
		n.executeCommand(res.Data)
	}
}

func (n *WorkerNode) executeCommand(command string) {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		fmt.Printf("Worker %s: Empty command received\n", n.config.ID)
		return
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Worker %s: Command failed: %v\n", n.config.ID, err)
		return
	}

	fmt.Printf("Worker %s: Command output:\n%s\n", n.config.ID, string(output))
}

var workerNode *WorkerNode

func GetWorkerNode() *WorkerNode {
	if workerNode == nil {
		workerNode = NewWorkerNode(DefaultWorkerConfig())
	}
	return workerNode
}
