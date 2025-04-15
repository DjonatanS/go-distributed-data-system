package launcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/DjonatanS/go-distribuited-data-system/core"
)

type ClusterConfig struct {
	MasterPort     int
	ApiPort        int
	WorkerCount    int
	MasterHost     string
	BaseWorkerID   string
	ReconnectDelay time.Duration
}

func DefaultClusterConfig() ClusterConfig {
	return ClusterConfig{
		MasterPort:     50051,
		ApiPort:        9092,
		WorkerCount:    3,
		MasterHost:     "localhost",
		BaseWorkerID:   "worker",
		ReconnectDelay: 5 * time.Second,
	}
}

type Cluster struct {
	config   ClusterConfig
	master   *core.MasterNode
	workers  []*core.WorkerNode
	ctx      context.Context
	cancel   context.CancelFunc
	workerWg sync.WaitGroup
}

func NewCluster(config ClusterConfig) *Cluster {
	ctx, cancel := context.WithCancel(context.Background())
	return &Cluster{
		config:  config,
		ctx:     ctx,
		cancel:  cancel,
		workers: make([]*core.WorkerNode, 0, config.WorkerCount),
	}
}

func (c *Cluster) StartMaster() {
	fmt.Println("Starting master node...")

	c.master = core.GetMasterNode()

	go c.master.Start()
	fmt.Printf("Master node started on ports %d (gRPC) and %d (API)\n",
		c.config.MasterPort, c.config.ApiPort)
}

func (c *Cluster) StartWorkers() {
	fmt.Printf("Starting %d worker nodes...\n", c.config.WorkerCount)

	masterAddr := fmt.Sprintf("%s:%d", c.config.MasterHost, c.config.MasterPort)

	for i := 1; i <= c.config.WorkerCount; i++ {
		workerID := fmt.Sprintf("%s-%d", c.config.BaseWorkerID, i)

		config := core.WorkerConfig{
			ID:             workerID,
			MasterAddr:     masterAddr,
			ReconnectDelay: c.config.ReconnectDelay,
		}

		worker := core.NewWorkerNode(config)
		c.workers = append(c.workers, worker)

		c.workerWg.Add(1)
		go func(w *core.WorkerNode) {
			defer c.workerWg.Done()
			if err := w.Start(); err != nil && err != context.Canceled {
				fmt.Printf("Worker error: %v\n", err)
			}
		}(worker)
	}

	fmt.Printf("Started %d worker nodes\n", c.config.WorkerCount)
}

func (c *Cluster) Start() {
	c.StartMaster()

	time.Sleep(1 * time.Second)

	c.StartWorkers()
}

func (c *Cluster) Stop() {
	fmt.Println("Stopping cluster...")

	c.cancel()

	c.workerWg.Wait()

	fmt.Println("Cluster stopped")
}

func StartCluster() *Cluster {
	cluster := NewCluster(DefaultClusterConfig())
	cluster.Start()
	return cluster
}
