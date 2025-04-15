package core

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type NodeServiceGrpcServer struct {
	UnimplementedNodeServiceServer
	CmdChannel   chan string
	workersMutex sync.RWMutex
	workers      []NodeService_AssignTaskServer
	currentIdx   int
}

func (s *NodeServiceGrpcServer) ReportStatus(ctx context.Context, request *Request) (*Response, error) {
	return &Response{Data: "ok"}, nil
}

func (s *NodeServiceGrpcServer) AssignTask(request *Request, server NodeService_AssignTaskServer) error {
	s.workersMutex.Lock()
	s.workers = append(s.workers, server)
	workerIdx := len(s.workers) - 1
	s.workersMutex.Unlock()

	<-server.Context().Done()

	s.workersMutex.Lock()
	if workerIdx == len(s.workers)-1 {
		s.workers = s.workers[:workerIdx]
	} else {
		s.workers[workerIdx] = s.workers[len(s.workers)-1]
		s.workers = s.workers[:len(s.workers)-1]
	}
	s.workersMutex.Unlock()

	return nil
}

func (s *NodeServiceGrpcServer) DistributeTask(cmd string) bool {
	s.workersMutex.RLock()
	defer s.workersMutex.RUnlock()

	if len(s.workers) == 0 {
		return false
	}

	workerCount := len(s.workers)

	for i := 0; i < workerCount; i++ {
		idx := (s.currentIdx + i) % workerCount

		if err := s.workers[idx].Send(&Response{Data: cmd}); err == nil {
			s.currentIdx = (idx + 1) % workerCount
			return true
		}
	}

	return false
}

func (s *NodeServiceGrpcServer) StartTaskDistribution() {
	for cmd := range s.CmdChannel {
		s.DistributeTask(cmd)
	}
}

var nodeServiceGrpcServer *NodeServiceGrpcServer

func GetNodeServiceGrpcServer() *NodeServiceGrpcServer {
	if nodeServiceGrpcServer == nil {
		nodeServiceGrpcServer = &NodeServiceGrpcServer{
			CmdChannel: make(chan string),
			workers:    make([]NodeService_AssignTaskServer, 0),
		}
		go nodeServiceGrpcServer.StartTaskDistribution()
	}
	return nodeServiceGrpcServer
}

type MasterNode struct {
	api     *gin.Engine
	ln      net.Listener
	svr     *grpc.Server
	nodeSvr *NodeServiceGrpcServer
}

func (n *MasterNode) Init() (err error) {
	n.ln, err = net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	n.svr = grpc.NewServer()
	n.nodeSvr = GetNodeServiceGrpcServer()
	RegisterNodeServiceServer(n.svr, n.nodeSvr)
	n.api = gin.Default()
	n.api.POST("/tasks", func(c *gin.Context) {
		var payload struct {
			Cmd string `json:"cmd"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		n.nodeSvr.CmdChannel <- payload.Cmd
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	n.api.GET("/workers", func(c *gin.Context) {
		n.nodeSvr.workersMutex.RLock()
		workerCount := len(n.nodeSvr.workers)
		n.nodeSvr.workersMutex.RUnlock()

		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"workerCount": workerCount,
		})
	})

	return nil
}

func (n *MasterNode) Start() {
	go n.svr.Serve(n.ln)
	_ = n.api.Run(":9092")
	n.svr.Stop()
}

var node *MasterNode

func GetMasterNode() *MasterNode {
	if node == nil {
		node = &MasterNode{}
		if err := node.Init(); err != nil {
			panic(err)
		}
	}
	return node
}
