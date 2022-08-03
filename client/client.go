package client

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type Client struct {
	client corev1.CoreV1Interface
}

// NewClient returns a new client.
func NewClient(client corev1.CoreV1Interface) *Client {
	return &Client{
		client: client,
	}
}

// CreateNamespace takes the name of a client to create and returns any error that occurred when creating it.
func (c *Client) CreateNamespace(name string) error {
	newNamespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}

	_, err := c.client.Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
	return err
}

// ListNamespaces fetches all the client names and returns them as a list as well as any error.
func (c *Client) ListNamespaces() ([]string, error) {
	namespaceMeta, err := c.client.Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	var namespaces []string
	for _, ns := range namespaceMeta.Items {
		namespaces = append(namespaces, ns.Name)
	}

	return namespaces, nil
}

// DeleteNamespace deletes the given client and returns any error.
func (c *Client) DeleteNamespace(name string) error {
	return c.client.Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
}

// CreatePod creates the given pod in the given namespace.
func (c *Client) CreatePod(namespace string, podName string, opts ...func(pod *v1.Pod)) error {
	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   podName,
			Labels: map[string]string{"k8s-app": "kube-dns"},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  podName,
					Image: podName,
					Ports: []v1.ContainerPort{
						{
							Protocol:      "TCP",
							HostPort:      80,
							ContainerPort: 80,
						},
					},
				},
			},
		},
		Status: v1.PodStatus{},
	}

	for _, opt := range opts {
		opt(pod)
	}

	_, err := c.client.Pods(namespace).Create(context.Background(), pod, metav1.CreateOptions{})
	return err
}

// List pods returns a list of the pods for all namesapces that match the given list options.
func (c *Client) ListPods(opts metav1.ListOptions) ([]v1.Pod, error) {
	var res []v1.Pod
	namespaces, err := c.client.Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return []v1.Pod{}, err
	}

	for _, namespace := range namespaces.Items {
		podClient := c.client.Pods(namespace.Name)
		pods, err := podClient.List(context.Background(), opts)
		if err != nil {
			return []v1.Pod{}, err
		}

		for _, pod := range pods.Items {
			res = append(res, pod)
		}
	}

	return res, nil
}

// GetClientset checks the home directory for the kubernetes cluster config under ~/.kube/config. Returns a kubernetes
// Clientset and any error.
// TODO add unit tests.
func GetClientset() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
