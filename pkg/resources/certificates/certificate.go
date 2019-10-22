package certificates

import (
	"fmt"
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/certificates/providers"
	apiextv1b1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	sigsconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	certmgr_detected *bool
	Provider         CertificateProvider
	log              = logf.Log.WithName("certificates")
)

type CertificateProvider interface {
	AddToScheme(scheme *runtime.Scheme) error
	GetSchemaGroupVersion() schema.GroupVersion
	GetCertManagerAPIVersion() string
	GetCertManagerAPIGroup() string
	GetCertManagerAPIGroupVersion() string
	GetCertificateName(certificate interface{}) string
	GetCASecretName(caIssuer interface{}) string
	UpdateCASecretName(caIssuer interface{}, secret string)
	NewIssuer() runtime.Object
	NewCertificate() runtime.Object
	NewSelfSignedIssuerForCR(m *v1alpha1.Interconnect) metav1.Object
	NewCAIssuerForCR(m *v1alpha1.Interconnect, secret string) metav1.Object
	NewCAIssuer(name string, namespace string, secret string) metav1.Object
	NewSelfSignedCACertificateForCR(m *v1alpha1.Interconnect) metav1.Object
	NewCertificateForCR(m *v1alpha1.Interconnect, profileName string, certName string, issuer string) metav1.Object
	NewCACertificateForCR(m *v1alpha1.Interconnect, name string) metav1.Object
}

func DetectCertmgrIssuer() bool {
	// find certmanager issuer crd
	if certmgr_detected == nil {
		iscm := detectCertmgr()
		certmgr_detected = &iscm
	}
	return *certmgr_detected
}

// populateCrtmgrProviders will initialize the certmgrProviders slice
// with all valid implementations.
func GetCrtmgrProviders() []CertificateProvider {
	return []CertificateProvider{
		&providers.CertificateProviderV1alpha2{},
		&providers.CertificateProviderV1alpha1{},
	}
}

func detectCertmgr() bool {
	config, err := sigsconfig.GetConfig()
	if err != nil {
		log.Error(err, "Error getting config: %v")
		return false
	}

	// create a client set that includes crd schema
	extClient, err := apiextclientset.NewForConfig(config)
	if err != nil {
		log.Error(err, "Error getting ext client set: %v")
		return false
	}

	crd := &apiextv1b1.CustomResourceDefinition{}
	for _, certmgrProvider := range GetCrtmgrProviders() {
		crd, err = extClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get("issuers."+certmgrProvider.GetCertManagerAPIGroup(), metav1.GetOptions{})
		if err == nil {
			log.Info(fmt.Sprintf("Detected certmanager issuer crd (%s)", certmgrProvider.GetCertManagerAPIGroup()), "issuer", crd)
			Provider = certmgrProvider
			return true
		}
		log.Error(err, "Error retrieving CRDs")
	}

	return false

}
