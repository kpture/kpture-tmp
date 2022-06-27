package server

func (s *Server) RegisterRoutes() {
	s.Echo.Static("/captures/", s.storagePath)
	s.Echo.POST("/api/v1/kpture", s.startKpture)
	s.Echo.PUT("/api/v1/kpture/:uuid/stop", s.stopKpture)
	s.Echo.GET("/api/v1/kpture/:uuid", s.getKpture)
	s.Echo.GET("/api/v1/kpture/ws/:uuid", s.kptureWebSocket)

	s.Echo.GET("/api/v1/kpture/:uuid/download", s.downLoadKpture)
	s.Echo.GET("/api/v1/kptures", s.getKptures)

	s.Echo.POST("/api/v1/kpture/k8s/namespace", s.startNamespacedKpture)
	s.Echo.GET("/api/v1/kpture/k8s/namespaces", s.getKubernetesEnabledNs)

	s.Echo.GET("/api/v1/agents", s.getAgents)

	s.Echo.GET("/api/v1/wireshark/hostfile", s.getHostFile)

	s.Echo.GET("/api/v1/profiles", s.getProfiles)
	s.Echo.POST("/api/v1/profile/:profileName", s.createProfile)
}
