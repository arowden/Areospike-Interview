package client

import (
	"errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sexample/client/clientfakes"
	"reflect"
	"testing"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 k8s.io/client-go/kubernetes/typed/core/v1.NamespaceInterface
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 k8s.io/client-go/kubernetes/typed/core/v1.CoreV1Interface
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 k8s.io/client-go/kubernetes/typed/core/v1.PodInterface

func TestClient_Create(t *testing.T) {
	coreClient, _, nsClient := getClients()
	client := NewClient(coreClient)
	expectedNamespace := "nsClient a"
	nsClient.CreateReturns(nil, errors.New(""))
	err := client.CreateNamespace(expectedNamespace)
	if err == nil {
		t.Log("Expected error was not returned.")
		t.Fail()
	}

	nsClient.CreateReturns(nil, nil)
	err = client.CreateNamespace(expectedNamespace)
	if err != nil {
		t.Log("Unexpected error returned.")
		t.Fail()
	}

	_, gotNamespace, _ := nsClient.CreateArgsForCall(0)
	if gotNamespace.Name != expectedNamespace {
		t.Logf("CreateNamespace failed. Got nsClient: %s Expected nsClient: %s", gotNamespace, expectedNamespace)
	}
}

func TestClient_Delete(t *testing.T) {
	coreClient, _, nsClient := getClients()
	client := NewClient(coreClient)
	expectedNamespaceToDelete := "client a"
	nsClient.DeleteReturns(errors.New(""))
	err := client.DeleteNamespace(expectedNamespaceToDelete)
	if err == nil {
		t.Log("Expected error was not returned.")
		t.Fail()
	}

	nsClient.DeleteReturns(nil)
	err = client.CreateNamespace(expectedNamespaceToDelete)
	if err != nil {
		t.Log("Unexpected error returned.")
		t.Fail()
	}

	_, gotNamespace, _ := nsClient.DeleteArgsForCall(0)
	if gotNamespace != expectedNamespaceToDelete {
		t.Logf("CreateNamespace failed. Got client: %s Expected client: %s", gotNamespace, expectedNamespaceToDelete)
	}
}

func TestClient_List(t *testing.T) {
	coreClient, _, nsClient := getClients()
	client := NewClient(coreClient)
	results := &v1.NamespaceList{
		Items: []v1.Namespace{
			{ObjectMeta: metav1.ObjectMeta{Name: "Name1"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "Name2"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "Name3"}},
		},
	}

	expected := []string{
		"Name1",
		"Name2",
		"Name3",
	}

	nsClient.ListReturns(results, nil)
	got, err := client.ListNamespaces()
	if err != nil {
		t.Logf("Received an unexpected error: %s", err.Error())
		t.Fail()
	}

	if !reflect.DeepEqual(expected, got) {
		t.Logf("Expected: %v Got: %v", expected, got)
		t.Fail()
	}
}

func getClients() (*clientfakes.FakeCoreV1Interface, *clientfakes.FakePodInterface, *clientfakes.FakeNamespaceInterface) {
	coreClient := &clientfakes.FakeCoreV1Interface{}
	nsClient := &clientfakes.FakeNamespaceInterface{}
	podClient := &clientfakes.FakePodInterface{}
	coreClient.NamespacesReturns(nsClient)
	coreClient.PodsReturns(podClient)
}
