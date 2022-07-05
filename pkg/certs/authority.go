package certs

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type AuthorithyCertificate struct {
	cert   *x509.Certificate
	CaByte []byte
	key    *rsa.PrivateKey
	Pem    *bytes.Buffer
}

func (a AuthorithyCertificate) storeK8s(client kubernetes.Interface) error {
	logrus.Info("Generating certificate authority secret")
	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: rootCAsecretName,
		},
		Data: map[string][]byte{
			"rootca.crt": a.CaByte,
			"rootca.pem": a.Pem.Bytes(),
		},
		Type: v1.ServiceAccountRootCAKey,
	}

	_, err := client.CoreV1().Secrets("kpture").Create(context.Background(), &secret, metav1.CreateOptions{})

	return err
}

func genAuthorithyCerts() (AuthorithyCertificate, error) {
	var (
		err   error
		aCert AuthorithyCertificate
	)

	aCert.cert = getCACertificate()

	aCert.key, err = rsa.GenerateKey(cryptorand.Reader, rsaKeylen)
	if err != nil {
		return aCert, errors.WithMessage(err, "could not generate rsa key")
	}

	// Self signed CA certificate
	aCert.CaByte, err = x509.CreateCertificate(cryptorand.Reader, aCert.cert, aCert.cert, &aCert.key.PublicKey, aCert.key)
	if err != nil {
		return aCert, errors.WithMessage(err, "could not generate rsa key")
	}

	// PEM encode CA cert
	aCert.Pem = new(bytes.Buffer)
	err = pem.Encode(aCert.Pem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: aCert.CaByte,
	})
	if err != nil {
		return aCert, errors.WithMessage(err, "could not generate rsa key")
	}
	return aCert, nil
}

func getCACertificate() *x509.Certificate {
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
