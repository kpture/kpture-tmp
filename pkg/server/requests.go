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
// @Failure      500  {object}  serverError
// @Success      200  {object}  []agent.Info
// @Router       /api/v1/agents [get]
func (s *Server) getAgents(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "getAgents")

	resp := []*agent.Info{}

	if err := s.RegisterK8sAgents(); err != nil {
		sErr := serverError{errors.WithMessage(err, "could not fetch k8s agents").Error()}

		return echo.NewHTTPError(http.StatusInternalServerError, sErr)
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
// @Failure      500  {object}  serverError
// @Failure      400  {object}  serverError
// @Success      200   {object}  capture.Kpture
// @Header       200   {string}  Websocket  ""
// @Router       /api/v1/kpture/k8s/namespace [post]
func (s *Server) startNamespacedKpture(c echo.Context) error {
	s.logger.Debug("startNamespacedKpture")

	u := &KptureNamespaceRequest{}
	if err := c.Bind(u); err != nil {
		s.logger.Errorf("error binding request: %v", err)
		sErr := serverError{errors.WithMessage(err, "cannot bind http request body").Error()}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	if u.KptureName == "" {
		sErr := serverError{errors.New("invalid kptureName").Error()}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	agents := []capture.Agent{}

	for _, agent := range s.Agents {
		if agent.Info().Metadata.Namespace == u.KptureNamespace {
			agents = append(agents, agent)
		}
	}

	if len(agents) == 0 {
		sErr := serverError{errors.New(fmt.Sprintf("no agent in namesapce %s", u.KptureNamespace)).Error()}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	k, err := s.StartKpture(u.KptureName, agents)
	if err != nil {
		s.logger.Errorf("error starting capture: %v", err)
		sErr := serverError{errors.WithMessage(err, "error starting kpture").Error()}

		return echo.NewHTTPError(http.StatusInternalServerError, sErr)
	}

	return c.JSON(http.StatusOK, k)
}

// @Summary      Get enabled kubernetes namespaces
// @Description  Get enabled kubernetes namespaces
// @Accept       json
// @Produce      json
// @Tags         kubernetes
// @Failure      500  {object}  serverError
// @Success      200   {object}  []string
// @Router       /api/v1/kpture/k8s/namespaces [GET]
func (s *Server) getKubernetesEnabledNs(c echo.Context) error {
	s.logger.Debug("getKubernetesEnabledNs")

	namespaces := []string{}
	label := fmt.Sprintf("%s=%s", mutation.NameSpaceSelectorLabel, mutation.NameSpaceSelectorValue)
	opts := v1.ListOptions{LabelSelector: label}
	nsCli := s.kubeclient.CoreV1().Namespaces()

	ns, err := nsCli.List(c.Request().Context(), opts)
	if err != nil {
		sErr := serverError{errors.WithMessage(err, "could not fetch kubernetes api").Error()}

		return c.JSON(http.StatusInternalServerError, sErr)
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
// @Failure      500  {object}  serverError
// @Success      200   {object}  capture.Kpture
// @Header       200   {string}  Websocket  ""
// @Router       /api/v1/kpture [post]
func (s *Server) startKpture(c echo.Context) error {
	s.logger.Debug("startKpture")

	u := &KptureRequest{}

	if err := c.Bind(u); err != nil {
		s.logger.Errorf("error binding request: %v", err)
		sErr := serverError{errors.WithMessage(err, "cannot bind request").Error()}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	agents := []capture.Agent{}

	for _, pod := range u.AgentsRequest {
		for _, agent := range s.Agents {
			if agent.Info().Metadata.Name == pod.Name && agent.Info().Metadata.Namespace == pod.Namespace {
				agents = append(agents, agent)
			}
		}
	}

	k, err := s.StartKpture(u.KptureName, agents)
	if err != nil {
		s.logger.Errorf("error starting capture: %v", err)
		sErr := serverError{errors.WithMessage(err, "error starting kpture").Error()}

		return c.JSON(http.StatusInternalServerError, sErr)
	}

	return c.JSON(http.StatusOK, k)
}

// @Summary      Stop Kpture
// @Description  Stop Kpture
// @Produce      json
// @Tags         kptures
// @Param        uuid  path      string  true  "capture uuid"
// @Failure      500   {object}  serverError
// @Failure      404  {object}  serverError
// @Success      200   {object}  capture.Kpture
// @Router       /api/v1/kpture/{uuid}/stop [put]
func (s *Server) stopKpture(c echo.Context) error {
	s.logger.Debug("stopKpture")

	uuid := c.Param("uuid")

	k := s.GetKpture(uuid)
	if k == nil {
		sErr := serverError{errors.New("capture not found").Error()}

		return c.JSON(http.StatusNotFound, sErr)
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
// @Failure      500   {object}  serverError
// @Failure      404   {object}  serverError
// @Router       /api/v1/kpture/{uuid}/download [get]
func (s *Server) downLoadKpture(c echo.Context) error {
	s.logger.Debug("downloadKpture")

	uuid := c.Param("uuid")

	kpture := s.GetKpture(uuid)
	if kpture == nil {
		sErr := serverError{errors.New("capture not found").Error()}

		return c.JSON(http.StatusNotFound, sErr)
	}

	if kpture.Status != capture.KptureStatusTerminated {
		sErr := serverError{errors.New("capture is running").Error()}

		return c.JSON(http.StatusBadRequest, sErr)
	}

	return c.Redirect(http.StatusMovedPermanently, filepath.Join("/captures", kpture.UUID, kpture.Name+".tar"))
}

// @Summary      Get kapture
// @Description  Get kapture
// @Produce      json
// @Tags         kptures
// @Param        uuid  path      string  true  "capture uuid"
// @Failure      500   {object}  serverError
// @Failure      404   {object}  serverError
// @Success      200   {object}  capture.Kpture
// @Router       /api/v1/kpture/{uuid} [get]
func (s *Server) getKpture(c echo.Context) error {
	s.logger.Debug("getKpture")

	uuid := c.Param("uuid")
	kpture := s.GetKpture(uuid)

	if kpture == nil {
		sErr := serverError{errors.New("capture not found").Error()}

		return c.JSON(http.StatusNotFound, sErr)
	}

	return c.JSON(http.StatusOK, kpture)
}

// @Summary      Get kaptures
// @Description  Get kaptures
// @Tags         kptures
// @Produce      json
// @Failure      500   {string}  string
// @Failure      404   {object}  serverError
// @Success      200  {object}  map[string]capture.Kpture
// @Router       /api/v1/kptures [get]
func (s *Server) getKptures(c echo.Context) error {
	s.logger.Debug("getKptures")

	return c.JSON(http.StatusOK, s.kptures)
}
