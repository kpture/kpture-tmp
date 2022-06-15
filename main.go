package main

import (
	"newproxy/pkg/capture"
	"newproxy/pkg/httpserver"
	"newproxy/pkg/logger"
	"newproxy/pkg/mutation"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {

	//HTTPS Server to receive the mutation request from K8s admission webhook
	mainLog := logger.NewLogger("main")
	httpsServer, err := mutation.NewMutationWebHookServer()
	if err != nil {
		mainLog.Errorf("error creating mutation webhook server: %v", err)
	}
	go httpsServer.Start()

	k8sconfig, err := rest.InClusterConfig()
	if err != nil {
		mainLog.Error(err)
		return
	}

	// create the clientset
	kubeclient, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		mainLog.Error(err)
		return
	}

	manager, err := capture.NewCaptureManager(kubeclient, "/tmp/captures/")
	if err != nil {
		mainLog.Error(err)
		return
	}
	httpserver := httpserver.NewServer(kubeclient, manager, "/tmp/captures/")
	httpserver.Start()

}
