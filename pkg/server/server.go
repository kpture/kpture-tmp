package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"newproxy/pkg/agent"
	"newproxy/pkg/capture"
	"newproxy/pkg/certs"
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
	profiles    map[string]profile
	Echo        *echo.Echo
	kubeclient  kubernetes.Interface
	storagePath string
	certs       *certs.CertificatesHandler
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
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))

	serv := &Server{
		Echo:        echoServer,
		logger:      logger.NewLogger("http"),
		kubeclient:  k8sClient,
		storagePath: storagePath,
		profiles:    make(map[string]profile),
		Agents:      []capture.Agent{},
		certs:       certs.NewCertificatesHandler(k8sClient),
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

func (s *Server) Start() error {
	tlsConfig, err := s.certs.Get()
	if err != nil {
		return err
	}

	serverhttps := http.Server{
		Addr:      "0.0.0.0:443",
		Handler:   s.Echo, // set Echo as handler
		TLSConfig: tlsConfig,
		// ReadTimeout: 30 * time.Second, // use custom timeouts
	}

	go func() {
		s.logger.Debug("Starting httpsbackend server")
		if err := serverhttps.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
			s.Echo.Logger.Fatal(err)
		}
	}()

	serverhttp := http.Server{
		Addr:    "0.0.0.0:80",
		Handler: s.Echo, // set Echo as handler
		// ReadTimeout: 30 * time.Second, // use custom timeouts
	}

	s.logger.Debug("Starting http backend server")
	if err := serverhttp.ListenAndServe(); err != http.ErrServerClosed {
		s.Echo.Logger.Fatal(err)
	}

	return err
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
