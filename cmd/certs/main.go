package main

import (
	"newproxy/pkg/certs"
	"newproxy/pkg/logger"
	"newproxy/pkg/mutation"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// //HTTPS Server to receive the mutation request from K8s admission webhook
	mainLog := logger.NewLogger("main")

	k8sconfig, err := rest.InClusterConfig()
	if err != nil {
		mainLog.Error(err)

		return
	}

	kubeclient, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		mainLog.Error(err)

		return
	}
	cert := certs.NewCertificatesHandler(kubeclient)

	_, err = cert.IsKubeSecretValid()
	if err == nil {
		logrus.Info("Certificates already exist and are  valid")
		return
	}

	cert.ForceDeleteCerts()

	_, err = cert.Gen()
	if err != nil {
		logrus.Fatal(err)
	}

	err = mutation.CreateMutationConfig(cert.AuthCert.Pem.Bytes())
	if err != nil {
		logrus.Fatal(err)
	}

	return
}
