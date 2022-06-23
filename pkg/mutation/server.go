package mutation

import (
	"encoding/json"
	"log"
	"net/http"
	"newproxy/pkg/logger"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
)

type Server struct {
	logger *logrus.Entry
	http   *http.Server
}

func NewMutationWebHookServer() (*Server, error) {
	s := &Server{logger: logger.NewLogger("mutation")}

	tlsConfig, err := genCerts()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", handleMutate)

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

func handleMutate(w http.ResponseWriter, r *http.Request) {
	l := logger.NewLogger("handleMutate")

	admReview, err := admissionReviewFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		l.Error(errors.WithMessage(err, "could not validate webook request"))
	}

	admResp, err := admissionResponseFromReview(admReview)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		l.Error(errors.WithMessage(err, "could not write webook http response"))
	}

	// the final response will be another admission review
	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admResp
	admissionReviewResponse.SetGroupVersionKind(admReview.GroupVersionKind())
	admissionReviewResponse.Response.UID = admReview.Request.UID

	resp, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			l.Error(errors.WithMessage(err, "could not write webook http response"))
		}

		l.Error(errors.WithMessage(err, "could not marshal admission response"))
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		l.Error(errors.WithMessage(err, "could not write webook http response"))
	}
}
