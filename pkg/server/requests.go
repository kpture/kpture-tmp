package server

import (
	"fmt"
	"net/http"
	"newproxy/pkg/agent"
	"newproxy/pkg/capture"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// getAgents godoc
// @Summary      Get Agents
// @Description  Get Agents
// @Produce      json
// @Tags         agents
// @Failure      500  {string}  string
// @Success      200  {object}  []agent.AgentInfo
// @Router       /api/v1/agents [get]
func (s *Server) getAgents(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "getAgents")
	resp := []agent.Info{}
	for _, agent := range s.Agents {
		resp = append(resp, agent.Info())
	}
	return c.JSON(http.StatusOK, resp)
}

type KptureRequest struct {
	Name          string        `json:"name,omitempty"`
	AgentsRequest AgentsRequest `json:"agents,omitempty"`
}
type AgentsRequest []struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// @Summary      Start Kpture
// @Description  Start Kpture
// @Accept       json
// @Produce      json
// @Tags         kptures
// @Param        data  body      KptureRequest  true  "selected agents for capture"
// @Failure      500  {string}  string
// @Success      200   {object}  capture.Kpture
// @Header       200   {string}  Websocket  ""
// @Router       /api/v1/kpture [post]
func (s *Server) startKpture(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "StartPodCapture")
	u := &KptureRequest{}

	if err := c.Bind(u); err != nil {
		s.logger.Errorf("error binding request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Println(u)

	agents := []capture.Agent{}
	for _, pod := range u.AgentsRequest {
		for _, agent := range s.Agents {
			if agent.Info().Name == pod.Name && agent.Info().Namespace == pod.Namespace {
				agents = append(agents, agent)
			}
		}
	}
	k, err := s.StartKpture(u.Name, agents)
	if err != nil {
		s.logger.Errorf("error starting capture: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, k)
}

// @Summary      Stop Kpture
// @Description  Stop Kpture
// @Produce      json
// @Tags         kptures
// @Param        uuid  path      string  true  "capture uuid"
// @Failure      500   {string}  string
// @Failure      404  {string}  string
// @Success      200   {object}  capture.Kpture
// @Router       /api/v1/kpture/{uuid}/stop [put]
func (s *Server) stopKpture(c echo.Context) (err error) {
	c.Request().Header.Set(echo.HeaderXRequestID, "StopPodCapture")
	uuid := c.Param("uuid")
	k := s.GetKpture(uuid)
	if k == nil {
		return c.JSON(http.StatusNotFound, "capture not found")
	}

	err = k.Stop()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, k)
}

// @Summary      Download Kpture
// @Description  Download Kpture
// @Produce      json
// @Tags         kptures
// @Param        uuid  path      string  true  "capture uuid"
// @Failure      500   {string}  string
// @Failure      404   {string}  string
// @Router       /api/v1/kpture/{uuid}/download [get]
func (s *Server) downLoadKpture(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "DownLoadCapture")
	uuid := c.Param("uuid")
	kpture := s.GetKpture(uuid)
	if kpture == nil {
		return c.JSON(http.StatusNotFound, "capture not found")
	}
	if kpture.Status != capture.KptureStatusTerminated {
		return c.JSON(http.StatusBadRequest, "capture not finished")
	}
	// fmt.Println(filepath.Join("captures", kpture.UUID, kpture.Name+".tar.gzip"))
	return c.Redirect(http.StatusMovedPermanently, filepath.Join("/captures", kpture.UUID, kpture.Name+".tar.gzip"))
}

// @Summary      get K8S capture
// @Description  get K8S capture
// @Produce      json
// @Tags         kptures
// @Param        uuid  path      string  true  "capture uuid"
// @Failure      500   {string}  string
// @Failure      404   {string}  string
// @Success      200   {object}  capture.Kpture
// @Router       /api/v1/kpture/{uuid} [get]
func (s *Server) getKpture(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "DownLoadCapture")
	uuid := c.Param("uuid")
	kpture := s.GetKpture(uuid)
	if kpture == nil {
		return c.JSON(http.StatusNotFound, "capture not found")
	}
	return c.JSON(http.StatusOK, kpture)
}

// @Summary      get K8S capture
// @Description  get K8S capture
// @Tags         kptures
// @Produce      json
// @Failure      500   {string}  string
// @Failure      404   {string}  string
// @Success      200  {object}  map[string]capture.Kpture
// @Router       /api/v1/kptures [get]
func (s *Server) getKptures(c echo.Context) error {
	return c.JSON(http.StatusOK, s.kptures)
}
