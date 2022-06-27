package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"newproxy/pkg/agent"
	"newproxy/pkg/capture"
	"newproxy/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	logger      *logrus.Entry
	Agents      []capture.Agent
	kptures     map[string]*capture.Kpture
	profiles    map[string]profile
	Echo        *echo.Echo
	kubeclient  kubernetes.Interface
	storagePath string
}

func NewServer(k8sClient kubernetes.Interface, storagePath string) (*Server, error) {
	if k8sClient == nil {
		return nil, errors.New("empty kubernetes clientset")
	}

	if len(storagePath) == 0 {
		return nil, errors.New("storagePath not provided")
	}

	echoServer := echo.New()
	echoServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))

	serv := &Server{
		Echo:        echoServer,
		logger:      logger.NewLogger("http"),
		kubeclient:  k8sClient,
		storagePath: storagePath,
		kptures:     make(map[string]*capture.Kpture),
		profiles:    make(map[string]profile),
		Agents:      []capture.Agent{},
	}

	err := os.MkdirAll(storagePath, fs.ModePerm)
	if err != nil {
		return nil, errors.WithMessage(err, "error create kpture directory")
	}

	err = serv.LoadCaptures()
	if err != nil {
		serv.logger.Warn("some captures could not be loaded, discarding ....")
	}

	if _, ok := serv.profiles[defaultProfile]; !ok {
		serv.profiles[defaultProfile] = profile{Name: defaultProfile, kptures: make(map[string]*capture.Kpture)}
	}

	serv.RegisterRoutes()

	return serv, nil
}

func (s *Server) Start() {
	s.logger.Debug("Starting http backend server")
	s.Echo.Logger.Fatal(s.Echo.Start(":8080"))
}

func (s *Server) RegisterK8sAgents() error {
	list, err := s.kubeclient.CoreV1().
		Pods("").
		List(context.Background(), v1.ListOptions{LabelSelector: agentPodSelector})
	if err != nil {
		s.logger.Error(err)

		return errors.WithMessage(err, "error fetching k8s api")
	}

	agents := []capture.Agent{}

	for _, pod := range list.Items {
		metadata := agent.Metadata{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Type:      agent.TypeKubernetes,
			TargetURL: fmt.Sprintf("%s:%d", pod.Status.PodIP, agentPort),
		}

		agents = append(agents, agent.NewCaptureSocket(metadata))
	}

	s.Agents = agents

	return nil
}

func (s *Server) StopKpture(name string) error {
	if k, ok := s.kptures[name]; ok {
		return errors.WithMessage(k.Stop(), "error stopping kpture")
	}

	return errors.New("kpture not found")
}

func (s *Server) GetKpture(name string) *capture.Kpture {
	return s.kptures[name]
}

func (s *Server) GetKptures() map[string]*capture.Kpture {
	return s.kptures
}

func (s *Server) LoadCaptures() error {
	files, err := ioutil.ReadDir(s.storagePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		s.logger.Infof("creating profile %s", file.Name())
		s.profiles[file.Name()] = profile{
			Name:    file.Name(),
			kptures: make(map[string]*capture.Kpture),
		}

		err := filepath.Walk(filepath.Join(s.storagePath, file.Name()),
			func(path string, info os.FileInfo, err error) error {
				curr := capture.Kpture{}

				if info.Name() == "descriptor.json" {
					dat, err := os.ReadFile(path)
					if err != nil {
						return errors.WithMessage(err, "error reading kpture descriptor")
					}
					err = json.Unmarshal(dat, &curr)
					if err != nil {
						return errors.WithMessage(err, "error unmarshaling kpture descriptor")
					}
					s.logger.Infof("adding capture %s to profile %s", curr.Name, curr.ProfileName)
					s.profiles[curr.ProfileName].kptures[curr.UUID] = &curr
					// t.kptures[curr.UUID] = &curr
				}

				return nil
			})
		if err != nil {
			return errors.WithMessage(err, "error loading captures")
		}
	}

	return errors.WithMessage(err, "error walking recursively trhough files")
}
