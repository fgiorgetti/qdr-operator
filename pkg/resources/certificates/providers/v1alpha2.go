package providers

import (
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/configs"
	cmv1alpha2 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha2"
	v1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type CertificateProviderV1alpha2 struct{}

func (c *CertificateProviderV1alpha2) GetCertificateName(certificate interface{}) string {
	return certificate.(*cmv1alpha2.Certificate).GetName()
}

func (c *CertificateProviderV1alpha2) AddToScheme(scheme *runtime.Scheme) error {
	return cmv1alpha2.AddToScheme(scheme)
}

func (c *CertificateProviderV1alpha2) GetSchemaGroupVersion() schema.GroupVersion {
	return cmv1alpha2.SchemeGroupVersion
}

func (c *CertificateProviderV1alpha2) NewIssuer() runtime.Object {
	return &cmv1alpha2.Issuer{}
}

func (c *CertificateProviderV1alpha2) NewCertificate() runtime.Object {
	return &cmv1alpha2.Certificate{}
}

func (c *CertificateProviderV1alpha2) GetCASecretName(caIssuer interface{}) string {
	return caIssuer.(*cmv1alpha2.Issuer).Spec.IssuerConfig.CA.SecretName
}

func (c *CertificateProviderV1alpha2) UpdateCASecretName(caIssuer interface{}, secret string) {
	caIssuer.(*cmv1alpha2.Issuer).Spec.IssuerConfig.CA.SecretName = secret
}

func (c *CertificateProviderV1alpha2) GetCertManagerAPIVersion() string {
	return "v1alpha2"
}

func (c *CertificateProviderV1alpha2) GetCertManagerAPIGroup() string {
	return "cert-manager.io"
}

func (c *CertificateProviderV1alpha2) GetCertManagerAPIGroupVersion() string {
	return c.GetCertManagerAPIGroup() + "/" + c.GetCertManagerAPIVersion()
}

func (c *CertificateProviderV1alpha2) NewSelfSignedIssuerForCR(m *v1alpha1.Interconnect) metav1.Object {
	issuer := &cmv1alpha2.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-selfsigned",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha2.IssuerSpec{
			IssuerConfig: cmv1alpha2.IssuerConfig{
				SelfSigned: &cmv1alpha2.SelfSignedIssuer{},
			},
		},
	}
	return issuer
}

func (c *CertificateProviderV1alpha2) NewCAIssuerForCR(m *v1alpha1.Interconnect, secret string) metav1.Object {
	issuer := &cmv1alpha2.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-ca",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha2.IssuerSpec{
			IssuerConfig: cmv1alpha2.IssuerConfig{
				CA: &cmv1alpha2.CAIssuer{
					SecretName: secret,
				},
			},
		},
	}
	return issuer
}

func (c *CertificateProviderV1alpha2) NewCAIssuer(name string, namespace string, secret string) metav1.Object {
	issuer := &cmv1alpha2.Issuer{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Issuer",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: cmv1alpha2.IssuerSpec{
			IssuerConfig: cmv1alpha2.IssuerConfig{
				CA: &cmv1alpha2.CAIssuer{
					SecretName: secret,
				},
			},
		},
	}
	return issuer
}

func (c *CertificateProviderV1alpha2) NewSelfSignedCACertificateForCR(m *v1alpha1.Interconnect) metav1.Object {
	cert := &cmv1alpha2.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-selfsigned",
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha2.CertificateSpec{
			SecretName: m.Name + "-selfsigned",
			CommonName: m.Name + "." + m.Namespace + ".svc.cluster.local",
			IsCA:       true,
			IssuerRef: v1.ObjectReference{
				Name: m.Name + "-selfsigned",
			},
		},
	}

	return cert
}

func (c *CertificateProviderV1alpha2) NewCertificateForCR(m *v1alpha1.Interconnect, profileName string, certName string, issuer string) metav1.Object {
	issuerName := issuer
	if issuer == "" {
		issuerName = m.Name + "-ca"
	}

	hostNames := configs.GetInterconnectExposedHostnames(m, profileName)
	cert := &cmv1alpha2.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      certName,
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha2.CertificateSpec{
			SecretName: certName,
			CommonName: m.Name,
			DNSNames:   hostNames,
			IssuerRef: v1.ObjectReference{
				Name: issuerName,
			},
		},
	}
	return cert
}

func (c *CertificateProviderV1alpha2) NewCACertificateForCR(m *v1alpha1.Interconnect, name string) metav1.Object {
	cert := &cmv1alpha2.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: c.GetCertManagerAPIGroupVersion(),
			Kind:       "Certificate",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: m.Namespace,
		},
		Spec: cmv1alpha2.CertificateSpec{
			SecretName: name,
			CommonName: name,
			IsCA:       true,
			IssuerRef: v1.ObjectReference{
				Name: m.Name + "-selfsigned",
			},
		},
	}
	return cert
}
