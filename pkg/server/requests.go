package server

import (
	"fmt"
	"net/http"
	"newproxy/pkg/agent"
	"newproxy/pkg/capture"
	"newproxy/pkg/mutation"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getAgents godoc
// @Summary      Get Agents
// @Description  Get Agents
// @Produce      json
// @Tags         agents
// @Failure      500  {string}  string
// @Success      200  {object}  []agent.Info
// @Router       /api/v1/agents [get]
func (s *Server) getAgents(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "getAgents")

	resp := []*agent.Info{}

	if err := s.RegisterK8sAgents(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.WithMessage(err, "could not fetch k8s agents"))
	}

	for _, agent := range s.Agents {
		resp = append(resp, agent.Info())
	}

	return c.JSON(http.StatusOK, resp)
}

type KptureRequest struct {
	KptureName    string        `json:"kptureName,omitempty"`
	AgentsRequest AgentsRequest `json:"agentsRequest,omitempty"`
}

type KptureNamespaceRequest struct {
	KptureName      string `json:"kptureName,omitempty"`
	KptureNamespace string `json:"kptureNamespace,omitempty"`
}

type AgentsRequest []struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// @Summary      Start namespace kpture
// @Description  Start namespace kpture
// @Accept       json
// @Produce      json
// @Tags         kptures kubernetes
// @Param        data  body      KptureNamespaceRequest  true  "namespace for capture"
// @Failure      500  {string}  string
// @Success      200   {object}  capture.Kpture
// @Header       200   {string}  Websocket  ""
// @Router       /api/v1/kpture/k8s/namespace [post]
func (s *Server) startNamespacedKpture(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "startNamespacedKpture")

	u := &KptureNamespaceRequest{}
	if err := c.Bind(u); err != nil {
		s.logger.Errorf("error binding request: %v", err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if u.KptureName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "kpture name is empty")
	}

	agents := []capture.Agent{}

	for _, agent := range s.Agents {
		if agent.Info().Namespace == u.KptureNamespace {
			agents = append(agents, agent)
		}
	}

	if len(agents) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "no agents in namespace "+u.KptureNamespace)
	}

	k, err := s.StartKpture(u.KptureName, agents)
	if err != nil {
		s.logger.Errorf("error starting capture: %v", err)

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, k)
}

// @Summary      Get enabled kubernetes namespaces
// @Description  Get enabled kubernetes namespaces
// @Accept       json
// @Produce      json
// @Tags         kubernetes
// @Failure      500  {string}  string
// @Success      200   {object}  []string
// @Router       /api/v1/kpture/k8s/namespaces [GET]
func (s *Server) getKubernetesEnabledNs(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "getKubernetesEnabledNs")

	namespaces := []string{}
	label := fmt.Sprintf("%s=%s", mutation.NameSpaceSelectorLabel, mutation.NameSpaceSelectorValue)
	opts := v1.ListOptions{LabelSelector: label}
	nsCli := s.kubeclient.CoreV1().Namespaces()

	ns, err := nsCli.List(c.Request().Context(), opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.WithMessage(err, "could not fetch namespaces from k8s api"))
	}

	for _, currns := range ns.Items {
		namespaces = append(namespaces, currns.Name)
	}

	return c.JSON(http.StatusOK, namespaces)
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
	c.Request().Header.Set(echo.HeaderXRequestID, "startkpture")

	u := &KptureRequest{}

	if err := c.Bind(u); err != nil {
		s.logger.Errorf("error binding request: %v", err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	agents := []capture.Agent{}

	for _, pod := range u.AgentsRequest {
		for _, agent := range s.Agents {
			if agent.Info().Name == pod.Name && agent.Info().Namespace == pod.Namespace {
				agents = append(agents, agent)
			}
		}
	}

	k, err := s.StartKpture(u.KptureName, agents)
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
func (s *Server) stopKpture(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "StopPodCapture")
	uuid := c.Param("uuid")

	k := s.GetKpture(uuid)
	if k == nil {
		return c.JSON(http.StatusNotFound, "capture not found")
	}

	if err := k.Stop(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, k)
}

// @Summary      Download kpture
// @Description  Download kpture
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

	return c.Redirect(http.StatusMovedPermanently, filepath.Join("/captures", kpture.UUID, kpture.Name+".tar"))
}

// @Summary      Get kapture
// @Description  Get kapture
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

// @Summary      Get kaptures
// @Description  Get kaptures
// @Tags         kptures
// @Produce      json
// @Failure      500   {string}  string
// @Failure      404   {string}  string
// @Success      200  {object}  map[string]capture.Kpture
// @Router       /api/v1/kptures [get]
func (s *Server) getKptures(c echo.Context) error {
	return c.JSON(http.StatusOK, s.kptures)
}
