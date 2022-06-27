package server_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"newproxy/pkg/server"
	testdescriptors_test "newproxy/testsutils/testdescriptors"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestNewServer(t *testing.T) {
	t.Parallel()

	k8sClient := fake.NewSimpleClientset()
	output := t.TempDir()

	t.Run("Test_new_server_nil_client", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		got, err := server.NewServer(nil, output)
		assert.NotNil(err, "err shall not be nil")
		assert.Nil(got)
	})

	t.Run("Test_new_server_empty_storagePath", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		got, err := server.NewServer(k8sClient, "")
		assert.NotNil(err, "err shall not be nil")
		assert.Nil(got)
	})

	t.Run("Test_new_server_permission_error_storagePath", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		got, err := server.NewServer(k8sClient, "/kpture")
		assert.NotNil(err, "err shall not be nil")
		assert.Nil(got)
	})

	// This test is about an invalid capture loading from filesystem
	// This could be an invalid kpture json descriptor, due to an struct/Api update
	// In this case we discard the error, the capture will not loaded
	t.Run("Test_new_server_invalid_capture_fs", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		location := t.TempDir()
		uuid := "18388a08-1845-4cfc-891b-29095027babe"
		err := os.MkdirAll(filepath.Join(location, uuid), fs.ModePerm)
		assert.Nil(err)

		f, err := os.OpenFile(filepath.Join(location, "descriptor.json"), os.O_CREATE|os.O_RDWR, fs.ModePerm)
		assert.Nil(err)
		assert.NotNil(f)
		_, err = f.WriteString(testdescriptors_test.InvalidDescriptor)
		assert.Nil(err)

		got, err := server.NewServer(k8sClient, location)
		assert.Nil(err, "err shall not be nil")
		assert.NotNil(got)
	})
}

func getPodsList() *v1.PodList {
	return &v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx",
					Namespace: "testing",
					Labels: map[string]string{
						"kpture-agent": "true",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx2",
					Namespace: "testing",
					Labels: map[string]string{
						"kpture-agent": "true",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx2",
					Namespace: "testing3",
				},
			},
		},
	}
}

func TestServer_RegisterK8sAgents(t *testing.T) {
	t.Parallel()

	output := t.TempDir()

	t.Run("Test_registerk8sAgents", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		k8sClient := fake.NewSimpleClientset()
		server, err := server.NewServer(k8sClient, output)
		assert.Nil(err)
		assert.NotNil(server)

		fakeCoreV1, ok := k8sClient.CoreV1().(*fakecorev1.FakeCoreV1)
		assert.True(ok)

		fakeCoreV1.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
			return true, getPodsList(), nil
		})

		err = server.RegisterK8sAgents()
		assert.Equal(len(server.Agents), 2)
		assert.Nil(err)
	})

	t.Run("Test_errorK8s", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)

		k8sClient := fake.NewSimpleClientset()
		server, err := server.NewServer(k8sClient, output)
		assert.Nil(err)
		assert.NotNil(server)

		fakeCoreV1, ok := k8sClient.CoreV1().(*fakecorev1.FakeCoreV1)
		assert.True(ok)

		fakeCoreV1.PrependReactor("list", "pods", func(action ktesting.Action) (bool, runtime.Object, error) {
			return true, &v1.PodList{}, errors.New("k8s error fetch")
		})

		err = server.RegisterK8sAgents()
		assert.NotNil(err)
	})
}

func TestServer_LoadCaptures(t *testing.T) {
	t.Parallel()

	t.Run("Test_new_server_invalid_capture_fs", func(t *testing.T) {
		t.Parallel()

		k8sClient := fake.NewSimpleClientset()
		assert := assert.New(t)

		location := t.TempDir()
		uuid := "18388a08-1845-4cfc-891b-29095027babe"
		err := os.MkdirAll(filepath.Join(location, uuid), fs.ModePerm)
		assert.Nil(err)

		f, err := os.OpenFile(filepath.Join(location, "descriptor.json"), os.O_CREATE|os.O_RDWR, fs.ModePerm)
		assert.Nil(err)
		assert.NotNil(f)
		_, err = f.WriteString(testdescriptors_test.InvalidDescriptor)
		assert.Nil(err)

		got, err := server.NewServer(k8sClient, location)
		assert.Nil(err, "err shall not be nil")
		assert.NotNil(got)
	})
}
