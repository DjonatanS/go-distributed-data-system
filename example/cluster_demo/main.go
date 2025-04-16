package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/DjonatanS/go-distributed-data-system/launcher"
)

func main() {
	// Inicia o cluster (master + workers)
	cluster := launcher.StartCluster()
	defer cluster.Stop()

	time.Sleep(2 * time.Second) // Aguarda o master e workers iniciarem

	// Envia algumas tarefas para o master via API HTTP
	tasks := []string{
		"echo 'Hello from worker 1'",
		"date",
		"uname -a",
	}

	for _, cmd := range tasks {
		fmt.Printf("Enviando tarefa: %s\n", cmd)
		payload := map[string]string{"cmd": cmd}
		body, _ := json.Marshal(payload)
		resp, err := http.Post("http://localhost:9092/tasks", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("Erro ao enviar tarefa: %v\n", err)
			continue
		}
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("Resposta: %s\n", string(respBody))
	}

	fmt.Println("Aguardando execução das tarefas...")
	time.Sleep(5 * time.Second)
	fmt.Println("Finalizando demo do cluster.")
}
