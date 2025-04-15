package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DjonatanS/go-distribuited-data-system/core"
	"github.com/DjonatanS/go-distribuited-data-system/launcher"
)

func main() {
	nodeTypePtr := flag.String("type", "", "Node type: master, worker, or cluster")
	workerIDPtr := flag.String("id", "", "Worker ID (only used with worker type)")
	masterAddrPtr := flag.String("master", "localhost:50051", "Master address (only used with worker type)")
	workerCountPtr := flag.Int("workers", 3, "Number of workers to start (only used with cluster type)")

	flag.Parse()

	nodeType := *nodeTypePtr
	if nodeType == "" && len(os.Args) > 1 && (os.Args[1] == "master" || os.Args[1] == "worker") {
		nodeType = os.Args[1]
	}

	if nodeType == "" {
		fmt.Println("Node type required: -type=master, -type=worker, or -type=cluster")
		flag.PrintDefaults()
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	switch nodeType {
	case "master":
		fmt.Println("Starting master node...")
		go core.GetMasterNode().Start()

		<-sigCh
		fmt.Println("Master node stopping...")

	case "worker":
		workerID := *workerIDPtr
		if workerID == "" {
			workerID = "worker-default"
		}

		config := core.WorkerConfig{
			ID:             workerID,
			MasterAddr:     *masterAddrPtr,
			ReconnectDelay: core.DefaultWorkerConfig().ReconnectDelay,
		}

		worker := core.NewWorkerNode(config)
		fmt.Printf("Starting worker node %s connecting to %s...\n", workerID, *masterAddrPtr)

		errCh := make(chan error, 1)
		go func() {
			errCh <- worker.Start()
		}()

		select {
		case <-sigCh:
			fmt.Println("Stopping worker...")
			worker.Stop()
		case err := <-errCh:
			if err != nil {
				fmt.Printf("Worker stopped with error: %v\n", err)
				os.Exit(1)
			}
		}

	case "cluster":
		config := launcher.DefaultClusterConfig()
		config.WorkerCount = *workerCountPtr

		cluster := launcher.NewCluster(config)
		cluster.Start()

		<-sigCh
		fmt.Println("Stopping cluster...")
		cluster.Stop()

	default:
		fmt.Printf("Invalid node type: %s\n", nodeType)
		flag.PrintDefaults()
		os.Exit(1)
	}
}
