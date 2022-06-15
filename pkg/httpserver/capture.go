package httpserver

import (
	"fmt"
	"net/http"
	"newproxy/pkg/capture"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type captureRequest struct {
	PodMetadata []capture.PodMetadata `json:"pods,omitempty"`
	Name        string                `json:"name,omitempty"`
}

// @Summary      Start K8S capture
// @Description  Start K8S capture
// @Tags         kubernetes
// @Accept       json
// @Produce      json
// @Param        data  body      captureRequest  true  "selected pods for capture"
// @Failure      500  {string}  string
// @Success      200   {object}  capture.kpture
// @Header       200   {string}  Websocket  ""
// @Router       /kubernetes/capture [post]
func (s *KptureServer) startPodCapture(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "StartPodCapture")
	u := new(captureRequest)

	if err := c.Bind(u); err != nil || len(u.PodMetadata) == 0 {
		s.logger.Errorf("error binding request: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Println(u.PodMetadata)
	kpture, err := s.CaptureManager.StartCapture(u.PodMetadata, u.Name)
	if err != nil {
		s.logger.Errorf("error starting capture: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Logger().Info(u)
	return c.JSON(http.StatusOK, kpture)
}

// @Summary      Get K8S capture
// @Description  Get K8S capture
// @Tags         kubernetes
// @Produce      json
// @Param        uuid  path      string  true  "capture uuid"
// @Failure      500   {string}  string
// @Failure      404   {string}  string
// @Success      200   {object}  capture.kpture
// @Router       /kubernetes/capture/{uuid} [get]
func (s *KptureServer) getPodCapture(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "GetPodCapture")
	uuid := c.Param("uuid")
	s.Echo.Logger.Info(uuid)
	kpture := s.CaptureManager.GetCapture(uuid)
	if kpture.Status.CaptureState == capture.CaptureStatusReady {
		kpture.BufferSizeStr = ""
	}
	if kpture == nil {
		return c.JSON(http.StatusNotFound, "capture not found")
	}
	kpture.UpdatePkNumber()

	return c.JSON(http.StatusOK, kpture)
}

// getCaptures godoc
// @Summary      Get K8S captures
// @Description  Get K8S captures
// @Tags         kubernetes
// @Produce      json
// @Failure      500   {string}  string
// @Success      200  {object}  map[string]capture.kpture
// @Router       /kubernetes/captures [get]
func (s *KptureServer) getCaptures(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "GetCaptures")
	return c.JSON(http.StatusOK, s.CaptureManager.GetCaptures())
}

// @Summary      Stop K8S capture
// @Description  Stop K8S capture
// @Tags         kubernetes
// @Produce      json
// @Param        uuid  path      string  true  "capture uuid"
// @Failure      500   {string}  string
// @Failure      404   {string}  string
// @Success      200   {object}  capture.kpture
// @Router       /kubernetes/capture/{uuid}/stop [put]
func (s *KptureServer) stopPodCapture(c echo.Context) (err error) {
	c.Request().Header.Set(echo.HeaderXRequestID, "StopPodCapture")
	uuid := c.Param("uuid")
	kpture := s.CaptureManager.GetCapture(uuid)
	if kpture == nil {
		return c.JSON(http.StatusNotFound, "capture not found")
	}
	if kpture.Status.CaptureState == capture.CaptureStatusStarted {
		kpture, err = s.CaptureManager.StopCapture(uuid)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, kpture)
}

// @Summary      Download K8S capture
// @Description  Download K8S capture
// @Tags         kubernetes
// @Produce      json
// @Param        uuid  path      string  true  "capture uuid"
// @Failure      500   {string}  string
// @Failure      404   {string}  string
// @Success      200   {object}  capture.kpture
// @Router       /kubernetes/capture/{uuid}/download [get]
func (s *KptureServer) downLoadCapture(c echo.Context) error {
	c.Request().Header.Set(echo.HeaderXRequestID, "DownLoadCapture")
	uuid := c.Param("uuid")
	kpture := s.CaptureManager.GetCapture(uuid)
	if kpture == nil {
		return c.JSON(http.StatusNotFound, "capture not found")
	}
	if kpture.Status.CaptureState != capture.CaptureStatusReady {
		return c.JSON(http.StatusBadRequest, "capture not finished")
	}
	return c.Redirect(http.StatusFound, filepath.Join("/captures/", kpture.ArchiveLocation))
}
