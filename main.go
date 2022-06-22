package main

import (
	"fmt"
	_ "newproxy/docs"
	"newproxy/pkg/logger"
	"newproxy/pkg/server"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {

	// //HTTPS Server to receive the mutation request from K8s admission webhook
	mainLog := logger.NewLogger("main")
	// httpsServer, err := mutation.NewMutationWebHookServer()
	// if err != nil {
	// 	mainLog.Errorf("error creating mutation webhook server: %v", err)
	// }
	// go httpsServer.Start()

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

	// manager, err := capture.NewCaptureManager(kubeclient, "/tmp/captures/")
	// if err != nil {
	// 	mainLog.Error(err)
	// 	return
	// }

	server := server.NewServer(kubeclient, "./outputs/")

	fmt.Println(server.GetKptures())
	err = server.RegisterK8sAgents()
	if err != nil {
		mainLog.Error(err)
		return
	}
	server.Start()
}

// a1 := testutils.NewFakeAgent("agent1")
// a1.OpenCapture("/Users/stephaneguillemot/Documents/dev/tls/testutils/test.pcap", context.Background(), nil)

// a2 := testutils.NewFakeAgent("agent2")
// a2.OpenCapture("/Users/stephaneguillemot/Documents/dev/tls/testutils/test2.pcap", context.Background(), nil)

// t := []capture.Agent{
// 	a1, a2,
// }

// kpture, err := capture.NewKpture("kpturetest", t)
// if err != nil {
// 	panic(err)
// }

// kpture.Start()
// time.Sleep(3 * time.Second)
// kpture.Stop()
