package server

import (
	"net/http"

	"newproxy/pkg/agent"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// getAgents godoc
// @Summary      Get Agents
// @Description  Get Agents
// @Produce      json
// @Tags         agents
// @Failure      500  {object}  serverError
// @Success      200  {object}  []agent.Metadata
// @Router       /api/v1/agents [get]
func (s *Server) getAgents(context echo.Context) error {
	context.Request().Header.Set(echo.HeaderXRequestID, "getAgents")

	resp := []*agent.Metadata{}

	if err := s.RegisterK8sAgents(); err != nil {
		sErr := serverError{errors.WithMessage(err, "could not fetch k8s agents").Error()}

		return echo.NewHTTPError(http.StatusInternalServerError, sErr)
	}

	for _, agent := range s.Agents {
		resp = append(resp, &agent.Info().Metadata)
	}

	if err := context.JSON(http.StatusOK, resp); err != nil {
		return errors.WithMessage(err, "error writing http response")
	}

	return nil
}
