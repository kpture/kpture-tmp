package mutation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func admissionReviewFromRequest(r *http.Request) (*v1.AdmissionReview, error) {
	// Validate that the incoming content type is correct.
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("expected application/json content-type")
	}

	admissionReviewRequest := &v1.AdmissionReview{}

	err := json.NewDecoder(r.Body).Decode(&admissionReviewRequest)
	if err != nil {
		return nil, err
	}

	return admissionReviewRequest, nil
}

func admissionResponseFromReview(admReview *v1.AdmissionReview) (*v1.AdmissionResponse, error) {
	// check if valid pod resource
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admReview.Request.Resource != podResource {
		err := fmt.Errorf("did not receive pod, got %s", admReview.Request.Resource.Resource)

		return nil, err
	}

	admissionResponse := &v1.AdmissionResponse{}

	// Decode the pod from the AdmissionReview.
	rawRequest := admReview.Request.Object.Raw
	pod := corev1.Pod{}

	err := json.NewDecoder(bytes.NewReader(rawRequest)).Decode(&pod)
	if err != nil {
		return nil, errors.WithMessage(err, "error decoding raw pod")
	}

	patchType := v1.PatchTypeJSONPatch

	admissionResponse.Allowed = true
	admissionResponse.PatchType = &patchType
	admissionResponse.Patch = []byte(injectionPatch)

	return admissionResponse, nil
}
