# Distributed System

A simple distributed system implementation with a master-worker architecture using gRPC for communication between nodes.

## Project Overview

This project implements a basic distributed system with two types of nodes:

1. **Master Node**: Coordinates tasks and distributes them to worker nodes
2. **Worker Node**: Executes tasks received from the master node

The system uses gRPC for communication between nodes and provides a REST API to submit tasks to the master.

### Components

- **Master Node**: 
  - Listens for gRPC connections from worker nodes
  - Provides an HTTP API for task submission
  - Distributes tasks to connected worker nodes using round-robin load balancing

- **Worker Node**: 
  - Connects to the master via gRPC
  - Receives and executes tasks (commands)
  - Reports status back to the master
  - Supports unique worker IDs and configurable master address

## Getting Started

### Prerequisites

- Go 1.18 or later
- protoc (Protocol Buffers compiler)

### Installation

```bash
# Clone the repository
git clone https://github.com/DjonatanS/go-distributed-data-system.git
cd go-distributed-data-system

# Install dependencies
go mod tidy
```

## Usage

O arquivo principal agora está em `cmd/server/main.go` e o exemplo de cluster/demo está em `example/cluster_demo/main.go`.

### Start the Master Node

```bash
# Usando argumentos por flag
go run cmd/server/main.go -type=master

# Para compatibilidade antiga (se aplicável)
go run cmd/server/main.go master
```

Isso irá iniciar:
- Um servidor gRPC na porta 50051 para conexões de workers
- Uma API REST na porta 9092 para submissão de tarefas

### Start a Worker Node

```bash
# Usando argumentos por flag com ID customizado
go run cmd/server/main.go -type=worker -id=worker1 -master=localhost:50051

# Para compatibilidade antiga (se aplicável)
go run cmd/server/main.go worker
```

### Start a Complete Cluster (Both Master and Workers)

```bash
# Inicia um master e 3 workers (padrão)
go run cmd/server/main.go -type=cluster

# Inicia um master com 5 workers
go run cmd/server/main.go -type=cluster -workers=5
```

### Exemplo de Demonstração Automática

Execute o exemplo de cluster/demo, que inicia o cluster e envia tarefas automaticamente:

```bash
go run example/cluster_demo/main.go
```

### Submit Tasks

Você pode submeter tarefas para o master usando a API REST:

```bash
curl -X POST http://localhost:9092/tasks \
  -H "Content-Type: application/json" \
  -d '{"cmd":"echo hello world"}'
```

## API Reference

### REST API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/tasks` | POST   | Submit a task to be executed by worker nodes. Request body: `{"cmd":"command_to_execute"}` |
| `/workers` | GET   | Get count of currently connected workers |

### gRPC Services

The system defines the following gRPC services:

```protobuf
service NodeService {
    rpc ReportStatus(Request) returns (Response){};
    rpc AssignTask(Request) returns (stream Response){};
}
```

## How It Works

1. **Initialization**:
   - The master node starts and listens for gRPC connections
   - Worker nodes connect to the master via gRPC and register themselves

2. **Task Assignment**:
   - Tasks are submitted via HTTP to the master's REST API
   - The master distributes tasks to workers using round-robin load balancing
   - Each task is sent to exactly one worker

3. **Command Execution**:
   - Workers receive commands and execute them using the system's shell
   - The execution output is captured and printed on the worker's console
   - Workers automatically reconnect if the connection to the master is lost

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go         # Entry point principal (CLI)
├── core/
│   ├── node.go             # Master node and service implementations
│   ├── node.pb.go          # Generated protobuf message definitions
│   ├── node_grpc.pb.go     # Generated gRPC service definitions
│   ├── node.proto          # Protocol buffer definitions
│   ├── worker_node.go      # Worker node implementation
│   ├── node_test.go        # Unit tests for the master node
│   └── worker_node_test.go # Unit tests for the worker node
├── example/
│   └── cluster_demo/
│       └── main.go         # Exemplo de uso programático do cluster
├── launcher/
│   ├── launcher.go         # Package to start both master and workers together
│   └── launcher_test.go    # Tests for the launcher package
├── go.mod
├── go.sum
└── README.md
```

## Technical Details

### Master Node Implementation

The master node:
- Creates a gRPC server for worker connections
- Maintains a thread-safe registry of connected workers
- Provides a REST API for task submission (using Gin)
- Distributes commands to workers using round-robin load balancing
- Monitors worker connections and removes disconnected workers

### Worker Node Implementation

The worker node:
- Establishes a gRPC connection to the master
- Has configurable parameters (ID, master address, reconnect delay)
- Creates a streaming connection for receiving tasks
- Executes commands using Go's `exec` package and captures output
- Handles reconnection to the master if the connection is lost
- Reports status to the master with its worker ID

### Launcher Package

The launcher package provides:
- A simple way to start both master and multiple workers from a single import
- Configurable number of workers and other settings
- Clean shutdown handling with proper context cancellation

### Communication Flow

1. Client submits a task to master's REST API
2. Master selects a worker using round-robin and sends the task via gRPC stream
3. Worker executes the command and displays the output
4. If a worker disconnects, the master removes it from the pool

## Programmatic Usage

You can also use the launcher package to programmatically start a cluster:

```go
import "github.com/DjonatanS/go-distribuited-data-system/launcher"

// Start a cluster with default configuration (1 master + 3 workers)
cluster := launcher.StartCluster()

// Or customize the configuration
config := launcher.DefaultClusterConfig()
config.WorkerCount = 5
cluster := launcher.NewCluster(config)
cluster.Start()

// Later, when finished
cluster.Stop()
```

## Configuration Reference

### Worker Configuration

```go
type WorkerConfig struct {
    ID             string        // Worker identifier
    MasterAddr     string        // Address of the master node (host:port)
    ReconnectDelay time.Duration // Delay between reconnection attempts
}
```

### Cluster Configuration

```go
type ClusterConfig struct {
    MasterPort     int           // Port for the master's gRPC server (default: 50051)
    ApiPort        int           // Port for the master's REST API (default: 9092)
    WorkerCount    int           // Number of workers to start (default: 3)
    MasterHost     string        // Host for the master node (default: "localhost")
    BaseWorkerID   string        // Base prefix for worker IDs (default: "worker")
    ReconnectDelay time.Duration // Delay between worker reconnection attempts (default: 5s)
}
```

## Testing

The system includes both unit and integration tests:

```bash
# Run all tests
go test ./...

# Run unit tests only
go test ./core -run "^Test[^Integration]"
```

## Development

### Regenerating Protocol Buffers

If you modify the `.proto` files, regenerate the Go code:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    core/node.proto
```

## Future Improvements

- More advanced task scheduling and prioritization
- Worker capabilities registration and task routing based on capabilities
- Authentication and security measures
- Better error handling and retry mechanisms
- Persistence of tasks and results
- Web-based administrative UI

## License

[License information]
