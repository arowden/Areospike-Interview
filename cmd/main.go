/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8sexample/client"
	"log"
	"os"
	"time"
)

func main() {
	clientset, err := client.GetClientset()
	if err != nil {
		log.Printf("Error occurred while fetching clientset. %s", err.Error())
		os.Exit(1)
	}

	nsManager := client.NewClient(clientset.CoreV1())
	newNamespace := flag.String("newNamespace", "", "A namespace to create")
	flag.Parse()

	if *newNamespace == "" {
		log.Print("Namespace must be set using the --newNamespace flag")
		os.Exit(1)
	}

	err = Execute(nsManager, *newNamespace)
	if err != nil {
		os.Exit(1)
	}
}

func Execute(client *client.Client, namespaceToCreate string) error {
	err := printNamespaces(client)
	if err != nil {
		return err
	}

	err = client.CreateNamespace(namespaceToCreate)
	if err != nil {
		log.Printf("Failed to create new namespace. Error: %s", err.Error())
		return err
	}

	err = printNamespaces(client)
	if err != nil {
		return err
	}

	err = client.CreatePod(namespaceToCreate, "hello-world")
	if err != nil {
		log.Printf("Failed to create pod. Error: %s", err.Error())
		return err
	}

	pods, err := client.ListPods(metav1.ListOptions{LabelSelector: "k8s-app=kube-dns"})
	if err != nil {
		log.Printf("Failed to list pods. Error: %s", err.Error())
		return err
	}

	for _, pod := range pods {
		log.Printf("Found pod. Namespace: %s, Pod Name: %s", pod.ObjectMeta.Namespace, pod.Name)
	}

	err = client.DeleteNamespace(namespaceToCreate)
	if err != nil {
		log.Printf("Failed to delete namespace. Error: %s)", err.Error())
	}

	// Wait for namespaces to be deleted
	time.Sleep(time.Second * 50)
	return printNamespaces(client)
}

func printNamespaces(nsClient *client.Client) error {
	namespaces, err := nsClient.ListNamespaces()
	if err != nil {
		log.Printf("Failed to list namespaces. Error: %s", err.Error())
		return err
	}

	log.Printf("Existing Namespaces: %v", namespaces)
	return nil
}
