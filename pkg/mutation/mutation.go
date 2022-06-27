package mutation

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
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

	admclient := kubeClient.
		AdmissionregistrationV1().
		MutatingWebhookConfigurations()

	wh, err := admclient.Get(context.Background(), webHookName, metav1.GetOptions{})
	if err != nil {
		return errors.WithMessage(err, "could not create admission configuration")
	}

	if wh != nil {
		if err := admclient.Delete(context.Background(), webHookName, metav1.DeleteOptions{}); err != nil {
			return errors.WithMessage(err, "could not delete admission configuration")
		}
	}

	mutateconfig := getMutationConfig(caCert.Bytes())

	_, err = admclient.
		Create(context.Background(), mutateconfig, metav1.CreateOptions{})

	if err != nil && !k8serr.IsAlreadyExists(err) {
		return errors.WithMessage(err, "could not create admission configuration")
	}

	return nil
}
