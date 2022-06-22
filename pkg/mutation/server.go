package mutation

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"newproxy/pkg/logger"

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
	log.Println("Starting mutation webhook server")
	err := s.http.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleMutate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WEBHOOK DEMAND")
	admReview, err := admissionReviewFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("error getting admission review from request: %v", err)
		return
	}
	admResp, err := admissionResponseFromReview(admReview)
	if err != nil {
		w.WriteHeader(400)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	// the final response will be another admission review
	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admResp
	admissionReviewResponse.SetGroupVersionKind(admReview.GroupVersionKind())
	admissionReviewResponse.Response.UID = admReview.Request.UID
	resp, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		msg := fmt.Errorf("error marshaling response: %v", err)
		log.Println(msg)
		w.WriteHeader(500)
		_, err = w.Write([]byte(msg.Error()))
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf("allowing pod as %v", string(resp))
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println(err)
	}
}
