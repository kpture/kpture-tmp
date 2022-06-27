package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (s *Server) getProfile(context echo.Context) (*profile, error) {
	profilename, _, valid := context.Request().BasicAuth()

	if !valid || len(profilename) == 0 {
		defaultProfile, ok := s.profiles[defaultProfile]
		if !ok {
			return nil, fmt.Errorf("could not load default profile")
		}
		return &defaultProfile, nil
	}

	profile, ok := s.profiles[profilename]
	if !ok {
		return nil, fmt.Errorf("could not load profile")
	}

	return &profile, nil
}
