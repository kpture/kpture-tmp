package server

import (
	"fmt"
	"net/http"

	"newproxy/pkg/capture"
	"newproxy/pkg/mutation"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary      Get enabled kubernetes namespaces
// @Description  Get enabled kubernetes namespaces
// @Accept       json
// @Produce      json
// @Tags         kubernetes
// @Failure      500  {object}  serverError
// @Success      200  {object}  []string
// @Router       /api/v1/kpture/k8s/namespaces [GET]
func (s *Server) getKubernetesEnabledNs(context echo.Context) error {
	s.logger.Debug("getKubernetesEnabledNs")

	namespaces := []string{}
	label := fmt.Sprintf("%s=%s", mutation.NameSpaceSelectorLabel, mutation.NameSpaceSelectorValue)
	opts := v1.ListOptions{LabelSelector: label}
	nsCli := s.kubeclient.CoreV1().Namespaces()

	nsResponse, err := nsCli.List(context.Request().Context(), opts)
	if err != nil {
		sErr := serverError{errors.WithMessage(err, "could not fetch kubernetes api").Error()}

		if err := context.JSON(http.StatusInternalServerError, sErr); err != nil {
			return errors.WithMessage(err, "error writing http response")
		}

		return nil
	}

	for _, currns := range nsResponse.Items {
		namespaces = append(namespaces, currns.Name)
	}

	if err := context.JSON(http.StatusOK, namespaces); err != nil {
		return errors.WithMessage(err, "error writing http response")
	}

	return nil
}

// @Summary      Start namespace kpture
// @Description  Start namespace kpture
// @Accept       json
// @Produce      json
// @Tags         kubernetes
// @Param        data  body      kptureNamespaceRequest  true  "namespace for capture"
// @Failure      500   {object}  serverError
// @Failure      400   {object}  serverError
// @Success      200   {object}  capture.Kpture
// @Header       200   {string}  Websocket  ""
// @Router       /api/v1/kpture/k8s/namespace [post]
func (s *Server) startNamespacedKpture(context echo.Context) error {
	s.logger.Debug("startNamespacedKpture")

	profile, err := s.getProfile(context)
	if err != nil {
		sErr := serverError{"profile does not exist"}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return err
	}

	kptureReq := &kptureNamespaceRequest{}
	if err := context.Bind(kptureReq); err != nil {
		s.logger.Errorf("error binding request: %v", err)
		sErr := serverError{errors.WithMessage(err, "cannot bind http request body").Error()}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	if kptureReq.KptureName == "" {
		sErr := serverError{errors.New("invalid kptureName").Error()}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	agents := []capture.Agent{}

	for _, agent := range s.Agents {
		if agent.Info().Metadata.Namespace == kptureReq.KptureNamespace {
			agents = append(agents, agent)
		}
	}

	if len(agents) == 0 {
		sErr := serverError{"no agent in selected namespace"}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	kpture, err := capture.NewKpture(kptureReq.KptureName, profile.Name, s.storagePath, agents)
	if err != nil {
		s.logger.Errorf("error starting capture: %v", err)
		sErr := serverError{errors.WithMessage(err, "error starting kpture").Error()}

		if err := context.JSON(http.StatusInternalServerError, sErr); err != nil {
			return errors.WithMessage(err, "could not write HTTP response")
		}

		return nil
	}
	profile.kptures[kpture.UUID] = kpture
	kpture.Start()

	if err := context.JSON(http.StatusOK, kpture); err != nil {
		return errors.WithMessage(err, "could not write HTTP response")
	}

	return nil
}
