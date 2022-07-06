package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @Summary      Get hostfile
// @Description  Get hostfile
// @Produce      plain
// @Tags         wireshark
// @Failure      500  {object}  serverError
// @Success      200  {string}  string
// @Router       /wireshark/hostfile [get]
func (s *Server) getHostFile(context echo.Context) error {
	s.logger.Debug("getHostFile")

	var response string

	pods, err := s.kubeclient.CoreV1().Pods("").List(context.Request().Context(), v1.ListOptions{})
	if err != nil {
		sErr := serverError{errors.WithMessage(err, "could not fetch kubernetes api").Error()}

		if err := context.JSON(http.StatusInternalServerError, sErr); err != nil {
			return errors.WithMessage(err, "error writing http response")
		}

		return nil
	}

	for _, pod := range pods.Items {
		response += fmt.Sprintf("%s %s\n", pod.Status.PodIP, pod.Name)
	}

	if err := context.String(http.StatusOK, response); err != nil {
		return errors.WithMessage(err, "error writing http response")
	}

	return nil
}
