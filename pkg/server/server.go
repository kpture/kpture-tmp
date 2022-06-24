package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"newproxy/pkg/agent"
	"newproxy/pkg/capture"
	"newproxy/pkg/logger"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	logger      *logrus.Entry
	Agents      []capture.Agent
	kptures     map[string]*capture.Kpture
	Echo        *echo.Echo
	kubeclient  kubernetes.Interface
	storagePath string
}

func NewServer(kc *kubernetes.Clientset, storagePath string) *Server {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))

	serv := &Server{
		Echo:        e,
		logger:      logger.NewLogger("http"),
		kubeclient:  kc,
		storagePath: storagePath,
		kptures:     make(map[string]*capture.Kpture),
		Agents:      []capture.Agent{},
	}

	err := os.MkdirAll(storagePath, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = serv.LoadCaptures()
	if err != nil {
		log.Fatal(err)
	}

	serv.RegisterRoutes()

	return serv
}

func (s *Server) Start() {
	s.logger.Debug("Starting http backend server")
	s.Echo.Logger.Fatal(s.Echo.Start(":8080"))
}

func (s *Server) RegisterK8sAgents() error {
	list, err := s.kubeclient.CoreV1().
		Pods("").
		List(context.Background(), v1.ListOptions{LabelSelector: "kpture-agent=true"})
	if err != nil {
		s.logger.Error(err)

		return err
	}

	agents := []capture.Agent{}

	for _, pod := range list.Items {
		metadata := agent.Metadata{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Type:      agent.TypeKubernetes,
			TargetURL: pod.Status.PodIP + ":10000",
		}

		agents = append(agents, agent.NewCaptureSocket(metadata))
	}

	s.Agents = agents

	return nil
}

func (s *Server) StartKpture(name string, agents []capture.Agent) (*capture.Kpture, error) {
	k, err := capture.NewKpture(name, s.storagePath, agents)
	if err != nil {
		return nil, err
	}

	s.kptures[k.UUID] = k
	k.Start()

	return k, nil
}

func (s *Server) StopKpture(name string) error {
	if k, ok := s.kptures[name]; ok {
		return k.Stop()
	} else {
		return fmt.Errorf("kpture %s not found", name)
	}
}

func (s *Server) GetKpture(name string) *capture.Kpture {
	return s.kptures[name]
}

func (s *Server) GetKptures() map[string]*capture.Kpture {
	return s.kptures
}

func (s *Server) LoadCaptures() error {
	curr := capture.Kpture{}
	err := filepath.Walk(s.storagePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Name() == "descriptor.json" {
				dat, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				err = json.Unmarshal(dat, &curr)
				if err != nil {
					return err
				}
				s.kptures[curr.UUID] = &curr
			}

			return nil
		})

	return err
}
