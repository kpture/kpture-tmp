package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"newproxy/pkg/capture"

	corev1 "k8s.io/api/core/v1"

	agentCapture "github.com/kpture/agent/api/capture"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary      Retrieve all pods in cluster
// @Description  Retrieve all pods in cluster
// @Tags         kubernetes
// @Accept       */*
// @Produce      json
// @Success      200  {object}  []capture.Pod
// @Failure      500  {string}  string
// @Router       /kubernetes/pods [get]
func (s *KptureServer) getAllPods(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "GetAllPods")
	list, err := s.kubeclient.CoreV1().Pods("").List(context.Background(), v1.ListOptions{
		LabelSelector: "kpture-agent=true",
	})
	if err != nil {
		s.logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	response := []capture.Pod{}
	for _, pod := range list.Items {

		agentStatus := ""
		if agentHC(pod.Status.PodIP) == nil {
			agentStatus = "healthy"
		} else {
			agentStatus = "unhealthy"
		}

		response = append(response, capture.Pod{
			AgentStatus: agentStatus,
			PodMetadata: capture.PodMetadata{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
			PodStatus: capture.PodStatus{
				Phase:   string(pod.Status.Phase),
				Started: time.Since(time.Unix(pod.Status.StartTime.Unix(), 0)).String(),
			},
		})
	}
	return c.JSON(http.StatusOK, response)
}

func agentHC(podIP string) error {
	dial, err := grpc.Dial(fmt.Sprintf("%s:%d", podIP, 10000), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	_, err = agentCapture.NewKptureClient(dial).Health(context.Background(), &agentCapture.Empty{})
	return err
}

//GenerateHostFile generate a host file with all kubernetes endpoints
//It is used for wireshark hardcoded dns resolution
func (s *KptureServer) getHostMap() (map[string]string, error) {
	var hostmap = make(map[string]string)
	pods, err := s.kubeclient.CoreV1().Pods("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, pod := range pods.Items {
		hostmap[pod.Status.PodIP] = "pod-" + pod.Name
	}

	svc, err := s.kubeclient.CoreV1().Services("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, service := range svc.Items {
		if service.Spec.Type == corev1.ServiceTypeClusterIP && service.Spec.ClusterIP != "None" {
			hostmap[service.Spec.ClusterIP] = "svc-" + service.Name
		}
	}
	return hostmap, nil
}
