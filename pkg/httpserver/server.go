package httpserver

import (
	_ "net/http/pprof"
	"newproxy/pkg/logger"

	"newproxy/pkg/capture"

	echoPrometheus "github.com/globocom/echo-prometheus"
	echopprof "github.com/hiko1129/echo-pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type KptureServer struct {
	Echo             *echo.Echo
	kubeclient       kubernetes.Interface
	CaptureManager   *capture.Manager
	capturesBasePath string
	logger           *logrus.Entry
}

func NewServer(kc *kubernetes.Clientset, ch *capture.Manager, bp string) *KptureServer {
	e := echo.New()
	echopprof.Wrap(e)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: utils.MiddleWareFormat,
	// }))

	serv := &KptureServer{
		Echo:             e,
		logger:           logger.NewLogger("http"),
		kubeclient:       kc,
		CaptureManager:   ch,
		capturesBasePath: bp,
	}
	serv.registerRoutes()
	return serv
}

func (s *KptureServer) Start() {

	s.Echo.Logger.Fatal(s.Echo.Start(":8080"))
}

func (s *KptureServer) registerRoutes() {
	//Serve frontend website and captures file system

	s.Echo.Use(echoPrometheus.MetricsMiddleware())
	s.Echo.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	s.Echo.Static("/captures/", s.capturesBasePath)
	s.Echo.Static("/", "build")
	s.Echo.GET("/api/v2/kubernetes/pods", s.getAllPods)
	s.Echo.GET("/ws", s.hello)
	s.Echo.POST("/api/v2/kubernetes/capture", s.startPodCapture)
	s.Echo.GET("/api/v2/kubernetes/capture/:uuid", s.getPodCapture)
	s.Echo.PUT("/api/v2/kubernetes/capture/:uuid/stop", s.stopPodCapture)
	s.Echo.GET("/api/v2/kubernetes/capture/:uuid/download", s.downLoadCapture)
	s.Echo.GET("/api/v2/kubernetes/captures", s.getCaptures)
}
