package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"newproxy/pkg/capture"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// @Summary                    Start Kpture
// @Description                Start Kpture
// @Accept                     json
// @Produce      json
// @Tags         kptures
// @Param                      data  body      kptureRequest  true  "selected agents for capture"
// @Failure      500   {object}  serverError
// @Success                    200   {object}  capture.Kpture
// @Header                     200   {string}  Websocket  ""
// @Router                     /kpture [post]
// @securityDefinitions.basic  BasicAuth
func (s *Server) startKpture(context echo.Context) error {
	s.logger.Debug("startKpture")

	profile, err := s.getProfile(context)
	if err != nil {
		sErr := serverError{"profile does not exist"}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return err
	}

	kptureReq := &kptureRequest{}

	if err := context.Bind(kptureReq); err != nil {
		s.logger.Errorf("error binding request: %v", err)
		sErr := serverError{errors.WithMessage(err, "cannot bind request").Error()}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	agents := []capture.Agent{}

	for _, pod := range kptureReq.AgentsRequest {
		for _, agent := range s.Agents {
			if agent.Info().Metadata.Name == pod.Name && agent.Info().Metadata.Namespace == pod.Namespace {
				agents = append(agents, agent)
			}
		}
	}
	kpture, err := capture.NewKpture(kptureReq.KptureName, profile.Name, s.storagePath, agents, s.kubeclient)
	if err != nil {
		return errors.WithMessage(err, "error creating kpture")
	}
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

// @Summary                    Stop Kpture
// @Description                Stop Kpture
// @Produce                    json
// @Tags                       kptures
// @Param        uuid  path      string  true  "capture uuid"
// @Failure                    500   {object}  serverError
// @Failure      404   {object}  serverError
// @Success                    200   {object}  capture.Kpture
// @Router                     /kpture/{uuid}/stop [put]
// @securityDefinitions.basic  BasicAuth
func (s *Server) stopKpture(context echo.Context) error {
	s.logger.Debug("stopKpture")
	uuid := context.Param("uuid")

	profile, err := s.getProfile(context)
	if err != nil {
		sErr := serverError{"profile does not exist"}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return err
	}

	kpture := profile.kptures[uuid]
	if kpture == nil {
		sErr := serverError{errors.New("capture not found").Error()}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return nil
	}

	if err := kpture.Stop(); err != nil {
		return errors.WithMessage(err, "error stop kpture")
	}

	if err := context.JSON(http.StatusOK, kpture); err != nil {
		return errors.WithMessage(err, "could not write http response")
	}

	return nil
}

// @Summary                    Download kpture
// @Description                Download kpture
// @Produce                    json
// @Tags                       kptures
// @Param                      uuid  path      string  true  "capture uuid"
// @Failure                    500   {object}  serverError
// @Failure                    404   {object}  serverError
// @Router                     /kpture/{uuid}/download [get]
// @securityDefinitions.basic  BasicAuth
func (s *Server) downLoadKpture(context echo.Context) error {
	s.logger.Debug("downloadKpture")
	uuid := context.Param("uuid")
	profile, err := s.getProfile(context)
	if err != nil {
		sErr := serverError{"profile does not exist"}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return err
	}

	fmt.Println(profile.kptures)
	kpture := profile.kptures[uuid]
	if kpture == nil {
		sErr := serverError{errors.New("capture not found").Error()}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return nil
	}

	if kpture.Status != capture.KptureStatusTerminated {
		sErr := serverError{errors.New("capture is running").Error()}

		if err := context.JSON(http.StatusBadRequest, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return nil
	}

	redirect := filepath.Join("/captures", profile.Name, kpture.UUID, kpture.Name+".tar")
	fmt.Println(redirect)

	if err := context.Redirect(http.StatusFound, redirect); err != nil {
		return errors.WithMessage(err, "could not write http response")
	}

	return nil
}

// @Summary                    Get kapture
// @Description                Get kapture
// @Produce                    json
// @Tags                       kptures
// @Param                      uuid  path      string  true  "capture uuid"
// @Failure                    500   {object}  serverError
// @Failure                    404   {object}  serverError
// @Success                    200   {object}  capture.Kpture
// @Router                     /kpture/{uuid} [get]
// @securityDefinitions.basic  BasicAuth
func (s *Server) getKpture(context echo.Context) error {
	s.logger.Debug("getKpture")

	uuid := context.Param("uuid")

	profile, err := s.getProfile(context)
	if err != nil {
		sErr := serverError{"profile does not exist"}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return err
	}

	kpture := profile.kptures[uuid]

	if kpture == nil {
		sErr := serverError{errors.New("capture not found").Error()}
		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return nil
	}

	if err := context.JSON(http.StatusOK, kpture); err != nil {
		return errors.WithMessage(err, "could not write http response")
	}

	return nil
}

// @Summary      Delete kapture
// @Description  Delete kapture
// @Produce                    json
// @Tags                       kptures
// @Param                      uuid  path      string  true  "capture uuid"
// @Failure                    500   {object}  serverError
// @Failure                    404   {object}  serverError
// @Success      204
// @Router       /kpture/{uuid} [delete]
func (s *Server) deleteKpture(context echo.Context) error {
	s.logger.Debug("deletekpture")

	uuid := context.Param("uuid")

	profile, err := s.getProfile(context)
	if err != nil {
		sErr := serverError{"profile does not exist"}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return err
	}

	kpture := profile.kptures[uuid]

	if kpture == nil {
		sErr := serverError{errors.New("capture not found").Error()}
		logrus.Error(sErr)
		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return nil
	}

	err = os.RemoveAll(filepath.Join(s.storagePath, profile.Name, uuid))
	logrus.Error(err)
	if err != nil {
		logrus.Error(err)
		if err := context.JSON(http.StatusInternalServerError, err); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}
	}

	delete(profile.kptures, uuid)

	if err := context.JSON(http.StatusNoContent, kpture); err != nil {
		return errors.WithMessage(err, "could not write http response")
	}

	return nil
}

// @Summary                    Get kaptures
// @Description                Get kaptures
// @Tags                       kptures
// @Produce                    json
// @Failure                    500  {string}  string
// @Failure                    404  {object}  serverError
// @Success                    200  {object}  map[string]capture.Kpture
// @Router                     /kptures [get]
// @securityDefinitions.basic  BasicAuth
func (s *Server) getKptures(context echo.Context) error {
	s.logger.Debug("getKptures")

	profile, err := s.getProfile(context)
	if err != nil {
		sErr := serverError{"profile does not exist"}

		if err := context.JSON(http.StatusNotFound, sErr); err != nil {
			return errors.WithMessage(err, "could not write http response")
		}

		return err
	}

	if err := context.JSON(http.StatusOK, profile.kptures); err != nil {
		return errors.WithMessage(err, "could not write http response")
	}

	return nil
}
