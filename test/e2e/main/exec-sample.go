package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/router-mgmt/entities"
	"os/exec"
	"time"
)

type Command struct {
	Cmd *exec.Cmd
	Timeout time.Time
}



func main() {
	// Needed as not starting the framework
	framework.TestContext.KubectlPath = "/usr/local/bin/kubectl"

	timeout := time.Duration(10 * time.Second)
	//kubectl := framework.NewKubectlCommandTimeout(timeout, "get", "pods")
	kubectl := framework.NewKubectlExecCommand("router-75cd5c66f5-tvgfl", timeout, "qdmanage", "query", "--type=connection")
	stdout, err := kubectl.Exec()
	if err != nil {
		fmt.Println("Error executing kubectl command", err)
		return
	}

	fmt.Println("Command executed successfully")
	fmt.Println("STDOUT")
	fmt.Println(stdout)

	var connections []entities.Connection
	json.Unmarshal([]byte(stdout), &connections)

	for _, c := range(connections) {
		fmt.Println("Connection:", c.Name, c.Identity)
	}
	//manual_timeout()

}

func manual_timeout() {
	c1 := exec.Command("sleep", "3")
	var stdout, stderr bytes.Buffer
	c1.Stdout, c1.Stderr = &stdout, &stderr
	if err := c1.Start(); err != nil {
		fmt.Print("Error executing command:", err)
	}
	// Create a done channel
	done := make(chan error)
	go func() {
		done <- c1.Wait()
	}()
	// Create a timeout channel
	timeout := time.After(5 * time.Second)
	select {
	case <-timeout:
		_ = c1.Process.Kill()
		fmt.Print("Process has timed out")
	case <-done:
		fmt.Printf("Stdout: %s\n", string(stdout.Bytes()))
		fmt.Printf("Stderr: %s\n", string(stderr.Bytes()))
	}
}
