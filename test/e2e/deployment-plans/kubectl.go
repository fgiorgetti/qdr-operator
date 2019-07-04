package e2e

import (
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	router_mgmt "github.com/interconnectedcloud/qdr-operator/test/e2e/router-mgmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[kubectl] Interconnect Kubectl tests", func() {

	f := framework.NewFramework("basic-kubectl", nil)

	It("Should be able to execute commands using kubectl", func() {
		testKubectl(f)
	})

})

func testKubectl(f *framework.Framework) {

	By("Creating an edge interconnect with default size")
	ei, err := f.CreateInterconnect(f.Namespace, 0, func(ei *v1alpha1.Interconnect) {
		ei.Name = "basic-kubectl"
		ei.Spec.DeploymentPlan.Role = "edge"
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating a Deployment with 1 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "basic-kubectl", 1, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	By("Retrieving the deployment")
	dep, err := f.GetDeployment("basic-kubectl")
	Expect(err).NotTo(HaveOccurred())
	Expect(*dep.Spec.Replicas).To(Equal(int32(1)))

	By("Retriving all pods")
	pods, err := f.ListPodsForDeployment(dep)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(pods.Items)).To(Equal(1))

	By("Retrieving all pods with kubectl")
	timeout := time.Duration(10 * time.Second)
	kubectl := framework.NewKubectlCommandTimeout(timeout, "--namespace=" + f.Namespace, "get", "pods")
	stdout, err := kubectl.Exec()
	if err != nil {
		fmt.Println("Error executing kubectl command", err)
		return
	}
	fmt.Println("PODS STDOUT:", stdout)

	By("Validating connections on all pods")
	for _, pod := range pods.Items {

		conns, err := router_mgmt.QdmanageQueryConnections(f, pod.Name, nil)
		for _, c := range conns {
			fmt.Println("Connection Identity", c.Identity, "Error:", err)
		}
		//fmt.Println("Retrieving connections from POD:", pod.Name)
		//timeout := time.Duration(10 * time.Second)
		//kubectl := framework.NewKubectlExecCommand(f, pod.Name, timeout, "qdmanage", "query", "--type=connection")
		////kubectl := framework.NewKubectlCommandTimeout(timeout, "--namespace=" + f.Namespace, "exec", pod.Name, "--", "ls", "-l", "/")
		//stdout, err := kubectl.Exec()
		//if err != nil {
		//	fmt.Println("Error executing kubectl command", err)
		//	return
		//}
		//fmt.Println("LS -L / - STDOUT:", stdout)
		//
		//var connections []entities.Connection
		//_ = json.Unmarshal([]byte(stdout), &connections)
		//
		//for _, c := range(connections) {
		//	fmt.Println("Connection:", c.Name, c.Identity)
		//}
	}

}
