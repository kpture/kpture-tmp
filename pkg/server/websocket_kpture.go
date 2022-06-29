package server

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (s *Server) kptureWebSocket(context echo.Context) error {
	s.logger.Debug("getKpture")

	uuid := context.Param("uuid")
	profileName := context.Param("profileName")

	profile, ok := s.profiles[profileName]
	if !ok {
		return errors.New("profile not found")
	}

	kpture := profile.kptures[uuid]
	if kpture == nil {
		return errors.New("kpture not found")
	}

	upgrader := websocket.Upgrader{}

	webSocket, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
	if err != nil {
		return errors.WithMessage(err, "error upgrading websocket")
	}
	defer webSocket.Close()

	for {
		time.Sleep(time.Second)

		b, err := json.Marshal(kpture)
		if err != nil {
			return errors.WithMessage(err, "error marshaling json info")
		}

		err = webSocket.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			s.logger.Info("could not send websocket message")

			return errors.WithMessage(err, "error sending websocket msg")
		}
	}
}
