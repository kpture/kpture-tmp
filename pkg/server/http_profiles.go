package server

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"newproxy/pkg/capture"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type profileResponse struct {
	Profiles []string `json:"profiles,omitempty"`
}

// @Summary      Get profiles
// @Description  Get profiles
// @Produce      json
// @Tags         profiles
// @Failure      500  {object}  serverError
// @Success      200  {object}  []string{}
// @Router       /profiles [get]
func (s *Server) getProfiles(context echo.Context) error {
	s.logger.Debug("getHostFile")

	response := []string{}

	for name := range s.profiles {
		response = append(response, name)
	}

	if err := context.JSON(http.StatusOK, response); err != nil {
		return errors.WithMessage(err, "could not write http response")
	}

	return nil
}

// @Summary      Create profile
// @Description  Create profile
// @Tags         profiles
// @Failure      500  {object}  serverError
// @Success      200
// @Param        profileName  path  string  true  "profileName"
// @Router       /profile/{profileName} [POST]
func (s *Server) createProfile(context echo.Context) error {
	s.logger.Debug("getHostFile")

	profileName := context.Param("profileName")

	if len(profileName) == 0 || profileName == "default" {
		sErr := serverError{"invalid profile name"}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	if _, ok := s.profiles[profileName]; ok {
		sErr := serverError{"profileName already exist"}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	s.profiles[profileName] = profile{
		Name:    profileName,
		kptures: make(map[string]*capture.Kpture),
	}

	err := os.MkdirAll(filepath.Join(s.storagePath, profileName), fs.ModePerm)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := context.JSON(http.StatusOK, nil); err != nil {
		return errors.WithMessage(err, "could not write http response")
	}

	return nil
}

// @Summary      Delete profile
// @Description  Delete profile
// @Tags         profiles
// @Failure      500  {object}  serverError
// @Success      204
// @Param        profileName  path  string  true  "profileName"
// @Router       /profile/{profileName} [DELETE]
func (s *Server) deleteProfile(context echo.Context) error {
	s.logger.Debug("getHostFile")

	profileName := context.Param("profileName")

	if len(profileName) == 0 {
		sErr := serverError{"invalid profile name"}

		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	if _, ok := s.profiles[profileName]; !ok {
		sErr := serverError{"profileName does not exist"}
		return echo.NewHTTPError(http.StatusNotFound, sErr)
	}

	err := os.RemoveAll(filepath.Join(s.storagePath, profileName))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	delete(s.profiles, profileName)

	if err := context.JSON(http.StatusNoContent, nil); err != nil {
		return errors.WithMessage(err, "could not write http response")
	}

	return nil
}
