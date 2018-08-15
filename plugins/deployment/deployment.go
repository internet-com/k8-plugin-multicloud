package main

import (
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"

	pkgerrors "github.com/pkg/errors"

	appsV1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/shank7485/k8-plugin-multicloud/krd"
)

// CreateResource object in a specific Kubernetes Deployment
func CreateResource(kubedata *krd.GenericKubeResourceData, kubeclient *kubernetes.Clientset) (string, error) {
	if kubedata.Namespace == "" {
		kubedata.Namespace = "default"
	}

	if _, err := os.Stat(kubedata.YamlFilePath); err == nil {
		log.Println("Reading deployment YAML")

		rawBytes, err := ioutil.ReadFile(kubedata.YamlFilePath)
		if err != nil {
			return "", pkgerrors.Wrap(err, "Deployment YAML file read error")
		}

		log.Println("Decoding deployment YAML")

		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, _, err := decode(rawBytes, nil, nil)
		if err != nil {
			return "", pkgerrors.Wrap(err, "Deserialize deployment error")
		}

		switch o := obj.(type) {
		case *appsV1.Deployment:
			kubedata.DeploymentData = o
		}

		// cloud1-default-uuid-sisedeploy
		internalDeploymentName := kubedata.InternalVNFID + "-" + kubedata.DeploymentData.Name

		kubedata.DeploymentData.Namespace = kubedata.Namespace
		kubedata.DeploymentData.Name = internalDeploymentName

		result, err := kubeclient.AppsV1().Deployments(kubedata.Namespace).Create(kubedata.DeploymentData)
		if err != nil {
			return "", pkgerrors.Wrap(err, "Create Deployment error")
		}

		return result.GetObjectMeta().GetName(), nil
	}
	return "", pkgerrors.New("File " + kubedata.YamlFilePath + " not found")
}

// ListResources of existing deployments hosted in a specific Kubernetes Deployment
func ListResources(limit int64, namespace string, kubeclient *kubernetes.Clientset) (*[]string, error) {
	if namespace == "" {
		namespace = "default"
	}

	opts := metaV1.ListOptions{
		Limit: limit,
	}
	opts.APIVersion = "apps/v1"
	opts.Kind = "Deployment"

	list, err := kubeclient.AppsV1().Deployments(namespace).List(opts)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "Get Deployment list error")
	}

	result := make([]string, 0, limit)
	if list != nil {
		for _, deployment := range list.Items {
			result = append(result, deployment.Name)
		}
	}

	return &result, nil
}

// DeleteResource existing deployments hosting in a specific Kubernetes Deployment
func DeleteResource(name string, namespace string, kubeclient *kubernetes.Clientset) error {
	if namespace == "" {
		namespace = "default"
	}

	log.Println("Deleting deployment: " + name)

	deletePolicy := metaV1.DeletePropagationForeground
	err := kubeclient.AppsV1().Deployments(namespace).Delete(name, &metaV1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		return pkgerrors.Wrap(err, "Delete Deployment error")
	}

	return nil
}

// GetResource existing deployment hosting in a specific Kubernetes Deployment
func GetResource(name string, namespace string, kubeclient *kubernetes.Clientset) (string, error) {
	if namespace == "" {
		namespace = "default"
	}

	opts := metaV1.ListOptions{
		Limit: 10,
	}
	opts.APIVersion = "apps/v1"
	opts.Kind = "Deployment"

	list, err := kubeclient.AppsV1().Deployments(namespace).List(opts)
	if err != nil {
		return "", pkgerrors.Wrap(err, "Get Deployment error")
	}

	for _, deployment := range list.Items {
		if deployment.Name == name {
			return name, nil
		}
	}
	return "", nil
}
