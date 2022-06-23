package mutation

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

func genCerts() (*tls.Config, error) {
	const rsaKeylen = 4096

	var caPEM, serverCertPEM, serverPrivKeyPEM *bytes.Buffer

	ca := getCACertificate()
	// CA private key
	caPrivKey, err := rsa.GenerateKey(cryptorand.Reader, rsaKeylen)
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate rsa key")
	}
	// Self signed CA certificate
	caBytes, err := x509.CreateCertificate(cryptorand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate rsa key")
	}
	// PEM encode CA cert
	caPEM = new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	// server cert config
	cert := getserverCertificate()
	// server private key
	serverPrivKey, err := rsa.GenerateKey(cryptorand.Reader, rsaKeylen)
	if err != nil {
		return nil, errors.WithMessage(err, "could not generate private key")
	}

	// sign the server cert
	serverCertBytes, err := x509.CreateCertificate(cryptorand.Reader, cert, ca, &serverPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create x509 key")
	}

	// PEM encode the  server cert and key
	serverCertPEM = new(bytes.Buffer)
	_ = pem.Encode(serverCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCertBytes,
	})

	serverPrivKeyPEM = new(bytes.Buffer)
	_ = pem.Encode(serverPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverPrivKey),
	})

	c, err := tls.X509KeyPair(serverCertPEM.Bytes(), serverPrivKeyPEM.Bytes())
	if err != nil {
		return nil, err
	}

	err = createMutationConfig(caPEM)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{c}}, nil
}
