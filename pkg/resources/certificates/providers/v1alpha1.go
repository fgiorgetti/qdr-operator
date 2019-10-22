package providers

import (
	cmv1alpha1 "github.com/interconnectedcloud/cert-manager/pkg/apis/certmanager/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/configs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type CertificateProviderV1alpha1 struct{}

func (c *CertificateProviderV1alpha1) GetCertificateName(certificate interface{}) string {
	return certificate.(*cmv1alpha1.Certificate).GetName()
}

func (c *CertificateProviderV1alpha1) AddToScheme(scheme *runtime.Scheme) error {
	return cmv1alpha1.AddToScheme(scheme)
}

func (c *CertificateProviderV1alpha1) GetSchemaGroupVersion() schema.GroupVersion {
	return schema.GroupVersion{Group: c.GetCertManagerAPIGroup(), Version: c.GetCertManagerAPIVersion()}
}

func (c *CertificateProviderV1alpha1) NewIssuer() runtime.Object {
	return &cmv1alpha1.Issuer{}
}

func (c *CertificateProviderV1alpha1) NewCertificate() runtime.Object {
	return &cmv1alpha1.Certificate{}
}

func (c *CertificateProviderV1alpha1) GetCASecretName(caIssuer interface{}) string {
	return caIssuer.(*cmv1alpha1.Issuer).Spec.IssuerConfig.CA.SecretName
}

func (c *CertificateProviderV1alpha1) UpdateCASecretName(caIssuer interface{}, secret string) {
	caIssuer.(*cmv1alpha1.Issuer).Spec.IssuerConfig.CA.SecretName = secret
}

func (c *CertificateProviderV1alpha1) GetCertManagerAPIVersion() string {
	return "v1alpha1"
}

func (c *CertificateProviderV1alpha1) GetCertManagerAPIGroup() string {
	return "certmanager.k8s.io"
}

func (c *CertificateProviderV1alpha1) GetCertManagerAPIGroupVersion() string {
	return c.GetCertManagerAPIGroup() + "/" + c.GetCertManagerAPIVersion()
}

func (c *CertificateProviderV1alpha1) NewSelfSignedIssuerForCR(m *v1alpha1.Interconnect) metav1.Object {
	issuer := &cmv1alpha1.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-selfsigned",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.IssuerSpec{
			IssuerConfig: cmv1alpha1.IssuerConfig{
				SelfSigned: &cmv1alpha1.SelfSignedIssuer{},
			},
		},
	}
	return issuer
}

func (c *CertificateProviderV1alpha1) NewCAIssuerForCR(m *v1alpha1.Interconnect, secret string) metav1.Object {
	issuer := &cmv1alpha1.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-ca",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.IssuerSpec{
			IssuerConfig: cmv1alpha1.IssuerConfig{
				CA: &cmv1alpha1.CAIssuer{
					SecretName: secret,
				},
			},
		},
	}
	return issuer
}

func (c *CertificateProviderV1alpha1) NewCAIssuer(name string, namespace string, secret string) metav1.Object {
	issuer := &cmv1alpha1.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: cmv1alpha1.IssuerSpec{
			IssuerConfig: cmv1alpha1.IssuerConfig{
				CA: &cmv1alpha1.CAIssuer{
					SecretName: secret,
				},
			},
		},
	}
	return issuer
}

func (c *CertificateProviderV1alpha1) NewSelfSignedCACertificateForCR(m *v1alpha1.Interconnect) metav1.Object {
	cert := &cmv1alpha1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-selfsigned",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: m.Name + "-selfsigned",
			CommonName: m.Name + "." + m.Namespace + ".svc.cluster.local",
			IsCA:       true,
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: m.Name + "-selfsigned",
			},
		},
	}

	return cert
}

func (c *CertificateProviderV1alpha1) NewCertificateForCR(m *v1alpha1.Interconnect, profileName string, certName string, issuer string) metav1.Object {
	issuerName := issuer
	if issuer == "" {
		issuerName = m.Name + "-ca"
	}

	hostNames := configs.GetInterconnectExposedHostnames(m, profileName)
	cert := &cmv1alpha1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      certName,
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: certName,
			CommonName: m.Name,
			DNSNames:   hostNames,
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: issuerName,
			},
		},
	}
	return cert
}

func (c *CertificateProviderV1alpha1) NewCACertificateForCR(m *v1alpha1.Interconnect, name string) metav1.Object {
	cert := &cmv1alpha1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha1.CertificateSpec{
			SecretName: name,
			CommonName: name,
			IsCA:       true,
			IssuerRef: cmv1alpha1.ObjectReference{
				Name: m.Name + "-selfsigned",
			},
		},
	}
	return cert
}
