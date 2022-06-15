package mutation

import (
	"bytes"
	"context"
	"fmt"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func createMutationConfig(caCert *bytes.Buffer) error {

	config := ctrl.GetConfigOrDie()
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	path := "/mutate"
	fail := admissionregistrationv1.Fail
	se := admissionregistrationv1.SideEffectClassNone
	mutateconfig := &admissionregistrationv1.MutatingWebhookConfiguration{

		ObjectMeta: metav1.ObjectMeta{
			Name: "kpture-server-mutate",
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{{
			SideEffects:             &se,
			AdmissionReviewVersions: []string{"v1"},
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"kpture-agent": "enabled",
				},
			},
			Name: "tls.kpture.svc",
			ClientConfig: admissionregistrationv1.WebhookClientConfig{
				CABundle: caCert.Bytes(), // CA bundle created earlier
				Service: &admissionregistrationv1.ServiceReference{
					Name:      "tls",
					Namespace: "kpture",
					Path:      &path,
				},
			},
			Rules: []admissionregistrationv1.RuleWithOperations{{Operations: []admissionregistrationv1.OperationType{
				admissionregistrationv1.Create},
				Rule: admissionregistrationv1.Rule{
					APIGroups:   []string{""},
					APIVersions: []string{"v1"},
					Resources:   []string{"pods"},
				},
			}},
			FailurePolicy: &fail,
		}},
	}

	fmt.Println(kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(context.Background(), mutateconfig.Name, metav1.DeleteOptions{}))

	if _, err := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(context.Background(), mutateconfig, metav1.CreateOptions{}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Mutate config created")
	}

	return nil
}
