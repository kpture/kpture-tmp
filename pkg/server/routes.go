package server

import (
	"newproxy/pkg/mutation"
)

const basePath = "/kpture/api/v1/"

func (s *Server) RegisterRoutes() {
	s.Echo.Static(basePath+"captures/", s.storagePath)
	s.Echo.POST("/mutate", mutation.HandleMutate)

	s.Echo.POST(basePath+"kpture", s.startKpture)

	s.Echo.PUT(basePath+"kpture/:uuid/stop", s.stopKpture)
	s.Echo.GET(basePath+"kpture/:uuid", s.getKpture)
	s.Echo.DELETE(basePath+"kpture/:uuid", s.deleteKpture)

	s.Echo.GET(basePath+"kpture/ws/:profileName/:uuid", s.kptureWebSocket)

	s.Echo.GET(basePath+"kpture/:uuid/download", s.downLoadKpture)
	s.Echo.GET(basePath+"kptures", s.getKptures)

	s.Echo.POST(basePath+"kpture/k8s/namespace", s.startNamespacedKpture)
	s.Echo.GET(basePath+"kpture/k8s/namespaces", s.getKubernetesEnabledNs)
	s.Echo.GET(basePath+"k8s/namespaces", s.getNamespaces)
	s.Echo.POST(basePath+"k8s/namespaces/:namespace/inject", s.injectNamespace)

	s.Echo.GET(basePath+"agents", s.getAgents)

	s.Echo.GET(basePath+"wireshark/hostfile", s.getHostFile)

	s.Echo.GET(basePath+"profiles", s.getProfiles)
	s.Echo.POST(basePath+"profile/:profileName", s.createProfile)
	s.Echo.DELETE(basePath+"profile/:profileName", s.deleteProfile)
}
