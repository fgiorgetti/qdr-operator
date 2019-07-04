package main

import (
	"bytes"
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/router-mgmt/entities"
	"os/exec"
	"reflect"
	"time"
)

type Command struct {
	Cmd *exec.Cmd
	Timeout time.Time
}

func filter(i interface{}, fn func(i interface{}) bool) []interface{} {
	s := reflect.ValueOf(i)
	if s.Kind() != reflect.Slice {
		panic("Expecting a slice")
	}

	var ret []interface{}
	//ret := make([]interface{}, s.Len())
	//ret := make([]interface{}, 1)
	ri := 0
	for j := 0; j < s.Len(); j++ {
		ii := s.Index(j).Interface()
		if fn(ii) {
			ret = append(ret, ii)
			//ret[ri] = ii
			ri++
		}
	}

	return ret
}

func main() {

	var connections []entities.Connection
	c1 := entities.Connection{ User: "C1" }
	c2 := entities.Connection{ User: "C2" }
	c3 := entities.Connection{ User: "AA" }
	connections = append(connections, c1, c2, c3)

	f1 := func (i interface{}) bool {
		return true
		//c := i.(entities.Connection)
		//if c.User[0] == 'A' {
		//	return true
		//}
		//return false
	}

	//ret := filter(connections, func(i interface{}) bool {
	//	c := i.(entities.Connection)
	//	if c.User[0] == 'B' {
	//		return true
	//	}
	//	return false
	//})
	//
	//fmt.Println(ret)
	//fmt.Printf("%T\n", ret)
	//
	//fmt.Println("ATTEMPT2")
	ret := filter(connections, f1)
	//fmt.Println(ret)
	//fmt.Printf("BANANA %T\n", ret)
	//
	var ccc []entities.Connection
	ccc = nil
	for _, v := range(ret) {
		fmt.Printf("TYPE: %T | %v\n", v, v)
		ccc = append(ccc, v.(entities.Connection))
	}

	fmt.Println("CONNECTIONS =", ccc)

	for idx, c := range(ccc) {
		fmt.Println("Connection", idx, c.User)
	}

	//// Needed as not starting the framework
	//framework.TestContext.KubectlPath = "/usr/local/bin/kubectl"
	//
	//timeout := time.Duration(10 * time.Second)
	////kubectl := framework.NewKubectlCommandTimeout(timeout, "get", "pods")
	//kubectl := framework.NewKubectlExecCommand("router-75cd5c66f5-tvgfl", timeout, "qdmanage", "query", "--type=connection")
	//stdout, err := kubectl.Exec()
	//if err != nil {
	//	fmt.Println("Error executing kubectl command", err)
	//	return
	//}
	//
	//fmt.Println("Command executed successfully")
	//fmt.Println("STDOUT")
	//fmt.Println(stdout)
	//
	//var connections []entities.Connection
	//json.Unmarshal([]byte(stdout), &connections)
	//
	//for _, c := range(connections) {
	//	fmt.Println("Connection:", c.Name, c.Identity)
	//}
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
