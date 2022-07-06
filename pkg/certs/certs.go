package certs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CertificatesHandler struct {
	AuthCert   AuthorithyCertificate
	ServerCert ServerCert
	client     kubernetes.Interface
}

func NewCertificatesHandler(client kubernetes.Interface) *CertificatesHandler {
	return &CertificatesHandler{
		client: client,
	}
}

func (s *CertificatesHandler) ForceDeleteCerts() {
	logrus.Info("Deleting certificates")
	s.client.CoreV1().Secrets("kpture").Delete(context.Background(), rootCAsecretName, v1.DeleteOptions{})
	s.client.CoreV1().Secrets("kpture").Delete(context.Background(), serverCa, v1.DeleteOptions{})
}

func (s *CertificatesHandler) IsKubeSecretValid() (*tls.Config, error) {
	rootCert, err := s.client.CoreV1().Secrets("kpture").Get(context.Background(), rootCAsecretName, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	caPem, ok := rootCert.Data["rootca.pem"]
	if !ok {
		return nil, errors.New("invalid authority certificate")
	}

	ServerCert, err := s.client.CoreV1().Secrets("kpture").Get(context.Background(), serverCa, v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	serverPem, ok := ServerCert.Data["server.pem"]
	if !ok {
		return nil, errors.New("invalid authority certificate")
	}

	servercert, ok := ServerCert.Data["server.crt"]
	if !ok {
		return nil, errors.New("invalid authority certificate")
	}

	if !isChainValid(caPem, serverPem) {
		return nil, errors.New("server certificates and authority does not match")
	}

	c, err := tls.X509KeyPair(servercert, serverPem)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{c}}, nil
}

func isChainValid(rootPem, serverPem []byte) bool {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootPem)
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode(rootPem)
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
	}

	if _, err := cert.Verify(opts); err != nil {
		return false
	}

	return true
}

func (s *CertificatesHandler) Get() (*tls.Config, error) {
	return s.IsKubeSecretValid()
}

func (s *CertificatesHandler) Gen() (*tls.Config, error) {
	var err error

	s.AuthCert, err = genAuthorithyCerts()
	if err != nil {
		return nil, err
	}

	err = s.AuthCert.storeK8s(s.client)
	if err != nil {
		logrus.Error(err)

		return nil, err
	}

	s.ServerCert, err = genServerCerts(s.AuthCert)
	if err != nil {
		return nil, err
	}

	err = s.ServerCert.storeK8s(s.client)
	if err != nil {
		logrus.Error(err)

		return nil, err
	}

	return s.ServerCert.TlSConfig()
}
