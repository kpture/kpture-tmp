package capture

import (
	"context"
	"newproxy/pkg/logger"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Manager struct {
	//logger
	logger *logrus.Entry
	//kubeclient let us interact with K8S to fetch containers ID and nodes location
	kubeclient *kubernetes.Clientset
	//map of current capture
	kptures map[string]*kpture
	store   *store
}

func NewCaptureManager(kc *kubernetes.Clientset, basepath string) (*Manager, error) {
	c := new(Manager)
	c.logger = logger.NewLogger("Manager")
	c.kptures = make(map[string]*kpture)
	c.kubeclient = kc
	c.store = newStore(basepath)
	storeCapture, err := c.store.GetStoreKpture()
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	c.kptures = storeCapture
	return c, nil
}

func (c *Manager) GetCapture(uuid string) *kpture {
	return c.kptures[uuid]
}

func (c *Manager) GetCaptures() map[string]*kpture {
	return c.kptures
}

func (c *Manager) StopCapture(uuid string) (*kpture, error) {
	kpture := c.GetCapture(uuid)
	if kpture == nil {
		return nil, nil
	}

	kpture.UpdatePkNumber()
	kpture.StopKptures()
	kpture.Status.CaptureState = CaptureStatusWriting
	kpture.Status.Desciption = "Writing capture"

	// hostFile, err := c.daemonsetHandler.GenerateHostFile()
	// if err != nil {
	// 	c.logger.Error(err)
	// 	return nil, err
	// }
	// //Generate DNS file
	// dnsFile := AdditionalFile{
	// 	Path: filepath.Join("profiles", "hosts"),
	// 	Data: hostFile,
	// }

	go func() {
		err := c.store.Save(kpture)
		if err != nil {
			c.logger.Error(err)
			kpture.Status.CaptureState = CaptureStatusError
			kpture.Status.Desciption = err.Error()
		}
	}()

	return kpture, nil
}
func (c *Manager) StartCapture(p []PodMetadata, name string) (*kpture, error) {

	//Load a new kpture instance with the given name and the given pods
	k, err := c.newKpture(p, name)
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}
	//Dial grpc server for each capture pod
	err = k.OpenCaptureSockets()
	if err != nil {
		c.logger.Error(err)
		return nil, err
	}

	//Start listening for the grpc sockets
	k.StartListenSockets()

	return k, nil
}

func (c *Manager) newKpture(p []PodMetadata, name string) (*kpture, error) {
	k := NewKpture(name)
	c.kptures[k.UUID] = k
	//For each capture we need to fetch the containerID and the node location
	for _, pod := range p {
		c.logger.Debug("Requested capture ", pod.Name)
		k8spod, err := c.kubeclient.CoreV1().Pods(pod.Namespace).Get(context.Background(), pod.Name, v1.GetOptions{})
		if err != nil {
			c.logger.Error(err)
			return nil, err
		}
		pod := PodMetadata{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Filter:    pod.Filter,
		}

		socket, err := NewCaptureSocket(pod, k8spod.Status.PodIP)
		k.AddCapture(socket)
	}

	return k, nil
}
