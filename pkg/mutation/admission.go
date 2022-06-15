package mutation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func admissionReviewFromRequest(r *http.Request) (*admissionv1.AdmissionReview, error) {
	// Validate that the incoming content type is correct.
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("expected application/json content-type")
	}

	admissionReviewRequest := &admissionv1.AdmissionReview{}

	err := json.NewDecoder(r.Body).Decode(&admissionReviewRequest)
	if err != nil {
		return nil, err
	}
	return admissionReviewRequest, nil
}

func admissionResponseFromReview(admReview *admissionv1.AdmissionReview) (*admissionv1.AdmissionResponse, error) {
	// check if valid pod resource
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admReview.Request.Resource != podResource {
		err := fmt.Errorf("did not receive pod, got %s", admReview.Request.Resource.Resource)
		return nil, err
	}

	admissionResponse := &admissionv1.AdmissionResponse{}

	// Decode the pod from the AdmissionReview.
	rawRequest := admReview.Request.Object.Raw
	pod := corev1.Pod{}

	err := json.NewDecoder(bytes.NewReader(rawRequest)).Decode(&pod)
	if err != nil {
		err := fmt.Errorf("error decoding raw pod: %v", err)
		return nil, err
	}

	var patch string
	patchType := v1.PatchTypeJSONPatch

	log.Println("pod has following labels", pod.Labels)
	// if _, ok := pod.Labels["tcpdump-sidecar"]; ok {
	patch = `[{ "op": "add", "path": "/metadata/labels/kpture-agent", "value": "true" },{"op":"add","path":"/spec/containers/1","value":{"image":"kpture/agent:latest","imagePullPolicy":"IfNotPresent","name":"ubuntu","ports":[{"name":"agent","containerPort":10000,"protocol":"TCP"}]}}]`
	// }

	admissionResponse.Allowed = true
	if patch != "" {
		log.Println("patching the pod with:", patch)
		admissionResponse.PatchType = &patchType
		admissionResponse.Patch = []byte(patch)
	}

	return admissionResponse, nil
}
