package server

func (s *Server) RegisterRoutes() {
	s.Echo.Static("/api/v1/captures/", s.storagePath)
	s.Echo.POST("/api/v1/kpture", s.startKpture)
	s.Echo.PUT("/api/v1/kpture/:uuid/stop", s.stopKpture)
	s.Echo.GET("/api/v1/kpture/:uuid", s.getKpture)
	s.Echo.GET("/api/v1/kpture/:uuid/download", s.downLoadKpture)
	s.Echo.GET("/api/v1/kptures", s.getKptures)
	s.Echo.GET("/api/v1/agents", s.getAgents)
}