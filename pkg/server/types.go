package server

import "newproxy/pkg/capture"

const (
	agentPodSelector = "kpture-agent=true"
	agentPort        = 10000
	defaultProfile   = "default"
)

type (
	profile struct {
		Name    string
		kptures map[string]*capture.Kpture
	}

	kptureRequest struct {
		KptureName    string        `json:"kptureName,omitempty"`
		AgentsRequest agentsRequest `json:"agentsRequest,omitempty"`
	}

	kptureNamespaceRequest struct {
		KptureName      string `json:"kptureName,omitempty"`
		KptureNamespace string `json:"kptureNamespace,omitempty"`
	}

	agentsRequest []struct {
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}

	serverError struct {
		Message string `json:"message"`
	}
)
