package mutation

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const injectionPatch = `[
	{
	  "op":"add",
	  "path":"/metadata/labels/kpture-agent",
	  "value":"true"
	},
	{
	  "op":"add",
	  "path":"/spec/containers/1",
	  "value":{
		"image":"kpture/agent:latest",
		"imagePullPolicy":"IfNotPresent",
		"name":"kpture-agent",
		"ports":[
		  {
			"name":"agent",
			"containerPort":10000,
			"protocol":"TCP"
		  }
		]
	  }
	}
  ]`

func getserverCertificate() *x509.Certificate {
	const serialNumberSize = 1658

	return &x509.Certificate{
		DNSNames:     []string{"*.kpture.svc", "tls.kpture.svc"},
		SerialNumber: big.NewInt(serialNumberSize),
		Subject: pkix.Name{
			CommonName:   "tls.kpture.svc",
			Organization: []string{"kpture.io"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
}

func getCACertificate() *x509.Certificate {
	const serialNumberSize = 2020

	return &x509.Certificate{
		SerialNumber: big.NewInt(serialNumberSize),
		Subject: pkix.Name{
			Organization: []string{"kpture.io"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

const (
	NameSpaceSelectorLabel = "kpture-agent"
	NameSpaceSelectorValue = "enabled"
)

func getNamespaceLabelSelector() *metav1.LabelSelector {
	return &metav1.LabelSelector{
		MatchLabels: map[string]string{
			NameSpaceSelectorLabel: NameSpaceSelectorValue,
		},
	}
}

const (
	webHookName = "kpture-server-mutate"
)

func getMutationConfig(cacerts []byte) *admissionregistrationv1.MutatingWebhookConfiguration {
	path := "/mutate"
	fail := admissionregistrationv1.Ignore
	se := admissionregistrationv1.SideEffectClassNone

	return &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: webHookName,
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{{
			SideEffects:             &se,
			AdmissionReviewVersions: []string{"v1"},
			NamespaceSelector:       getNamespaceLabelSelector(),
			Name:                    "tls.kpture.svc",
			ClientConfig: admissionregistrationv1.WebhookClientConfig{
				CABundle: cacerts, // CA bundle created earlier
				Service: &admissionregistrationv1.ServiceReference{
					Name:      "tls",
					Namespace: "kpture",
					Path:      &path,
				},
			},
			Rules: []admissionregistrationv1.RuleWithOperations{{
				Operations: []admissionregistrationv1.OperationType{
					admissionregistrationv1.Create,
				},
				Rule: admissionregistrationv1.Rule{
					APIGroups:   []string{""},
					APIVersions: []string{"v1"},
					Resources:   []string{"pods"},
				},
			}},
			FailurePolicy: &fail,
		}},
	}
}
