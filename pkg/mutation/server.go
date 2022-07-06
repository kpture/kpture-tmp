package mutation

import (
	"encoding/json"
	"log"
	"net/http"

	"newproxy/pkg/certs"
	"newproxy/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	logger *logrus.Entry
	http   *http.Server
	client kubernetes.Interface
	certs  *certs.CertificatesHandler
}

func NewMutationWebHookServer(client kubernetes.Interface) (*Server, error) {
	s := &Server{
		logger: logger.NewLogger("mutation"),
		client: client,
		certs:  certs.NewCertificatesHandler(client),
	}

	tlsConfig, err := s.certs.Get()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	s.http = &http.Server{
		Addr:      "0.0.0.0:443",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	return s, nil
}

func (s *Server) Start() {
	s.logger.Debug("Starting k8s mutation server")

	if err := s.http.ListenAndServeTLS("", ""); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func HandleMutate(c echo.Context) error {
	l := logger.NewLogger("handleMutate")

	admReview, err := admissionReviewFromRequest(c.Request())
	if err != nil {
		http.Error(c.Response().Writer, err.Error(), http.StatusBadRequest)
		l.Error(errors.WithMessage(err, "could not validate webook request"))
	}

	admResp, err := admissionResponseFromReview(admReview)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusBadRequest)
		_, err = c.Response().Writer.Write([]byte(err.Error()))
		l.Error(errors.WithMessage(err, "could not write webook http response"))
	}

	// the final response will be another admission review
	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admResp
	admissionReviewResponse.SetGroupVersionKind(admReview.GroupVersionKind())
	admissionReviewResponse.Response.UID = admReview.Request.UID

	resp, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		c.Response().Writer.WriteHeader(http.StatusInternalServerError)

		_, err = c.Response().Writer.Write([]byte(err.Error()))
		if err != nil {
			l.Error(errors.WithMessage(err, "could not write webook http response"))
		}

		l.Error(errors.WithMessage(err, "could not marshal admission response"))
	}

	c.Response().Writer.Header().Set("Content-Type", "application/json")

	_, err = c.Response().Write(resp)
	if err != nil {
		l.Error(errors.WithMessage(err, "could not write webook http response"))
	}

	return nil
}
