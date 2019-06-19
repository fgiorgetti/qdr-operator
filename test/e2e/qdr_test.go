package e2e

import (
	goctx "context"
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/utils"
	"strings"
	"testing"
	"time"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	apis "github.com/interconnectedcloud/qdr-operator/pkg/apis"
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"k8s.io/apimachinery/pkg/types"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 600
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestInterconnect(t *testing.T) {
	interconnectList := &v1alpha1.InterconnectList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Interconnect",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, interconnectList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("interconnect-group", func(t *testing.T) {
		t.Run("Mesh", InterconnectCluster)
	})
}

func interconnectScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create interconnect customer resource
	exampleInterconnect := &v1alpha1.Interconnect{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Interconnect",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-interconnect",
			Namespace: namespace,
		},
		Spec: v1alpha1.InterconnectSpec{
			DeploymentPlan: v1alpha1.DeploymentPlanType{
				Size:      3,
				Image:     "quay.io/interconnectedcloud/qdrouterd:1.6.0",
				Role:      "interior",
				Placement: "Any",
			},
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleInterconnect, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-interconnect to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-interconnect", 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-interconnect", Namespace: namespace}, exampleInterconnect)
	if err != nil {
		return err
	}
	exampleInterconnect.Spec.DeploymentPlan.Size = 4
	err = f.Client.Update(goctx.TODO(), exampleInterconnect)
	if err != nil {
		return err
	}

	// wait for example-interconnect to reach 4 replicas
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-interconnect", 4, retryInterval, timeout)
}

// Validate generated roles
func interconnectValidateRoles(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Errorf("Could not get namespace: %v", err)
	}

	cliListOptions := client.ListOptions{
		Namespace: namespace,
	}

	roleList := rbacv1.RoleList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: rbacv1.SchemeGroupVersion.String(),
		},
	}

	err = f.Client.List(goctx.TODO(), &cliListOptions, &roleList)
	if err != nil {
		return err
	}

	expectedRoles := []string{"example-interconnect", "qdr-operator"}

	// Retrieving all roles found
	var rolesFound []string
	if len(roleList.Items) > 0 {
		for _, role := range roleList.Items {
			rolesFound = append(rolesFound, role.Name)
		}
	}

	// If roles found do not match expected
	if !utils.ContainsAll(utils.FromStrings(expectedRoles), utils.FromStrings(rolesFound)) {
		t.Error("Expected", expectedRoles, "Found", rolesFound)
	}

	return nil
}

// Validate generated role bindings
func interconnectValidateRoleBindings(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Errorf("Could not get namespace: %v", err)
	}

	cliListOptions := client.ListOptions{
		Namespace: namespace,
	}

	roleBindingList := rbacv1.RoleBindingList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: rbacv1.SchemeGroupVersion.String(),
		},
	}

	err = f.Client.List(goctx.TODO(), &cliListOptions, &roleBindingList)
	if err != nil {
		return err
	}

	expectedRb := []string{
		"example-interconnect",
		"qdr-operator",
	}
	rbFound := []string{}
	if len(roleBindingList.Items) > 0 {
		for _, rb := range roleBindingList.Items {
			rbFound = append(rbFound, rb.Name)
		}
	}

	if !utils.ContainsAll(utils.FromStrings(expectedRb), utils.FromStrings(rbFound)) {
		t.Error("Expected", expectedRb, "Found", rbFound)
	}

	return nil
}

// Validate generated service accounts
func interconnectValidateSvcAccounts(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Errorf("Could not get namespace: %v", err)
	}

	cliListOptions := client.ListOptions{
		Namespace: namespace,
	}

	svcAccountList := corev1.ServiceAccountList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: rbacv1.SchemeGroupVersion.String(),
		},
	}

	err = f.Client.List(goctx.TODO(), &cliListOptions, &svcAccountList)
	if err != nil {
		return err
	}

	expectedSvcAccounts := []string{
		"example-interconnect",
		"qdr-operator",
	}
	
	var svcAccountsFound []string
	if len(svcAccountList.Items) > 0 {
		for _, svcAccount := range svcAccountList.Items {
			svcAccountsFound = append(svcAccountsFound, svcAccount.Name)
		}
	}

	if !utils.ContainsAll(utils.FromStrings(expectedSvcAccounts), utils.FromStrings(svcAccountsFound)) {
		t.Error("Expected", expectedSvcAccounts, "Found", svcAccountsFound)
	}

	return nil
}

// Validate deployments
func interconnectValidateDeployments(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Errorf("Could not get namespace: %v", err)
	}

	cliListOptions := client.ListOptions{
		Namespace: namespace,
	}

	deploymentList := appsv1.DeploymentList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
	}

	err = f.Client.List(goctx.TODO(), &cliListOptions, &deploymentList)
	if err != nil {
		return err
	}

	if len(deploymentList.Items) == 0 {
		t.Errorf("No deployments found")
	}

	found := false
	for _, deployment := range deploymentList.Items {

		if deployment.Name != "example-interconnect" {
			continue
		}

		found = true

		// Validate replicas
		if deployment.Status.Replicas != 4 {
			t.Errorf("Invalid replica count. Expected: %v. Found: %v", 4, deployment.Status.Replicas)
		}

		curCondition := deployment.Status.Conditions[0]
		if !strings.EqualFold(string(curCondition.Type), "Available") ||
			!strings.EqualFold(string(curCondition.Status), "True") {
			t.Errorf("Expected condition/status: Available/True. Got: %v/%v", curCondition.Type, curCondition.Status)
		}
	}

	if !found {
		t.Errorf("Deployment name not found: %v", "example-interconnect")
	}

	return nil
}

// Validate pods
func interconnectValidatePods(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Errorf("Could not get namespace: %v", err)
	}

	cliListOptions := client.ListOptions{
		Namespace: namespace,
	}

	podList := corev1.PodList{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
	}

	err = f.Client.List(goctx.TODO(), &cliListOptions, &podList)
	if err != nil {
		return err
	}

	if len(podList.Items) == 0 {
		t.Errorf("No Pods found")
	}

	count := 0
	expEnvVars := []string{"APPLICATION_NAME", "QDROUTERD_CONF", "POD_COUNT"}
	for _, pod := range podList.Items {
		if strings.HasPrefix(pod.Name, "example-interconnect-") {
			count++
			if "Running" != pod.Status.Phase {
				t.Errorf("Invalid POD Status. Expected: %v. Found: %v", "Running", pod.Status.Phase)
			}

			// Validating QDROUTERD_CONF env var in containers
			var envVarsFound []string
			for _, c := range pod.Spec.Containers {
				for _, envVar := range c.Env {
					if len(envVar.Value) == 0 {
						continue
					}
					envVarsFound = append(envVarsFound, envVar.Name)
					switch envVar.Name {
					case "QDROUTERD_CONF":
						if !strings.Contains(envVar.Value, "router {") {
							t.Errorf("QDROUTERD_CONF does not define the router entity")
						}
						if !strings.Contains(envVar.Value, "listener {") {
							t.Errorf("QDROUTERD_CONF does not define any listener")
						}
					case "APPLICATION_NAME":
						if envVar.Value != "example-interconnect" {
							t.Errorf("APPLICATION_NAME does not match expected value: %v",
								"example-interconnect")
						}
					case "POD_COUNT":
						if envVar.Value != "3" {
							t.Errorf("POD_COUNT does not match expected value: %v. Found: %v", "3", envVar.Value)
						}
					}
				}
			}

			if !utils.ContainsAll(utils.FromStrings(expEnvVars), utils.FromStrings(envVarsFound)) {
				t.Errorf("Missing EnvVars in Pod. Expected: %v. Found: %v",
					expEnvVars, envVarsFound)
			}
		}

	}

	if count != 4 {
		t.Errorf("Expected pods: %d. Found: %d", 4, count)
	}

	return nil
}

// Validate services
func interconnectValidateServices(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Errorf("Could not get namespace: %v", err)
	}

	cliListOptions := client.ListOptions{
		Namespace: namespace,
	}

	svcList := corev1.ServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind: "Pod",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
	}

	err = f.Client.List(goctx.TODO(), &cliListOptions, &svcList)
	if err != nil {
		return err
	}

	if len(svcList.Items) == 0 {
		t.Errorf("No Services found")
	}

	found := false
	for _, svc := range svcList.Items {
		if svc.Name != "example-interconnect" {
			continue
		}
		found = true

		if "example-interconnect" != svc.ObjectMeta.OwnerReferences[0].Name {
			t.Errorf("Invalid service owner. Expected: %v. Got: %v",
				"example-interconnect",
				svc.ObjectMeta.OwnerReferences[0].Name)
		}

		expectedPorts := []int{5672, 8080, 55672, 45672}
		portsFound := getPorts(svc)
		if !utils.ContainsAll(utils.FromInts(expectedPorts), utils.FromInts(portsFound)) {
			t.Errorf("Expected ports not available. Expected: %v. Found: %v",
				expectedPorts,
				portsFound)
		}
	}

	if !found {
		t.Errorf("Service not found. Expected: %v.", "example-interconnect")
	}

	return nil
}

// Validate ConfigMaps
// Not using ConfigMaps at this point (certmanager probably does)
func interconnectValidateConfigMaps(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Errorf("Could not get namespace: %v", err)
	}

	cliListOptions := client.ListOptions{
		Namespace: namespace,
	}

	cfgMapList := corev1.ConfigMapList{
		TypeMeta: metav1.TypeMeta{
			Kind: "ConfigMap",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
	}

	err = f.Client.List(goctx.TODO(), &cliListOptions, &cfgMapList)
	if err != nil {
		return err
	}

	if len(cfgMapList.Items) == 0 {
		t.Errorf("No ConfigMap found")
	}

	//for _, cfgMap := range cfgMapList.Items {
	//	fmt.Println(cfgMap)
	//}

	return nil
}

func getPorts(service corev1.Service) []int {
	if  len(service.Spec.Ports) == 0 {
		return []int{}
	}
	var svcPorts []int
	for _, port := range service.Spec.Ports {
		svcPorts = append(svcPorts, int(port.Port))
	}
	return svcPorts
}


func InterconnectCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for qdr-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "qdr-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = interconnectScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}

	if err = interconnectValidateRoles(t, f, ctx); err != nil {
		t.Fatal(err)
	}

	if err = interconnectValidateRoleBindings(t, f, ctx); err != nil {
		t.Fatal(err)
	}

	if err = interconnectValidateSvcAccounts(t, f, ctx); err != nil {
		t.Fatal(err)
	}

	if err = interconnectValidateDeployments(t, f, ctx); err != nil {
		t.Fatal(err)
	}

	if err = interconnectValidatePods(t, f, ctx); err != nil {
		t.Fatal(err)
	}

	if err = interconnectValidateServices(t, f, ctx); err != nil {
		t.Fatal(err)
	}
	if err = interconnectValidateConfigMaps(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
