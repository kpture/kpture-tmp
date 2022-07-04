package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (s *Server) getProfile(context echo.Context) (*profile, error) {
	profilename, _, valid := context.Request().BasicAuth()

	fmt.Println(profilename)
	if !valid || len(profilename) == 0 {
		s.logger.Info("Using default profile")
		defaultProfile, ok := s.profiles[defaultProfile]
		if !ok {
			return nil, fmt.Errorf("could not load default profile")
		}
		return &defaultProfile, nil
	}

	s.logger.Info("Using profile", profilename)

	profile, ok := s.profiles[profilename]
	if !ok {
		return nil, fmt.Errorf("could not load profile")
	}

	fmt.Println(profile.Name)
	return &profile, nil
}
