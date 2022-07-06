package certs

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServerCert struct {
	cert *bytes.Buffer
	pem  *bytes.Buffer
}

func (s *ServerCert) TlSConfig() (*tls.Config, error) {
	c, err := tls.X509KeyPair(s.cert.Bytes(), s.pem.Bytes())
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{c}}, nil
}

func (a ServerCert) storeK8s(client kubernetes.Interface) error {
	logrus.Info("Generating server crt")
	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: serverCa,
		},
		Data: map[string][]byte{
			"server.crt": a.cert.Bytes(),
			"server.pem": a.pem.Bytes(),
		},
		Type: v1.ServiceAccountRootCAKey,
	}

	_, err := client.CoreV1().Secrets("kpture").Create(context.Background(), &secret, metav1.CreateOptions{})
	return err
}

func genServerCerts(authority AuthorithyCertificate) (ServerCert, error) {
	var Scerts ServerCert
	// server cert config
	cert := getserverCertificate()
	// server private key
	serverPrivKey, err := rsa.GenerateKey(cryptorand.Reader, rsaKeylen)
	if err != nil {
		return ServerCert{}, errors.WithMessage(err, "could not generate private key")
	}
	// sign the server cert
	serverCertBytes, err := x509.CreateCertificate(
		cryptorand.Reader,
		cert,
		authority.cert,
		&serverPrivKey.PublicKey, authority.key)
	if err != nil {
		return ServerCert{}, errors.WithMessage(err, "could not create x509 key")
	}
	// PEM encode the  server cert and key
	Scerts.cert = new(bytes.Buffer)
	_ = pem.Encode(Scerts.cert, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})
	Scerts.pem = new(bytes.Buffer)
	_ = pem.Encode(Scerts.pem, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverPrivKey),
	})

	return Scerts, nil
}

func getserverCertificate() *x509.Certificate {
	const serialNumberSize = 1658

	return &x509.Certificate{
		DNSNames:     []string{"*.kpture.svc", "kpture.kpture.svc"},
		SerialNumber: big.NewInt(serialNumberSize),
		Subject: pkix.Name{
			CommonName:   "kpture.kpture.svc",
			Organization: []string{"kpture.io"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
}
